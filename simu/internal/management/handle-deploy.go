package management

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"scada-simu/internal/config"
	"scada-simu/internal/device"
)

type DeploymentRequest struct {
	ConfigPath string `json:"config_path"`
	OutputDir  string `json:"output_dir"`
}

type DeploymentResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

func (s *ServiceServer) handleDeployStatus(w http.ResponseWriter, r *http.Request) {
    s.sendResponse(w, DeploymentResponse{
        Message: "Current deployment status",
        Status:  s.deploymentStatus,
    })
}

func (s *ServiceServer) runDeployment(cfg *config.Config, outputDir string) {
	slog.Info("Starting VM deployment", "outputDir", outputDir)
	manager := device.InitManager(cfg, outputDir)
	manager.Config = cfg
	manager.Deploy()
	manager.StartVMs()
	slog.Info("VM deployment completed", "outputDir", outputDir)
	s.deploymentStatus = "completed"
}

func (s *ServiceServer) handleDeploy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req DeploymentRequest
	if r.Header.Get("Content-Type") == "application/json" {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
	} else {
		req.ConfigPath = r.FormValue("config_path")
		req.OutputDir = r.FormValue("output_dir")
	}

	if req.ConfigPath == "" {
		s.sendResponse(w, DeploymentResponse{
			Message: "Configuration path is required",
			Status:  "error",
		})
		return
	}

	if req.OutputDir == "" {
		req.OutputDir = "tmp/"
	}

	configPath := s.Config.LocalPath + "/" + req.ConfigPath
	outputPath := s.Config.LocalPath + "/" + req.OutputDir
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		s.sendResponse(w, DeploymentResponse{
			Message: fmt.Sprintf("Configuration file not found: %s", configPath),
			Status:  "error",
		})
		return
	}
	slog.Info("Loading configuration", "configPath", configPath)

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		s.sendResponse(w, DeploymentResponse{
			Message: fmt.Sprintf("Failed to load configuration: %v", err),
			Status:  "error",
		})
		return
	}

	if err := os.MkdirAll(req.OutputDir, 0755); err != nil {
		s.sendResponse(w, DeploymentResponse{
			Message: fmt.Sprintf("Failed to create output directory: %v", err),
			Status:  "error",
		})
		return
	}

	go s.runDeployment(cfg, outputPath)
	s.deploymentStatus = "deploying"

	s.sendResponse(w, DeploymentResponse{
		Message: "Deployment started successfully",
		Status:  s.deploymentStatus,
	})
}
