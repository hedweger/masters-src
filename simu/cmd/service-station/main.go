package main

import (
	"log"
	"net/http"
	"scada-simu/internal/management"
)


func main() {
	server := management.NewServiceServer(management.ServiceConfig{
		LocalPath: "/home/th/workspace/masters/simu",
	})
	log.Printf("[LOG] starting service server on http://localhost:8080")
	if err := http.ListenAndServe(":8080", server); err != nil {
		log.Fatalf("[ERROR] Failed to start server: %v", err)
	}
}
