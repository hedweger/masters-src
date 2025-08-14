package templates

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"text/template"
)

type UserDataContext struct {
	Hostname string
	Password string
	Commands []string
	Packages []string
}

func WriteUserData(fp string, context UserDataContext) (string, error) {
	tmpl, err := template.ParseFS(templateFiles, "data/user-data.tmpl")
	if err != nil {
		return "", fmt.Errorf("failed to parse user-data template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, context); err != nil {
		return "", fmt.Errorf("failed to execute user-data template: %w", err)
	}

	resPath := fmt.Sprintf("%s/user-data", fp)
	if err := os.WriteFile(resPath, buf.Bytes(), 0644); err != nil {
		return "", fmt.Errorf("failed to write user-data: %w", err)
	}

	resPath, err = filepath.Abs(resPath)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path for user-data: %w", err) // is this hittable????
	}
	slog.Info("User data config written", "path", resPath)
	return resPath, nil
}
