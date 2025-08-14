package management

import (
	"log/slog"
	"net/http"
	"os"
)

func (s *ServiceServer) serveClientBinary(w http.ResponseWriter, r *http.Request) {
	path := s.Config.LocalPath + "/frontend/ied-client"
	slog.Info("Serving client binary", "path", path)
	if _, err := os.Stat(path); err != nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Disposition", `attachment; filename="ied-client"`)
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(w, r, path)
}

func (s *ServiceServer) serveServerBinary(w http.ResponseWriter, r *http.Request) {
	path := s.Config.LocalPath + "/frontend/ied-server"
	slog.Info("Serving server binary", "path", path)
	if _, err := os.Stat(path); err != nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Disposition", `attachment; filename="ied-server"`)
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(w, r, path)
}
