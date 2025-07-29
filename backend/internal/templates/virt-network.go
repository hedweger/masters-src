package templates

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

type VirtNetworkContext struct {
	Name string
	Bridge string
}

func WriteVirtNetwork(fp string, context VirtNetworkContext) (string, error) {
	os.Mkdir(fmt.Sprintf("%s/networks/", fp), 0755)
	tmpl, err := template.ParseFS(templateFiles, "data/virt-network.tmpl")
	if err != nil {
		return "", fmt.Errorf("Failed to parse virt-network template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, context); err != nil {
		return "", fmt.Errorf("Failed to execute virt-network template: %w", err)
	}

	resPath := fmt.Sprintf("%s/networks/%s.xml", fp, context.Name)
	if err := os.WriteFile(resPath, buf.Bytes(), 0644); err != nil {
		return "", fmt.Errorf("Failed to write virt-network config: %w", err)
	}

	resPath, err = filepath.Abs(resPath) 
	if err != nil {
		return "", fmt.Errorf("Failed to get absolute path for virt-network config: %w", err) // is this hittable????
	}
	log.Println("[INFO] Wrote virt-network config to", resPath)
	return resPath, nil
}
