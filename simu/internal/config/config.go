package config

import (
	"fmt"
	"net"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Network      NetworkConfig  `yaml:"network"`
	RTUs         []DeviceConfig `yaml:"rtus"`
	Switches     []DeviceConfig `yaml:"switches"`
	BinaryPath   string         `yaml:"binary_path"`
	FrontendPath string         `yaml:"frontend_path"`
}

type NetworkConfig struct {
	Address string `yaml:"address"`
	CIDR    *net.IPNet
}

type DeviceConfig struct {
	Name      string       `yaml:"name"`
	Address   string       `yaml:"address"`
	Connected []Connection `yaml:"connected,omitempty"`
}

type Connection struct {
	To      string `yaml:"to"`
	Gateway string `yaml:"gateway,omitempty"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	_, network, err := net.ParseCIDR(config.Network.Address)
	if err != nil {
		return nil, fmt.Errorf("failed to parse network CIDR: %w", err)
	}
	config.Network.CIDR = network
	// @TODO: Validate the config ??
	return &config, nil
}
