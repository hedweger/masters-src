package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"scada-simu/internal/config"
	"scada-simu/internal/device"
)

func main() {
	configPath := flag.String("cfg", "config.yaml", "Path to configuration file")
	outputDir := flag.String("output", "tmp/", "Output directory for generated files")
	flag.Parse()

	cfg, err := processConfig(*configPath, *outputDir)
	if err != nil {
		log.Fatalf("[ERROR] %v", err)
	}

	os.Mkdir(*outputDir, 0755)
	manager := device.InitManager(cfg, *outputDir)
	manager.StartVMs()
}

func processConfig(configPath string, outputDir string) (*config.Config, error) {
	if configPath == "" {
		return nil, fmt.Errorf("Configuration path cannot be empty")
	}
	if outputDir == "" {
		log.Println("No output directory specified, using default: tmp/")
		outputDir = "tmp/"
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to load configuration: %v", err)
	}

	log.Printf("[LOG] Loaded configuration: %+v\n", cfg)
	return cfg, nil
}
