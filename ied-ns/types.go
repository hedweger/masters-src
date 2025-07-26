package main

import (
	"net"
)

// DeviceType represents the type of device (equivalent to Python DeviceType enum)
type DeviceType string

const (
	DeviceTypeSW  DeviceType = "switch"
	DeviceTypeRTU DeviceType = "rtu"
)

// NetworkConnection represents a network connection for a device
type NetworkConnection struct {
	Name    string  `json:"name"`
	Iface   string  `json:"iface"`
	SrcIP   string  `json:"src_ip"`
	MAC     string  `json:"mac"`
	Gateway *string `json:"gateway,omitempty"`
}

// Device represents a network device (RTU or Switch)
type Device struct {
	DevType      DeviceType           `json:"dev_type"`
	Name         string               `json:"name"`
	Address      string               `json:"address"`
	Networks     []NetworkConnection  `json:"networks"`
	IfaceCount   int                  `json:"iface_count"`
	ImagePath    *string              `json:"image_path,omitempty"`
	SeedPath     *string              `json:"seed_path,omitempty"`
	UserDataPath *string              `json:"user_data_path,omitempty"`
	CloudDataPath *string             `json:"cloud_data_path,omitempty"`
}

// DeviceManager manages network devices and their configuration
type DeviceManager struct {
	Devices        map[string]*Device `json:"devices"`
	Context        string             `json:"context"`
	NetworkAddress *net.IPNet         `json:"network_address"`
	MacIter        *MacGenerator      `json:"-"`
	IP4Iter        *IPGenerator       `json:"-"`
}

// MacGenerator generates MAC addresses
type MacGenerator struct {
	current int
	start   int
	wrapAt  int
}

// IPGenerator generates IP addresses from a network
type IPGenerator struct {
	network *net.IPNet
	current net.IP
}

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