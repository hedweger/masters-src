package device

import (
	"fmt"
	"log"
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
	cfg            *config.Config
	outputDir      string
	cleanDrivePath string
	ifaceCount     int
}

func InitManager(cfg *config.Config, outputDir string) *Manager {
	ipgen, err := DefaultIPGenerator(cfg.Network.Address)
	if err != nil {
		log.Fatalf("[ERROR] %v", err)
	}
	manager := &Manager{
		Devices:        make(map[string]*Device),
		MacGen:         DefaultMACGenerator(),
		IpGen:          ipgen,
		Networks:       []string{},
		outputDir:      outputDir,
		cleanDrivePath: "../debian-12-genericcloud-amd64.qcow2",
		ifaceCount:     1,
	}

	return manager
}

func (m *Manager) Deploy() {
	if err := m.initializeRTUs(m.cfg); err != nil {
		log.Fatalf("[ERROR] %s", err)
	}

	if err := m.initializeSwitches(m.cfg); err != nil {
		log.Fatalf("[ERROR] %s", err)
	}
	m.prepareDevices()
}

func (m *Manager) StartVMs() {
	fmt.Println("")

	conn, err := libvirt.NewConnect("qemu:///system")
	if err != nil {
		log.Fatalf("[ERROR] Failed to connect to libvirt: %v", err)
	}
	defer conn.Close()

	for _, networkXmlPath := range m.Networks {
		networkName := strings.TrimSuffix(filepath.Base(networkXmlPath), ".xml")
		fmt.Printf("Defining/starting network: %s\n", networkName)

		netw, err := conn.LookupNetworkByName(networkName)
		if err != nil {
			xmlBytes, readErr := os.ReadFile(networkXmlPath)
			if readErr != nil {
				log.Fatalf("[ERROR] Failed to read network XML: %s: %v", networkXmlPath, readErr)
			}
			netw, err = conn.NetworkDefineXML(string(xmlBytes))
			if err != nil {
				log.Fatalf("[ERROR] Failed to define network %s: %v", networkName, err)
			}
			defer netw.Free()
		} else {
			defer netw.Free()
		}

		active, err := netw.IsActive()
		if err != nil {
			log.Fatalf("[ERROR] Failed to check if network %s is active: %v", networkName, err)
		}
		if !active {
			if err := netw.Create(); err != nil {
				log.Fatalf("[ERROR] Failed to start network %s: %v", networkName, err)
			}
		}
	}

	for _, device := range m.Devices {
		deviceXmlPath := filepath.Join(m.outputDir, device.Name, device.Name+".xml")
		fmt.Printf("Defining/starting device: %s\n", device.Name)

		dom, err := conn.LookupDomainByName(device.Name)
		if err == nil {
			defer dom.Free()
			active, _ := dom.IsActive()
			if active {
				if err := dom.Destroy(); err != nil {
					log.Printf("[WARN] Failed to destroy running domain %s: %v", device.Name, err)
				}
			}
			if err := dom.UndefineFlags(libvirt.DOMAIN_UNDEFINE_MANAGED_SAVE | libvirt.DOMAIN_UNDEFINE_SNAPSHOTS_METADATA | libvirt.DOMAIN_UNDEFINE_NVRAM); err != nil {
				log.Printf("[WARN] Failed to undefine domain %s: %v", device.Name, err)
			}
		}

		xmlBytes, readErr := os.ReadFile(deviceXmlPath)
		if readErr != nil {
			log.Fatalf("[ERROR] Failed to read domain XML: %s: %v", deviceXmlPath, readErr)
		}
		dom, err = conn.DomainDefineXML(string(xmlBytes))
		if err != nil {
			log.Fatalf("[ERROR] Failed to define domain %s: %v", device.Name, err)
		}
		defer dom.Free()

		if err := dom.Create(); err != nil {
			log.Fatalf("[ERROR] Failed to start domain %s: %v", device.Name, err)
		}
	}
}

func (m *Manager) prepareDevices() {
	for _, device := range m.Devices {
		device.CreateCloudInitConfig(m.outputDir)
		if err := device.CreateSeedImage(m.outputDir); err != nil {
			log.Fatalf("[ERROR] %v", err)
		}
		device.CreateSeedImage(m.outputDir)
		if err := device.CreateLibvirtConfig(m.outputDir); err != nil {
			log.Fatalf("[ERROR] %v", err)
		}
	}
}

func (m *Manager) initializeSwitches(cfg *config.Config) error {
	for i, sw := range cfg.Switches {
		if sw.Name == "" {
			sw.Name = fmt.Sprintf("sw%d", i+1)
		}

		// @TODO: Do switches need addresses?
		// For later management we should at least have a connection to historian.
		// if sw.Address != "" {
		// 	err := m.Netman.ValidateAddress(sw.Address)
		// 	if err != nil {
		// 		return fmt.Errorf("failed to validate address for switch %s: %w", sw.Name, err)
		// 	}
		// }

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
			ImagePath:  image,
			ifaceCount: 2,
		}

		for j, conn := range sw.Connected {
			if conn.To == "" {
				log.Println("[WARN] Switch", sw.Name, "connection", j+1, "has no target device, skipping")
				continue
			}
			net_name := fmt.Sprintf("%s-%s", sw.Name, conn.To)
			net_bridge := fmt.Sprintf("virbr%d", m.ifaceCount)
			device.AddNetworkConnection(net_name, device.Name, m.IpGen.NextWCidr(), m.MacGen.Next())
			dest := m.Devices[conn.To]
			if dest == nil {
				panic(fmt.Errorf("%s: 'To' device %s does not exist.", device.Name, conn.To))
			}
			dest.AddNetworkConnection(net_name, device.Name, m.IpGen.NextWCidr(), m.MacGen.Next())
			m.Networks = append(m.Networks, m.createNetwork(net_name, net_bridge))
			m.ifaceCount++
		}

		m.Devices[sw.Name] = device
	}
	return nil
}

func (m *Manager) createNetwork(name string, bridge string) string {
	path, err := templates.WriteVirtNetwork(m.outputDir, templates.VirtNetworkContext{
		Name:   name,
		Bridge: bridge,
	})
	if err != nil {
		log.Fatalf("[ERROR] %s: %v", name, err)
	}
	return path
}

func (m *Manager) initializeRTUs(cfg *config.Config) error {
	for i, rtu := range cfg.RTUs {
		if rtu.Name == "" {
			rtu.Name = fmt.Sprintf("rtu%d", i+1)
		}

		if rtu.Address == "" {
			panic("@TODO: empty address for RTUs")
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
			ImagePath:  image,
			ifaceCount: 2,
		}

		m.Devices[rtu.Name] = device
	}
	return nil
}
