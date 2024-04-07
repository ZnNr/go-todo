package main

import (
	"net/http"
	"os"
)

func main() {
	port := getPort()
	webDir := "./web"

	http.Handle("/", http.FileServer(http.Dir(webDir)))

	http.ListenAndServe(":"+port, nil)
}

func getPort() string {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}
	return port
}
