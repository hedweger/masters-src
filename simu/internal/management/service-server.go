package management

import (
	"encoding/json"
	"log"
	"net/http"
	"scada-simu/internal/config"
	"scada-simu/internal/device"
)

type ServiceConfig struct {
	LocalPath string
}

type ServiceServer struct {
	DeviceManager *device.Manager
	HttpClient    *http.Client
	Config        *ServiceConfig
	fileCache     string
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
	case "/api/scd/upload":
		s.uploadScd(w, r)
	case "/api/scd/model":
		s.viewModel(w, r)
	case "/api/scd/view":
		s.viewScd(w, r)
	case "/api/deploy":
		s.handleDeploy(w, r)
	case "/api/status":
		s.handleStatus(w, r)
	default:
		log.Printf("[LOG] 404 Not Found: %s", r.URL.Path)
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
	index_path := s.Config.LocalPath + "/frontend/index.html"
	log.Println("[LOG] Serving dashboard from " + index_path)
	http.ServeFile(w, r, index_path)
}

func NewServiceServer(cfg ServiceConfig) *ServiceServer {
	return &ServiceServer{
		DeviceManager: nil,
		HttpClient:    &http.Client{},
		Config:        &cfg,
		fileCache:     "",
	}
}
