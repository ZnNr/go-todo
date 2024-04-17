package nextdate

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

// ParseRequestParams парсирует параметры запроса и возвращает временные метки и повторение.
func ParseRequestParams(r *http.Request) (time.Time, time.Time, string, error) {
	nowStr := r.URL.Query().Get("now")    // Получение параметра "now" из URL запроса
	dateStr := r.URL.Query().Get("date")  // Получение параметра "date" из URL запроса
	repeat := r.URL.Query().Get("repeat") // Получение параметра "repeat" из URL запроса

	// Парсинг строки "now" в формат времени
	nowTime, err := parseDate(nowStr, "parsing string now to date")
	if err != nil {
		return time.Time{}, time.Time{}, "", fmt.Errorf("Failed to parse 'now' date: %v", err)
	}
	// Парсинг строки "date" в формат времени
	dateTime, err := parseDate(dateStr, "parsing string date to date")
	if err != nil {
		return time.Time{}, time.Time{}, "", fmt.Errorf("Failed to parse 'date' date: %v", err)
	}
	// Проверка наличия параметра "repeat" в запросе
	if repeat == "" {
		return time.Time{}, time.Time{}, "", errors.New("empty repeat")
	}

	return nowTime, dateTime, repeat, nil
}

// NextDate обрабатывает запрос, вычисляет следующую дату и записывает результат в ответ.
func NextDate(w http.ResponseWriter, r *http.Request) {
	nowTime, dateTime, repeat, err := ParseRequestParams(r)
	if err != nil {
		handleError(w, err, "Request parameter error")
		return
	}
	// Вычисление следующей даты на основе полученных параметров
	nextDate, err := CalculateNextDate(dateTime, nowTime, repeat)
	if err != nil {
		handleError(w, err, "Failed to calculate next date")
		return
	}

	logAndWriteResult(w, nowTime.Format("20060102"), dateTime.Format("20060102"), repeat, nextDate)
}

// parseDate парсит строку даты в формат времени с заданной целью.
func parseDate(dateStr, errorMsg string) (time.Time, error) {
	dateTime, err := time.Parse("20060102", dateStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("Failed: %s - %v", errorMsg, err)
	}
	return dateTime, nil
}

// logAndWriteResult выполняет логирование результатов и запись их в ответ.
func logAndWriteResult(w http.ResponseWriter, nowStr, dateStr, repeat string, nextDate time.Time) {
	log.Println("[Info] FOR now =", nowStr, "date =", dateStr, "repeat =", repeat)
	log.Println("[Info] nextDate =", nextDate.Format("20060102"))

	w.Write([]byte(nextDate.Format("20060102")))
}

// handleError обрабатывает ошибку, логирует сообщение об ошибке и отправляет его в ответ.
func handleError(w http.ResponseWriter, err error, errorMsg string) {
	log.Println("[Error] " + errorMsg)
	http.Error(w, err.Error(), http.StatusBadRequest)
}
