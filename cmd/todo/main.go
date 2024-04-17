package main

import (
	"fmt"
	"github.com/ZnNr/go-todo/internal/database"
	"github.com/ZnNr/go-todo/internal/nextdate"
	"log"

	"net/http"
	"os"
)

const (
	webDir      = "web"
	defaultPort = "7540"
)

func main() {
	database.InitializeDatabase("scheduler.db")

	addr := ":" + getPort()

	if err := initServer(addr); err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}
}

func initServer(addr string) error {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(webDir)))
	mux.HandleFunc("/api/nextdate", nextdate.NextDate) // Используем HandleFunc вместо Handle

	server := &http.Server{Addr: addr, Handler: mux}

	log.Printf("Server is listening on port %s...\n", addr)
	if err := server.ListenAndServe(); err != nil {
		return fmt.Errorf("failed to start server: %v", err)
	}
	return nil
}

func getPort() string {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = defaultPort
	}
	return port
}
