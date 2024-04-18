package nextdate

import (
	"encoding/json"
	"github.com/ZnNr/go-todo/internal/settings"
	"net/http"
	"time"
)

// ErrorResponse представляет структуру ошибки для кодирования в JSON.
type ErrorResponse struct {
	Message string `json:"message"`
}

// writeJSONError отправляет ответ с ошибкой в формате JSON.
func writeJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	errResponse := ErrorResponse{
		Message: message,
	}
	json.NewEncoder(w).Encode(errResponse)
}

// GetNextDate обрабатывает HTTP запрос и возвращает следующую дату на основе входных параметров.
func GetNextDate(w http.ResponseWriter, r *http.Request) {
	now, err := time.Parse(settings.DateFormat, r.URL.Query().Get("now"))
	if err != nil {
		writeJSONError(w, "Invalid 'now' parameter", http.StatusBadRequest)
		return
	}

	date := r.URL.Query().Get("date")
	if len(date) == 0 {
		writeJSONError(w, "Invalid 'date' parameter", http.StatusBadRequest)
		return
	}

	repeat := r.URL.Query().Get("repeat")

	ans, err := NextDate(now, date, repeat)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(ans))
	if err != nil {
		writeJSONError(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}
