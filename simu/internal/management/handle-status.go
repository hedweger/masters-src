package management

import "net/http"

func (s *ServiceServer) handleStatus(w http.ResponseWriter, r *http.Request) {
	// this could return the current deployment status
	// For now, return a simple status
	s.sendResponse(w, DeploymentResponse{
		Success: true,
		Message: "Service is running",
		Status:  "ready",
	})
}
