package main

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// NewDeviceManager creates a new DeviceManager instance
func NewDeviceManager() *DeviceManager {
	return &DeviceManager{
		Devices: make(map[string]*Device),
		MacIter: NewMacGenerator(0x22, 0x100),
	}
}

// NewMacGenerator creates a new MAC address generator
func NewMacGenerator(start, wrapAt int) *MacGenerator {
	return &MacGenerator{
		current: start,
		start:   start,
		wrapAt:  wrapAt,
	}
}

// Next generates the next MAC address
func (mg *MacGenerator) Next() string {
	xx := fmt.Sprintf("%02x", mg.current)
	mac := fmt.Sprintf("52:54:00:12:%s:01", xx)
	mg.current = (mg.current + 1) % mg.wrapAt
	return mac
}

// NewIPGenerator creates a new IP address generator
func NewIPGenerator(network *net.IPNet) *IPGenerator {
	// Get the first host IP in the network
	ip := make(net.IP, len(network.IP))
	copy(ip, network.IP)
	
	// Increment to get first host
	for i := len(ip) - 1; i >= 0; i-- {
		ip[i]++
		if ip[i] != 0 {
			break
		}
	}
	
	return &IPGenerator{
		network: network,
		current: ip,
	}
}

// Next generates the next IP address
func (ig *IPGenerator) Next() net.IP {
	result := make(net.IP, len(ig.current))
	copy(result, ig.current)
	
	// Increment current IP
	for i := len(ig.current) - 1; i >= 0; i-- {
		ig.current[i]++
		if ig.current[i] != 0 {
			break
		}
	}
	
	return result
}

// NewDevice creates a new Device instance
func NewDevice(devType DeviceType, name, address string) *Device {
	// Ensure address has CIDR notation
	if !strings.Contains(address, "/") {
		address = address + "/24"
	}
	
	return &Device{
		DevType:    devType,
		Name:       name,
		Address:    address,
		Networks:   make([]NetworkConnection, 0),
		IfaceCount: 0,
	}
}

// AddNetworkConnection adds a network connection to the device
func (d *Device) AddNetworkConnection(networkName, mac string, gateway *string) {
	d.IfaceCount++
	iface := fmt.Sprintf("ens%d", d.IfaceCount+1)
	
	conn := NetworkConnection{
		Name:    networkName,
		Iface:   iface,
		SrcIP:   d.Address,
		MAC:     mac,
		Gateway: gateway,
	}
	
	d.Networks = append(d.Networks, conn)
}

// Parse parses the configuration file and creates devices
func (dm *DeviceManager) Parse(cfgFilePath string, write bool) error {
	// Read the YAML configuration file
	data, err := os.ReadFile(cfgFilePath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}
	
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse YAML: %w", err)
	}
	
	// Create file path context for later use
	dir := filepath.Dir(cfgFilePath)
	dm.Context = filepath.Join(dir, "tmp")
	
	// Parse network address
	_, network, err := net.ParseCIDR(config.Network.Address)
	if err != nil {
		return fmt.Errorf("failed to parse network address: %w", err)
	}
	dm.NetworkAddress = network
	dm.IP4Iter = NewIPGenerator(network)
	
	// Create RTU devices
	for i, rtu := range config.RTUs {
		name := rtu.Name
		if name == "" {
			name = fmt.Sprintf("rtu%d", i)
		}
		
		address := rtu.Address
		if address == "" {
			nextIP := dm.IP4Iter.Next()
			address = fmt.Sprintf("%s/24", nextIP.String())
		}
		
		dm.Devices[name] = NewDevice(DeviceTypeRTU, name, address)
	}
	
	// Create Switch devices
	for i, sw := range config.Switches {
		name := sw.Name
		if name == "" {
			name = fmt.Sprintf("switch%d", i)
		}
		
		address := sw.Address
		if address == "" {
			nextIP := dm.IP4Iter.Next()
			address = fmt.Sprintf("%s/24", nextIP.String())
		}
		
		dm.Devices[name] = NewDevice(DeviceTypeSW, name, address)
		
		// Create network connections for switch
		for _, conn := range sw.Connected {
			dm.createNetwork(name, conn.To, "")
		}
	}
	
	if write {
		dm.createDevices()
	}
	
	return nil
}

