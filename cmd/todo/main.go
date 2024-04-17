package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	webDir := "web"
	port := getPort()
	http.Handle("/", http.FileServer(http.Dir(webDir)))
	fmt.Printf("Server is listening on port %s...\n", port)
	http.ListenAndServe(":"+port, nil)
}

func getPort() string {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}
	return port
}
