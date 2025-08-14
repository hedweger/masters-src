package management

/* this package handles and server SCD features, enabling:
1. capability of uploading a SCD file
2. running an analyzer on that file (modelviewer.jar)
3. changing the configuration of the RTU (probably with some default devs?)
*/

import (
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"text/template"

	"github.com/google/uuid"
)

func (s *ServiceServer) viewScd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if s.fileCache == "" {
		http.Error(w, "No SCD configuration uploaded!", http.StatusNotFound)
		slog.Error("no scd file uploaded")
		return
	}

	file, err := os.ReadFile(s.fileCache)
	if err != nil {
		http.Error(w, "Error while handling SCD file", http.StatusInternalServerError)
		slog.Error("scd file read failed", "fileCache", s.fileCache, "error", err)
		return
	}

	escaped := template.HTMLEscapeString(string(file))
	slog.Info("serving cached file to frontend", "fileCache", s.fileCache)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(escaped))
}


func (s *ServiceServer) viewModel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if s.fileCache == "" {
		http.Error(w, "No SCD configuration uploaded!", http.StatusNotFound)
		slog.Error("no scd file uploaded", "fileCache", s.fileCache)
		return
	}

	genconf := s.Config.LocalPath + "/rtu/utils/genconfig.jar"
	out := s.Config.LocalPath + "/server-cache/model/current.cfg"
	cmd := exec.Command("java", "-jar", genconf, s.fileCache, out)
	err := cmd.Run()
	if err != nil {
		http.Error(w, "failed to generate config", http.StatusInternalServerError)
		slog.Error(fmt.Sprintf("failed to run %s: %v", cmd.String(), err), "error", err)
		return
	}
	result, err := os.ReadFile(out)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read generated model file: %v", err), http.StatusInternalServerError)
		slog.Error("model file read failed", "output", out, "error", err)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Write(result)
}

// @TODO
// No validation, html should be enough (for now).
// I assume only one user will be connected to the server at a single time,
// so we can use a single uuid to identify the file. I also assume that the
// server will not be storing these files permanently, which it probably will
// and this will needs to be refactored later.
func (s *ServiceServer) uploadScd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	file, _, err := r.FormFile("scdfile")
	if err != nil {
		http.Error(w, "Failed to get file from form", http.StatusBadRequest)
		return
	}
	defer file.Close()
	err = s.cacheFile(file)
	if err != nil {
		http.Error(w, "Failed to cache file, try again!", http.StatusInternalServerError)
		slog.Error("failed to cache file", "error", err)
	}
}

func (s *ServiceServer) cacheFile(file multipart.File) error {
	fileId := uuid.New().String()
	filePath := s.Config.LocalPath + "/server-cache/scd/" + fileId + ".scd"
	os.MkdirAll(s.Config.LocalPath+"/server-cache/scd", os.ModePerm)
	outfile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer outfile.Close()
	_, err = io.Copy(outfile, file)
	s.fileCache = filePath
	slog.Info("cached SCD file at", "fileId", fileId)
	return nil
}
