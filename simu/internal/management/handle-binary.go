package management

import (
	"log"
	"net/http"
	"os"
)

func (s *ServiceServer) serveClientBinary(w http.ResponseWriter, r *http.Request) {
	path := s.Config.LocalPath + "/frontend/ied-client"
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
	path := s.Config.LocalPath + "/frontend/ied-server"
	log.Println("[LOG] Serving server binary from " + path)
	if _, err := os.Stat(path); err != nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Disposition", `attachment; filename="ied-server"`)
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(w, r, path)
}
