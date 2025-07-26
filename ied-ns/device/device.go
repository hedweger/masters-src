package device

import (
	"embed"
	"fmt"
	"ied-ns/cinit"
	"strings"
	"text/template"
)

//go:embed templates/*
var templates embed.FS

// DeviceType represents the type of device
type DeviceType string

const (
	SW  DeviceType = "switch"
	RTU DeviceType = "rtu"
)

// Device represents a network device
type Device struct {
	DevType      DeviceType
	Networks     []cinit.NetworkConnection
	Name         string
	Address      string
	IfaceCount   int
	ImagePath    *string
	SeedPath     *string
	UserDataPath *string
	CloudDataPath *string
}

// NewDevice creates a new Device instance
func NewDevice(devType DeviceType, name, address string) *Device {
	// Add /24 suffix if not present
	if !strings.Contains(address, "/") {
		address = address + "/24"
	}
	
	return &Device{
		DevType:    devType,
		Networks:   make([]cinit.NetworkConnection, 0),
		Name:       name,
		Address:    address,
		IfaceCount: 0,
	}
}

// AddNetworkConnection adds a network connection to the device
func (d *Device) AddNetworkConnection(networkName, mac string, gateway *string) {
	d.IfaceCount++
	conn := cinit.NetworkConnection{
		Name:    networkName,
		Iface:   fmt.Sprintf("ens%d", d.IfaceCount+1),
		SrcIP:   d.Address,
		MAC:     mac,
		Gateway: gateway,
	}
	d.Networks = append(d.Networks, conn)
}

// StartupCommands returns the startup commands for the device
func (d *Device) StartupCommands() []string {
	if d.DevType == SW {
		return []string{
			"ip link add name br0 type bridge",
			"ip link set dev ens2 master br0",
			"ip link set dev ens3 master br0",
			"ip link set dev br0 up",
			"sudo tcpdump -i br0 not arp and not llc",
		}
	} else {
		cmds := []string{
			"sudo wget https://github.com/hedweger/masters-src/releases/download/client/ied-client",
			"sudo chmod +x /ied-client",
		}
		if d.Name == "pc1" {
			cmds = append(cmds,
				"sudo wget https://github.com/hedweger/masters-src/releases/download/server/ied-server.tar",
				"sudo tar -xf /ied-server.tar",
				"sudo chmod +x /ied-server",
				"sudo /ied-server",
			)
		}
		return cmds
	}
}

// StartupFilewrites returns the file writes for device startup
func (d *Device) StartupFilewrites() []cinit.FileWrite {
	// For now, return empty slice (same as Python version)
	return []cinit.FileWrite{}
}

// LibvirtXMLData represents data for libvirt XML template
type LibvirtXMLData struct {
	DType string
	Name  string
	NRAM  string
	VCPU  string
	Disk  string
	Seed  string
	Nets  []NetworkInterface
}

// NetworkInterface represents a network interface for XML template
type NetworkInterface struct {
	Name string
	MAC  string
}

// LibvirtXML generates the libvirt XML configuration for the device
func (d *Device) LibvirtXML() (string, error) {
	tmpl, err := template.ParseFS(templates, "templates/virt_device.xml.tmpl")
	if err != nil {
		return "", fmt.Errorf("failed to parse libvirt template: %w", err)
	}
	
	// Convert networks to NetworkInterface format
	nets := make([]NetworkInterface, len(d.Networks))
	for i, net := range d.Networks {
		nets[i] = NetworkInterface{
			Name: net.Name,
			MAC:  net.MAC,
		}
	}
	
	data := LibvirtXMLData{
		DType: string(d.DevType),
		Name:  d.Name,
		NRAM:  "512", // for now
		VCPU:  "1",   // for now
		Disk:  *d.ImagePath,
		Seed:  *d.SeedPath,
		Nets:  nets,
	}
	
	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute libvirt template: %w", err)
	}
	
	return buf.String(), nil
}