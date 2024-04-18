package main

import (
	"github.com/ZnNr/go-todo/internal/settings"
	"net/http"
)

// FileServer обрабатывает запросы на статические файлы и отправляет их клиенту.
func FileServer(w http.ResponseWriter, r *http.Request) {
	handler := http.FileServer(http.Dir(settings.WebPath))
	handler.ServeHTTP(w, r)
}
