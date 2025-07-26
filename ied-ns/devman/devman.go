package devman

import (
	"embed"
	"fmt"
	"ied-ns/cinit"
	"ied-ns/device"
	"ied-ns/drive"
	"ied-ns/virtmac"
	"net"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"
)

//go:embed templates/*
var templates embed.FS

// Config represents the YAML configuration structure
type Config struct {
	Network struct {
		Address string `yaml:"address"`
	} `yaml:"network"`
	RTUs []struct {
		Name    string `yaml:"name"`
		Address string `yaml:"address"`
	} `yaml:"rtus"`
	Switches []struct {
		Name      string `yaml:"name"`
		Address   string `yaml:"address"`
		Connected []struct {
			To string `yaml:"to"`
		} `yaml:"connected"`
	} `yaml:"switches"`
}

// DeviceManager manages devices and their configurations
type DeviceManager struct {
	Devices        map[string]*device.Device
	Context        string
	NetworkAddress *net.IPNet
	MacGen         func() string
	IP4Iter        []net.IP
	IP4Index       int
}

// NewDeviceManager creates a new DeviceManager instance
func NewDeviceManager() *DeviceManager {
	return &DeviceManager{
		Devices: make(map[string]*device.Device),
		MacGen:  virtmac.Gen(),
	}
}

// Parse parses the configuration file and creates devices
func (dm *DeviceManager) Parse(cfgFP string, write bool) error {
	data, err := os.ReadFile(cfgFP)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse config YAML: %w", err)
	}

	// Create file path context for later use
	cfgDir := filepath.Dir(cfgFP)
	dm.Context = filepath.Join(cfgDir, "tmp")

	// Parse network address
	_, network, err := net.ParseCIDR(config.Network.Address)
	if err != nil {
		return fmt.Errorf("failed to parse network address: %w", err)
	}
	dm.NetworkAddress = network

	// Generate list of available IPs
	dm.IP4Iter = generateHosts(network)
	dm.IP4Index = 0

	// Create RTU devices
	for i, rtu := range config.RTUs {
		name := rtu.Name
		if name == "" {
			name = fmt.Sprintf("rtu%d", i)
		}

		address := rtu.Address
		if address == "" {
			address = fmt.Sprintf("%s/24", dm.nextIP())
		}

		dm.Devices[name] = device.NewDevice(device.RTU, name, address)
	}

	// Create switch devices
	for i, sw := range config.Switches {
		name := sw.Name
		if name == "" {
			name = fmt.Sprintf("switch%d", i)
		}

		address := sw.Address
		if address == "" {
			address = fmt.Sprintf("%s/24", dm.nextIP())
		}

		dm.Devices[name] = device.NewDevice(device.SW, name, address)

		// Create network connections
		for _, conn := range sw.Connected {
			if err := dm.createNetwork(name, conn.To, ""); err != nil {
				return fmt.Errorf("failed to create network connection: %w", err)
			}
		}
	}

	return dm.createDevices(write)
}

// nextIP returns the next available IP address
func (dm *DeviceManager) nextIP() string {
	if dm.IP4Index >= len(dm.IP4Iter) {
		return "192.168.1.1" // fallback
	}
	ip := dm.IP4Iter[dm.IP4Index]
	dm.IP4Index++
	return ip.String()
}

// generateHosts generates all host addresses in a network
func generateHosts(network *net.IPNet) []net.IP {
	var hosts []net.IP
	
	// Get the first and last IP
	ip := network.IP.Mask(network.Mask)
	
	// Increment through all IPs in the network
	for ip := ip.Mask(network.Mask); network.Contains(ip); inc(ip) {
		// Skip network and broadcast addresses
		if !ip.Equal(network.IP) && !isBroadcast(ip, network) {
			hosts = append(hosts, make(net.IP, len(ip)))
			copy(hosts[len(hosts)-1], ip)
		}
	}
	
	return hosts
}

// inc increments an IP address
func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

// isBroadcast checks if an IP is the broadcast address
func isBroadcast(ip net.IP, network *net.IPNet) bool {
	broadcast := make(net.IP, len(network.IP))
	copy(broadcast, network.IP)
	
	for i := 0; i < len(broadcast); i++ {
		broadcast[i] |= ^network.Mask[i]
	}
	
	return ip.Equal(broadcast)
}

