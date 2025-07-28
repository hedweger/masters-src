package device

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"scada-simu/internal/templates"
)

type DeviceType string

const (
	TypeRTU    DeviceType = "rtu"
	TypeSwitch DeviceType = "switch"
)

type NetworkConnection struct {
	Name      string
	Interface string
	SourceIP  string
	MAC       string
	Gateway   net.IP
}

type Device struct {
	Type              DeviceType
	Name              string
	Networks          []NetworkConnection
	ImagePath         string
	SeedImagePath     string
	UserDataPath      string
	MetaDataPath      string
	NetworkConfigPath string
	ifaceCount        int
}

func (d *Device) CreateLibvirtConfig(outputDir string) error {
	if d.ImagePath == "" {
		return fmt.Errorf("device %s does not have an image path set", d.Name)
	}
	if d.SeedImagePath == "" {
		return fmt.Errorf("device %s does not have a seed image path set", d.Name)
	}
	outpath := fmt.Sprintf("%s/%s/", outputDir, d.Name)

	deviceContext := templates.VirtDeviceContext{
		Name:     d.Name,
		RAM:      "512",
		VCPU:     1,
		DiskPath: d.ImagePath,
		SeedPath: d.SeedImagePath,
		Networks: make([]templates.VirtDevNetworkContext, 0, len(d.Networks)),
	}

	for _, conn := range d.Networks {
		deviceContext.Networks = append(deviceContext.Networks, templates.VirtDevNetworkContext{
			Name: conn.Name,
			MAC:  conn.MAC,
		})
	}

	var err error
	d.ImagePath, err = templates.WriteVirtDevice(outpath, deviceContext)
	if err != nil {
		log.Fatalf("%s: %v", d.Name, err)
	}
	return nil
}

func (d *Device) CreateCloudInitConfig(outputDir string) {
	seed_path := fmt.Sprintf("%s/%s/seed", outputDir, d.Name)
	os.MkdirAll(seed_path, 0755)
	var err error

	d.UserDataPath, err = templates.WriteUserData(seed_path, templates.UserDataContext{
		Hostname: d.Name,
		Password: "root",
		Commands: d.commands(),
	})
	if err != nil {
		log.Fatalf("%s: %s", d.Name, err)
	}

	netwCfgs := make([]templates.NetworkContext, 0, len(d.Networks))
	for _, conn := range d.Networks {
		netwCfgs = append(netwCfgs, templates.NetworkContext{
			Interface: conn.Interface,
			SourceIP:  conn.SourceIP,
			Gateway:   conn.Gateway.String(),
		})
	}
	d.NetworkConfigPath, err = templates.WriteNetworkConfig(seed_path, templates.NetworkConfigContext{
		Connections: netwCfgs,
		DeviceType:  string(d.Type),
	})
	if err != nil {
		log.Fatalf("%s: %s", d.Name, err)
	}

	metaDataPath := fmt.Sprintf("%s/meta-data", seed_path)
	if err := os.WriteFile(metaDataPath, []byte(""), 0644); err != nil {
		log.Fatalf("%s: %v", d.Name, err)
	}
	d.MetaDataPath, _ = filepath.Abs(metaDataPath)
}

func (d *Device) CreateSeedImage(outputDir string) error {
	outpath := fmt.Sprintf("%s/%s/seed.iso", outputDir, d.Name)
	outpath, err := filepath.Abs(outpath)
	if err != nil {
		return fmt.Errorf("Failed to find absolute path for %s: %w", outpath, err)
	}
	executor := "genisoimage"

	cmd := exec.Command(executor, "-o", outpath, "-volid", "cidata", "-joliet", "-rock", d.UserDataPath, d.NetworkConfigPath, d.MetaDataPath)
	log.Printf("[INFO] Executing %s", cmd.String())
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("Failed to execute %s in %s: %w", executor, outpath, err)
	}
	d.SeedImagePath, err = filepath.Abs(outpath)
	if err != nil {
		return fmt.Errorf("Failed to find %s, %w", d.SeedImagePath, err)
	}
	return nil
}

func (d *Device) AddNetworkConnection(name string, to string, src_ip string, mac string) error {
	if d.Type == TypeSwitch {
		src_ip = ""
	}

	iface := fmt.Sprintf("ens%d", d.ifaceCount)
	d.Networks = append(d.Networks, NetworkConnection{
		Name:      name,
		Interface: iface,
		SourceIP:  src_ip,
		Gateway:   nil,
		MAC:       mac,
	})
	log.Printf("[INFO] Added L2 connection %s to device %s on %s", name, d.Name, iface)
	d.ifaceCount++
	return nil
}

func (d *Device) commands() []string {
	return []string{"etst", "test"}
}
