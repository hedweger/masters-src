package device

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"scada-simu/internal/config"
	"scada-simu/internal/templates"
	"scada-simu/internal/virt"
	"strings"

	"libvirt.org/go/libvirt"
)

type Manager struct {
	Devices        map[string]*Device
	MacGen         *MACGen
	IpGen          *IPGen
	Networks       []string
	Config         *config.Config
	outputDir      string
	cleanDrivePath string
	ifaceCount     int
}

func InitManager(cfg *config.Config, outputDir string) *Manager {
	ipgen, err := DefaultIPGenerator(cfg.Network.Address)
	if err != nil {
		slog.Error("Failed to initialize IP generator", "error", err)
		return nil
	}
	manager := &Manager{
		Devices:        make(map[string]*Device),
		MacGen:         DefaultMACGenerator(),
		IpGen:          ipgen,
		Networks:       []string{},
		outputDir:      outputDir,
		cleanDrivePath: "/home/th/workspace/masters/debian-12-genericcloud-amd64.qcow2",
		ifaceCount:     1,
	}

	return manager
}

func (m *Manager) Deploy() {
	if err := m.initializeRTUs(m.Config); err != nil {
		slog.Error("Failed to initialize RTUs", "error", err)
		return
	}

	if err := m.initializeSwitches(m.Config); err != nil {
		slog.Error("Failed to initialize switches", "error", err)
		return
	}
	m.prepareDevices()
}

func (m *Manager) StartVMs() {
	conn, err := libvirt.NewConnect("qemu:///system")
	if err != nil {
		slog.Error("Failed to connect to libvirt", "error", err)
		return
	}
	defer conn.Close()

	for _, networkXmlPath := range m.Networks {
		networkName := strings.TrimSuffix(filepath.Base(networkXmlPath), ".xml")
		netw, err := conn.LookupNetworkByName(networkName)
		if err != nil {
			xmlBytes, readErr := os.ReadFile(networkXmlPath)
			if readErr != nil {
				slog.Error("Failed to read network XML", "path", networkXmlPath, "error", readErr)
				return
			}
			netw, err = conn.NetworkDefineXML(string(xmlBytes))
			if err != nil {
				slog.Error("Failed to define network", "name", networkName, "error", err)
				return
			}
			defer netw.Free()
		} else {
			defer netw.Free()
		}

		active, err := netw.IsActive()
		if err != nil {
			slog.Error("Failed to check network status", "name", networkName, "error", err)
			return
		}
		if !active {
			if err := netw.Create(); err != nil {
				slog.Error("Failed to start network", "name", networkName, "error", err)
				return
			}
		}
	}

	for _, device := range m.Devices {
		deviceXmlPath := filepath.Join(m.outputDir, device.Name, device.Name+".xml")
		dom, err := conn.LookupDomainByName(device.Name)
		if err == nil {
			defer dom.Free()
			active, _ := dom.IsActive()
			if active {
				if err := dom.Destroy(); err != nil {
					slog.Warn("Failed to destroy active domain", "name", device.Name, "error", err)
				}
			}
			if err := dom.UndefineFlags(libvirt.DOMAIN_UNDEFINE_MANAGED_SAVE | libvirt.DOMAIN_UNDEFINE_SNAPSHOTS_METADATA | libvirt.DOMAIN_UNDEFINE_NVRAM); err != nil {
				slog.Warn("Failed to undefine domain", "name", device.Name, "error", err)
			}
		}

		xmlBytes, readErr := os.ReadFile(deviceXmlPath)
		if readErr != nil {
			slog.Error("Failed to read device XML", "path", deviceXmlPath, "error", readErr)
			return
		}
		dom, err = conn.DomainDefineXML(string(xmlBytes))
		if err != nil {
			slog.Error("Failed to define domain", "name", device.Name, "error", err)
			return
		}
		defer dom.Free()

		if err := dom.Create(); err != nil {
			slog.Error("Failed to start domain", "name", device.Name, "error", err)
			return
		}
	}
}