// createDevices creates device configurations and files
func (dm *DeviceManager) createDevices() {
	for _, device := range dm.Devices {
		// Store original device address
		devAddr := device.Address
		
		// Set temporary address for default network
		lastChar := device.Name[len(device.Name)-1:]
		if num, err := strconv.Atoi(lastChar); err == nil {
			device.Address = fmt.Sprintf("192.168.122.%d/24", num+1)
		} else {
			device.Address = "192.168.122.2/24"
		}
		
		// Add default network connection
		gateway := "192.168.122.1"
		device.AddNetworkConnection("default", dm.MacIter.Next(), &gateway)
		
		// Restore original address
		device.Address = devAddr
		
		// Create device directory
		deviceDir := filepath.Join(dm.Context, device.Name)
		os.MkdirAll(deviceDir, 0755)
		
		// Set image path (simplified - just set the path without creating actual files)
		imagePath := filepath.Join(deviceDir, fmt.Sprintf("debian-12-%s.qcow2", device.Name))
		device.ImagePath = &imagePath
		
		// Set seed paths (simplified)
		seedPath := filepath.Join(deviceDir, "seed.iso")
		userDataPath := filepath.Join(deviceDir, "user-data")
		cloudDataPath := filepath.Join(deviceDir, "meta-data")
		
		device.SeedPath = &seedPath
		device.UserDataPath = &userDataPath
		device.CloudDataPath = &cloudDataPath
		
		// Create config.xml file (simplified)
		configPath := filepath.Join(deviceDir, "config.xml")
		configContent := dm.generateLibvirtXML(device)
		os.WriteFile(configPath, []byte(configContent), 0644)
	}
}

// createNetwork creates a network connection between two devices
func (dm *DeviceManager) createNetwork(srcName, dstName, gtwAddr string) {
	dstDevice, dstExists := dm.Devices[dstName]
	srcDevice, srcExists := dm.Devices[srcName]
	
	if !dstExists || !srcExists {
		return
	}
	
	gateway := gtwAddr
	if gateway != "" && !strings.Contains(gateway, "/") {
		gateway = gateway + "/24"
	}
	
	networkName := fmt.Sprintf("%s-%s", srcName, dstName)
	
	// Add network connections to both devices
	dstDevice.AddNetworkConnection(networkName, dm.MacIter.Next(), nil)
	srcDevice.AddNetworkConnection(networkName, dm.MacIter.Next(), nil)
	
	// Create network XML file (simplified)
	os.MkdirAll(dm.Context, 0755)
	networkPath := filepath.Join(dm.Context, fmt.Sprintf("%s.xml", networkName))
	networkContent := dm.generateNetworkXML(networkName)
	os.WriteFile(networkPath, []byte(networkContent), 0644)
}

// generateLibvirtXML generates a simple libvirt XML configuration
func (dm *DeviceManager) generateLibvirtXML(device *Device) string {
	return fmt.Sprintf(`<domain type='kvm'>
  <name>%s</name>
  <memory unit='MiB'>512</memory>
  <vcpu>1</vcpu>
  <os>
    <type arch='x86_64' machine='pc'>hvm</type>
  </os>
  <devices>
    <disk type='file' device='disk'>
      <source file='%s'/>
      <target dev='vda' bus='virtio'/>
    </disk>
    <disk type='file' device='cdrom'>
      <source file='%s'/>
      <target dev='sda' bus='sata'/>
    </disk>
  </devices>
</domain>`, device.Name, *device.ImagePath, *device.SeedPath)
}

// generateNetworkXML generates a simple network XML configuration
func (dm *DeviceManager) generateNetworkXML(networkName string) string {
	return fmt.Sprintf(`<network>
  <name>%s</name>
  <bridge name='%s-br'/>
</network>`, networkName, networkName)
}