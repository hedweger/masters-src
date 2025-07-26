package cinit

import (
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed templates/*
var templates embed.FS

// NetworkConnection represents a network connection configuration
type NetworkConnection struct {
	Name    string
	Iface   string
	SrcIP   string
	MAC     string
	Gateway *string
}

// CINITResult contains the paths to generated cloud-init files
type CINITResult struct {
	ISOP     string
	UserData string
	CloudData string
}

func (c CINITResult) String() string {
	return fmt.Sprintf("ISO: %s, User Data: %s, Cloud Data: %s", c.ISOP, c.UserData, c.CloudData)
}

// FileWrite represents a file to be written during cloud-init
type FileWrite struct {
	Path        string
	Owner       string
	Permissions string
	Content     string
}

// UserData contains the data for generating user-data cloud-init file
type UserData struct {
	DevType  string
	Hostname string
	Password string
	Commands []string
	Writes   []FileWrite
}

// NetworkData contains the data for generating network-config cloud-init file
type NetworkData struct {
	Connections []NetworkConnection
	DevType     string
}

// Prepare generates cloud-init configuration files
func Prepare(devType, devname string, cmds []string, flws []FileWrite, connections []NetworkConnection, fp string, write bool) (*CINITResult, error) {
	// Parse templates
	userDataTmpl, err := template.ParseFS(templates, "templates/user-data.tmpl")
	if err != nil {
		return nil, fmt.Errorf("failed to parse user-data template: %w", err)
	}

	networkConfigTmpl, err := template.ParseFS(templates, "templates/network-config.tmpl")
	if err != nil {
		return nil, fmt.Errorf("failed to parse network-config template: %w", err)
	}

	// Prepare user data
	userData := UserData{
		DevType:  devType,
		Hostname: devname,
		Password: "root", // Note: same as Python version
		Commands: cmds,
		Writes:   flws,
	}

	// Prepare network data
	networkData := NetworkData{
		Connections: connections,
		DevType:     devType,
	}

	var userDataContent, networkConfigContent string

	// Execute user-data template
	var userDataBuf, networkConfigBuf strings.Builder
	if err := userDataTmpl.Execute(&userDataBuf, userData); err != nil {
		return nil, fmt.Errorf("failed to execute user-data template: %w", err)
	}
	userDataContent = userDataBuf.String()

	// Execute network-config template
	if err := networkConfigTmpl.Execute(&networkConfigBuf, networkData); err != nil {
		return nil, fmt.Errorf("failed to execute network-config template: %w", err)
	}
	networkConfigContent = networkConfigBuf.String()

	result := &CINITResult{
		ISOP:      filepath.Join(fp, "cloudinit.iso"),
		UserData:  filepath.Join(fp, "user-data"),
		CloudData: filepath.Join(fp, "cloud-data"),
	}

	if write {
		// Create directories
		if err := os.MkdirAll(fp, 0755); err != nil {
			return nil, fmt.Errorf("failed to create directory %s: %w", fp, err)
		}
		seedDir := filepath.Join(fp, "seed")
		if err := os.MkdirAll(seedDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create seed directory: %w", err)
		}

		// Write user-data
		userDataPath := filepath.Join(seedDir, "user-data")
		if err := os.WriteFile(userDataPath, []byte(userDataContent), 0644); err != nil {
			return nil, fmt.Errorf("failed to write user-data: %w", err)
		}

		// Write network-config
		networkConfigPath := filepath.Join(seedDir, "network-config")
		if err := os.WriteFile(networkConfigPath, []byte(networkConfigContent), 0644); err != nil {
			return nil, fmt.Errorf("failed to write network-config: %w", err)
		}

		// Write meta-data (empty)
		metaDataPath := filepath.Join(seedDir, "meta-data")
		if err := os.WriteFile(metaDataPath, []byte(""), 0644); err != nil {
			return nil, fmt.Errorf("failed to write meta-data: %w", err)
		}

		// Generate ISO
		isoPath, err := filepath.Abs(filepath.Join(fp, "cloudinit.iso"))
		if err != nil {
			return nil, fmt.Errorf("failed to get absolute path for ISO: %w", err)
		}

		cmd := exec.Command("genisoimage",
			"-o", isoPath,
			"-volid", "cidata",
			"-joliet",
			"-rock",
			userDataPath,
			metaDataPath,
			networkConfigPath,
		)
		if err := cmd.Run(); err != nil {
			return nil, fmt.Errorf("failed to generate ISO: %w", err)
		}
	}

	return result, nil
}