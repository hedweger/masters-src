package device

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"scada-simu/internal/templates"
	"strconv"
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
	Gateway   string
}

type Device struct {
	Type              DeviceType
	Name              string
	Memory            int
	VCPU              int
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
		RAM:      strconv.Itoa(d.Memory),
		VCPU:     d.VCPU,
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
		return err
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
		Packages: d.packages(),
	})
	if err != nil {
		slog.Error("failed to write user data", "device", d.Name, "error", err)
		return
	}

	netwCfgs := make([]templates.NetworkContext, 0, len(d.Networks))
	for _, conn := range d.Networks {
		netwCfgs = append(netwCfgs, templates.NetworkContext{
			Interface: conn.Interface,
			SourceIP:  conn.SourceIP,
			Gateway:   conn.Gateway,
		})
	}
	d.NetworkConfigPath, err = templates.WriteNetworkConfig(seed_path, templates.NetworkConfigContext{
		Connections: netwCfgs,
		DeviceType:  string(d.Type),
	})
	if err != nil {
		slog.Error("failed to write network config", "device", d.Name, "error", err)
		return
	}

	metaDataPath := fmt.Sprintf("%s/meta-data", seed_path)
	if err := os.WriteFile(metaDataPath, []byte(""), 0644); err != nil {
		slog.Error("failed to write meta-data file", "device", d.Name, "error", err)
		return
	}
	d.MetaDataPath, _ = filepath.Abs(metaDataPath)
}

func (d *Device) CreateSeedImage(outputDir string) error {
	outpath := fmt.Sprintf("%s/%s/seed.iso", outputDir, d.Name)
	outpath, err := filepath.Abs(outpath)
	if err != nil {
		return err
	}
	executor := "genisoimage"

	cmd := exec.Command(executor, "-o", outpath, "-volid", "cidata", "-joliet", "-rock", d.UserDataPath, d.NetworkConfigPath, d.MetaDataPath)
	slog.Info("Creating seed image", "command", cmd.String(), "output", outpath)
	err = cmd.Run()
	if err != nil {
		return err
	}
	d.SeedImagePath, err = filepath.Abs(outpath)
	if err != nil {
		return err
	}
	return nil
}

func (d *Device) AddNetworkConnection(name string, to string, gw string, src_ip string, mac string) error {
	iface := fmt.Sprintf("ens%d", d.ifaceCount)
	d.Networks = append(d.Networks, NetworkConnection{
		Name:      name,
		Interface: iface,
		SourceIP:  src_ip,
		Gateway:   gw,
		MAC:       mac,
	})
	slog.Info("Added L2 connection", "device", d.Name, "connection", name, "interface", iface)
	d.ifaceCount++
	return nil
}

func (d *Device) packages() []string {
	return []string{
		"iperf3",
		"bash",
	}
}

func (d *Device) commands() []string {
	return []string{""}
}
