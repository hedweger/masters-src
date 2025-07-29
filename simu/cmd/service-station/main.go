package main

import (
	"flag"
	"log"
	"net/http"
	"scada-simu/internal/management"
)

func main() {
	configPath := flag.String("cfg", "config.yaml", "Path to configuration file")
	outputDir := flag.String("output", "tmp/", "Output directory for generated files")
	flag.Parse()

	if configPath == nil || *configPath == "" {
		log.Fatalf("[ERROR] Configuration path cannot be empty")
	}
	if outputDir == nil || *outputDir == "" {
		log.Println("[LOG] No output directory specified, using default: tmp/")
		*outputDir = "tmp/"
	}

	server := management.NewServiceServer(*configPath, *outputDir)
	log.Printf("[LOG] starting service server on localhost:8080")
	if err := http.ListenAndServe(":8080", server); err != nil {
		log.Fatalf("[ERROR] Failed to start server: %v", err)
	}
}