func (m *Manager) prepareDevices() {
	for _, device := range m.Devices {
		device.CreateCloudInitConfig(m.outputDir)
		if err := device.CreateSeedImage(m.outputDir); err != nil {
			slog.Error("Failed to create seed image", "device", device.Name, "error", err)
			return
		}
		device.CreateSeedImage(m.outputDir)
		if err := device.CreateLibvirtConfig(m.outputDir); err != nil {
			slog.Error("Failed to create libvirt config", "device", device.Name, "error", err)
			return
		}
	}
}

func (m *Manager) initializeSwitches(cfg *config.Config) error {
	for i, sw := range cfg.Switches {
		if sw.Name == "" {
			sw.Name = fmt.Sprintf("sw%d", i+1)
		}

		device_path := fmt.Sprintf("%s/%s", m.outputDir, sw.Name)
		os.Mkdir(device_path, 0755)
		image, err := virt.CreateQcow2Image(m.cleanDrivePath, device_path, sw.Name)
		if err != nil {
			return fmt.Errorf("Failed to copy .qcow2 image: %s, %s, %w", m.cleanDrivePath, m.outputDir, err)
		}
		image, _ = filepath.Abs(image) // file must exist because we just created it

		device := &Device{
			Type:       TypeSwitch,
			Name:       sw.Name,
			Memory:     sw.Memory,
			VCPU:       sw.VCPU,
			ImagePath:  image,
			ifaceCount: 2,
		}
		device.AddNetworkConnection("default", "default", "192.168.122.1", m.IpGen.NextWCidr(), m.MacGen.Next())

		for j, conn := range sw.Connected {
			if conn.To == "" {
				slog.Warn(fmt.Sprintf("%s: Connection %d has no 'To' device specified, skipping.", sw.Name, j+1))
				continue
			}
			net_name := fmt.Sprintf("%s-%s", sw.Name, conn.To)
			net_bridge := fmt.Sprintf("virbr%d", m.ifaceCount)
			// SW side of connection should not have a gateway, so we use an empty string
			device.AddNetworkConnection(net_name, device.Name, "", m.IpGen.NextWCidr(), m.MacGen.Next())
			dest := m.Devices[conn.To]
			if dest == nil {
				return fmt.Errorf("%s: 'To' device %s does not exist.", device.Name, conn.To)
			}
			dest.AddNetworkConnection(net_name, device.Name, "192.168.122.1", m.IpGen.NextWCidr(), m.MacGen.Next())
			netw, err := m.createNetwork(net_name, net_bridge)
			if err != nil {
				return fmt.Errorf("Failed to create network %s: %w", net_name, err)
			}
			m.Networks = append(m.Networks, netw)
			m.ifaceCount++
		}

		m.Devices[sw.Name] = device
	}
	return nil
}

func (m *Manager) createNetwork(name string, bridge string) (string, error) {
	path, err := templates.WriteVirtNetwork(m.outputDir, templates.VirtNetworkContext{
		Name:   name,
		Bridge: bridge,
	})
	if err != nil {
		return "", err
	}
	return path, nil
}

func (m *Manager) initializeRTUs(cfg *config.Config) error {
	for i, rtu := range cfg.RTUs {
		if rtu.Name == "" {
			rtu.Name = fmt.Sprintf("rtu%d", i+1)
		}

		if rtu.Address == "" {
			rtu.Address = m.IpGen.NextWCidr()
		}

		device_path := fmt.Sprintf("%s/%s", m.outputDir, rtu.Name)
		os.Mkdir(device_path, 0755)

		image, err := virt.CreateQcow2Image(m.cleanDrivePath, device_path, rtu.Name)
		if err != nil {
			return fmt.Errorf("Failed to copy .qcow2 image: %s, %s, %w", m.cleanDrivePath, m.outputDir, err)
		}
		image, _ = filepath.Abs(image) // file must exist because we just created it

		device := &Device{
			Type:       TypeRTU,
			Name:       rtu.Name,
			Memory:     rtu.Memory,
			VCPU:       rtu.VCPU,
			ImagePath:  image,
			ifaceCount: 2,
		}

		m.Devices[rtu.Name] = device
	}
	return nil
}
