package management

/* this package handles and server SCD features, enabling:
1. capability of uploading a SCD file
2. running an analyzer on that file (modelviewer.jar)
3. changing the configuration of the RTU (probably with some default devs?)
*/

import (
	"fmt"
	"io"
	"log"
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
		http.Error(w, "No SCD file uploaded", http.StatusNotFound)
		log.Printf("[ERROR] No SCD file uploaded")
		return
	}
	file, err := os.ReadFile(s.fileCache)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read cache at %s: %v", s.fileCache, err), http.StatusInternalServerError)
		log.Printf("[ERROR] Failed to read SCD file %s: %v", s.fileCache, err)
		return
	}
	escaped := template.HTMLEscapeString(string(file))
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to format SCD file: %v", err), http.StatusInternalServerError)
		log.Printf("[ERROR] Failed to format SCD file %s: %v", s.fileCache, err)
		return
	}
	log.Printf("[LOG] serving %s", s.fileCache)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(escaped))
}


func (s *ServiceServer) viewModel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if s.fileCache == "" {
		http.Error(w, "No SCD file uploaded", http.StatusNotFound)
		log.Printf("[ERROR] No SCD file uploaded")
		return
	}
	cmd := exec.Command("java", "-jar", s.Config.LocalPath+"/rtu/utils/modelviewer.jar", s.fileCache, "-s")
	result, err := cmd.Output()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to run modelviewer: %v", err), http.StatusInternalServerError)
		log.Printf("[ERROR] Failed to run %s: %v", cmd.String(), err)
		return
	}
	log.Printf("[LOG] %s ran successfully", cmd.String())
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
		http.Error(w, "Failed to cache file", http.StatusInternalServerError)
		log.Printf("[ERROR] Failed to cache SCD file: %v", err)
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
	log.Printf("[LOG] SCD file uploaded and saved as %s", filePath)
	return nil
}
