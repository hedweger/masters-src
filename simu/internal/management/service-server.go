package management

import (
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
	default:
		http.NotFound(w, r)
	}
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

func NewServiceServer(cfg_path string, outputDir string) *ServiceServer {
	cfg, err := config.LoadConfig(cfg_path)
	if err != nil {
		log.Fatalf("[ERROR] Failed to load configuration: %s | %v", cfg_path, err)
	}

	return &ServiceServer{
		deviceManager: device.InitManager(cfg, outputDir),
		httpClient:    &http.Client{},
		binaryPath:    cfg.BinaryPath,
		frontendPath:  cfg.FrontendPath,
	}
}
