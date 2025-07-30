package management

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"scada-simu/internal/config"
)

type DeploymentRequest struct {
	ConfigPath string `json:"config_path"`
	OutputDir  string `json:"output_dir"`
}

type DeploymentResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Status  string `json:"status"`
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
			Success: false,
			Message: "Configuration path is required",
			Status:  "error",
		})
		return
	}

	if req.OutputDir == "" {
		req.OutputDir = "tmp/"
	}

	if _, err := os.Stat(req.ConfigPath); os.IsNotExist(err) {
		s.sendResponse(w, DeploymentResponse{
			Success: false,
			Message: fmt.Sprintf("Configuration file not found: %s", req.ConfigPath),
			Status:  "error",
		})
		return
	}

	cfg, err := config.LoadConfig(req.ConfigPath)
	if err != nil {
		s.sendResponse(w, DeploymentResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to load configuration: %v", err),
			Status:  "error",
		})
		return
	}

	if err := os.MkdirAll(req.OutputDir, 0755); err != nil {
		s.sendResponse(w, DeploymentResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to create output directory: %v", err),
			Status:  "error",
		})
		return
	}

	go s.runDeployment(cfg, req.OutputDir)

	s.sendResponse(w, DeploymentResponse{
		Success: true,
		Message: "Deployment started successfully",
		Status:  "deploying",
	})
}
