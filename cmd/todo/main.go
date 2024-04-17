package main

import (
	"fmt"
	"github.com/ZnNr/go-todo/internal/database"
	"log"

	"net/http"
	"os"
)

func main() {
	database.InitializeDatabase()

	webDir := "web"
	port := getPort()
	addr := ":" + port

	http.Handle("/", http.FileServer(http.Dir(webDir)))

	server := &http.Server{Addr: addr}

	fmt.Printf("Server is listening on port %s...\n", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func getPort() string {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}
	return port
}
