package templates

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

type VirtDeviceContext struct {
	Name     string
	RAM      string
	VCPU     int
	DiskPath string
	SeedPath string
	Networks []VirtDevNetworkContext
}

type VirtDevNetworkContext struct {
	Name string
	MAC  string
}

func WriteVirtDevice(fp string, context VirtDeviceContext) (string, error) {
	tmpl, err := template.ParseFS(templateFiles, "data/virt-device.tmpl")
	if err != nil {
		return "", fmt.Errorf("Failed to parse virt-device template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, context); err != nil {
		return "", fmt.Errorf("Failed to execute virt-device template: %w", err)
	}

	resPath := fmt.Sprintf("%s/%s.xml", fp, context.Name)
	if err := os.WriteFile(resPath, buf.Bytes(), 0644); err != nil {
		return "", fmt.Errorf("Failed to write virt-device config: %w", err)
	}

	resPath, err = filepath.Abs(resPath) 
	if err != nil {
		return "", fmt.Errorf("Failed to get absolute path for virt-device config: %w", err) // is this hittable????
	}
	log.Println("[INFO] Wrote virt-device config to", resPath)
	return resPath, nil
}
