package virt

import (
	"fmt"
	"os/exec"
)

func CreateQcow2Image(fp string, imagePath string, deviceName string) (string, error) {
	resultPath := fmt.Sprintf("%s/debian-12-%s.qcow2", imagePath, deviceName)
	cmd := exec.Command("cp", fp, resultPath)
	return resultPath, cmd.Run()
}
