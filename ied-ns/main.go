package main

import (
	"log"
	"path/filepath"
)

func main() {
	// Get absolute path to config.yaml (equivalent to Python's os.path.abspath("./config.yaml"))
	configPath, err := filepath.Abs("../config.yaml")
	if err != nil {
		log.Fatalf("Failed to get absolute path for config.yaml: %v", err)
	}
	
	// Create DeviceManager instance (equivalent to Python's dman = DeviceManager())
	dman := NewDeviceManager()
	
	// Parse configuration with write=true (equivalent to Python's dman.parse(config, write=True))
	if err := dman.Parse(configPath, true); err != nil {
		log.Fatalf("Failed to parse configuration: %v", err)
	}
	
	log.Printf("Successfully parsed configuration and created %d devices", len(dman.Devices))
}