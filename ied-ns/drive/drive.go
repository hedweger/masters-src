package drive

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// QCOW2 creates a QCOW2 disk image for the specified device
func QCOW2(context, name string, write bool) (string, error) {
	destPath := filepath.Join(context, fmt.Sprintf("debian-12-%s.qcow2", name))
	
	if write {
		if err := os.MkdirAll(context, 0755); err != nil {
			return "", fmt.Errorf("failed to create directory %s: %w", context, err)
		}
		
		srcPath, err := filepath.Abs("debian-12-genericcloud-amd64.qcow2")
		if err != nil {
			return "", fmt.Errorf("failed to get absolute path: %w", err)
		}
		
		cmd := exec.Command("cp", srcPath, destPath)
		if err := cmd.Run(); err != nil {
			return "", fmt.Errorf("failed to copy qcow2 file: %w", err)
		}
	}
	
	return destPath, nil
}