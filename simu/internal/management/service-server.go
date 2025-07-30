package management

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"scada-simu/internal/config"
	"scada-simu/internal/device"
)

type ServiceServer struct {
	deviceManager *device.Manager
	httpClient    *http.Client
	frontendPath  string
	binaryPath    string
}

func (s *ServiceServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		http.Redirect(w, r, "/dashboard/", http.StatusFound)
	case "/dashboard/":
		s.serveDashboard(w, r)
	case "/rtu-client/":
		s.serveClientBinary(w, r)
	case "/rtu-server/":
		s.serveServerBinary(w, r)
	case "/api/deploy":
		s.handleDeploy(w, r)
	case "/api/status":
		s.handleStatus(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (s *ServiceServer) runDeployment(cfg *config.Config, outputDir string) {
	log.Printf("[LOG] Starting VM deployment with config and output dir: %s", outputDir)
	manager := device.InitManager(cfg, outputDir)
	manager.Config = cfg
	manager.Deploy()
	manager.StartVMs()
	log.Printf("[LOG] VM deployment completed")
}

func (s *ServiceServer) sendResponse(w http.ResponseWriter, resp DeploymentResponse) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *ServiceServer) serveDashboard(w http.ResponseWriter, r *http.Request) {
	log.Println("[LOG] Serving dashboard from " + s.frontendPath + "index.html")
	http.ServeFile(w, r, s.frontendPath+"index.html")
}

func (s *ServiceServer) serveClientBinary(w http.ResponseWriter, r *http.Request) {
	path := s.binaryPath + "ied-client"
	log.Println("[LOG] Serving client binary from " + path)
	if _, err := os.Stat(path); err != nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Disposition", `attachment; filename="ied-client"`)
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(w, r, path)
}

func (s *ServiceServer) serveServerBinary(w http.ResponseWriter, r *http.Request) {
	path := s.binaryPath + "ied-server"
	log.Println("[LOG] Serving server binary from " + path)
	if _, err := os.Stat(path); err != nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Disposition", `attachment; filename="ied-server"`)
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(w, r, path)
}

// @TODO avoid hardcoding paths
func NewServiceServer() *ServiceServer {
	return &ServiceServer{
		deviceManager: nil,
		httpClient:    &http.Client{},
		binaryPath:    "/home/th/workspace/masters/simu/rtu-bin/",
		frontendPath:  "/home/th/workspace/masters/simu/cmd/service-station/frontend/",
	}
}
