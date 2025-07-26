package main

import (
	"ied-ns/devman"
	"log"
	"path/filepath"
)

func main() {
	configPath, err := filepath.Abs("./config.yaml")
	if err != nil {
		log.Fatalf("Failed to get absolute path for config: %v", err)
	}

	if err := run(configPath); err != nil {
		log.Fatalf("Application failed: %v", err)
	}
}

func run(config string) error {
	dman := devman.NewDeviceManager()
	return dman.Parse(config, true)
}