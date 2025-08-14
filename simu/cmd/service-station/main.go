package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"scada-simu/internal/management"
)


func main() {
	server := management.NewServiceServer(management.ServiceConfig{
		LocalPath: "/home/th/workspace/masters/simu",
	})
	os.RemoveAll(server.Config.LocalPath + "/tmp")
	slog.Info("starting service server on http://localhost:8080")
	if err := http.ListenAndServe(":8080", server); err != nil {
		panic(fmt.Sprintf("failed to start server: %v", err))
	}
}
