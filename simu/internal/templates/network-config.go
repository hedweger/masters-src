package templates

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"text/template"
)

type NetworkConfigContext struct {
	Connections []NetworkContext
	DeviceType     string
}

type NetworkContext struct {
	Interface string
	SourceIP  string
	Gateway   string
}

func WriteNetworkConfig(fp string, context NetworkConfigContext) (string, error) {
	tmpl, err := template.ParseFS(templateFiles, "data/network-config.tmpl")
	if err != nil {
		return "", fmt.Errorf("failed to parse network-config template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, context); err != nil {
		return "", fmt.Errorf("failed to execute network-config template: %w", err)
	}

	resPath := fmt.Sprintf("%s/network-config", fp)
	if err := os.WriteFile(resPath, buf.Bytes(), 0644); err != nil {
		return "", fmt.Errorf("failed to write network-config: %w", err)
	}

	resPath, err = filepath.Abs(resPath)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path for network-config: %w", err) // is this hittable????
	}
	slog.Debug("Network config written", "path", resPath)
	return resPath, nil
}
