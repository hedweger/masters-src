package managment

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
	binaryPath    string
}

func (s *ServiceServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		http.Redirect(w, r, "/dashboard/", http.StatusFound)
	case "/dashboard/":
		s.serverDashboard(w, r)
	case "/rtu-bin/":
		s.serveBinaries(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (s *ServiceServer) serverDashboard(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "dashboard.html")
}

func (s *ServiceServer) serveBinaries(w http.ResponseWriter, r *http.Request) {
	binary := r.URL.Path[len("/rtu-bin/"):]
	if _, err := os.Stat(s.binaryPath + binary); err != nil {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, s.binaryPath+binary)
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
	}
}