// createDevices creates all device configurations
func (dm *DeviceManager) createDevices(write bool) error {
	for _, dev := range dm.Devices {
		// Store original address
		devAddr := dev.Address

		// Add default network connection
		dev.Address = fmt.Sprintf("192.168.122.%d/24", extractDeviceNumber(dev.Name)+1)
		gateway := "192.168.122.1"
		dev.AddNetworkConnection("default", dm.MacGen(), &gateway)

		// Restore original address
		dev.Address = devAddr

		// Create image path
		imagePath, err := drive.QCOW2(filepath.Join(dm.Context, dev.Name), dev.Name, write)
		if err != nil {
			return fmt.Errorf("failed to create QCOW2 image: %w", err)
		}
		dev.ImagePath = &imagePath

		// Prepare cloud-init
		seeds, err := cinit.Prepare(
			string(dev.DevType),
			dev.Name,
			dev.StartupCommands(),
			dev.StartupFilewrites(),
			dev.Networks,
			filepath.Join(dm.Context, dev.Name),
			write,
		)
		if err != nil {
			return fmt.Errorf("failed to prepare cloud-init: %w", err)
		}

		dev.SeedPath = &seeds.ISOP
		dev.UserDataPath = &seeds.UserData
		dev.CloudDataPath = &seeds.CloudData

		// Generate libvirt XML
		if write {
			xmlContent, err := dev.LibvirtXML()
			if err != nil {
				return fmt.Errorf("failed to generate libvirt XML: %w", err)
			}

			configPath := filepath.Join(dm.Context, dev.Name, "config.xml")
			if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
				return fmt.Errorf("failed to create config directory: %w", err)
			}

			if err := os.WriteFile(configPath, []byte(xmlContent), 0644); err != nil {
				return fmt.Errorf("failed to write config.xml: %w", err)
			}
		}
	}

	return nil
}

// extractDeviceNumber extracts the number from device name (e.g., "pc1" -> 1)
func extractDeviceNumber(name string) int {
	// Simple extraction - get last character if it's a digit
	if len(name) > 0 {
		lastChar := name[len(name)-1]
		if lastChar >= '0' && lastChar <= '9' {
			return int(lastChar - '0')
		}
	}
	return 0
}

// NetworkTemplateData represents data for network XML template
type NetworkTemplateData struct {
	Name string
	Brdg string
}

// createNetwork creates a network connection between two devices
func (dm *DeviceManager) createNetwork(srcName, dstName, gtwAddr string) error {
	dstDevice, exists := dm.Devices[dstName]
	if !exists {
		return fmt.Errorf("destination device %s not found", dstName)
	}

	srcDevice, exists := dm.Devices[srcName]
	if !exists {
		return fmt.Errorf("source device %s not found", srcName)
	}

	gateway := gtwAddr
	if gateway != "" && !strings.Contains(gateway, "/") {
		gateway = gateway + "/24"
	}

	networkName := fmt.Sprintf("%s-%s", srcName, dstName)

	// Add network connections to both devices
	var gatewayPtr *string
	if gateway != "" {
		gatewayPtr = &gateway
	}

	dstDevice.AddNetworkConnection(networkName, dm.MacGen(), gatewayPtr)
	srcDevice.AddNetworkConnection(networkName, dm.MacGen(), gatewayPtr)

	// Generate network XML template (similar to Python version)
	tmpl, err := template.ParseFS(templates, "templates/virt_network.xml.tmpl")
	if err != nil {
		return fmt.Errorf("failed to parse network template: %w", err)
	}

	data := NetworkTemplateData{
		Name: networkName,
		Brdg: fmt.Sprintf("%s-br", networkName),
	}

	if err := os.MkdirAll(dm.Context, 0755); err != nil {
		return fmt.Errorf("failed to create context directory: %w", err)
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to execute network template: %w", err)
	}

	networkXMLPath := filepath.Join(dm.Context, fmt.Sprintf("%s.xml", networkName))
	if err := os.WriteFile(networkXMLPath, []byte(buf.String()), 0644); err != nil {
		return fmt.Errorf("failed to write network XML: %w", err)
	}

	return nil
}