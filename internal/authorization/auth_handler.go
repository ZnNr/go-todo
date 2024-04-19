package authorization

import (
	"bytes"
	"encoding/json"
	"github.com/ZnNr/go-todo/internal/errorutil"
	"net/http"
)

var Service SignService

// Auth проверяет аутентификацию и переходит к следующему обработчику в цепочке.
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		cookie, err := r.Cookie("token")
		if err != nil {
			writeErrorAndRespond(w, http.StatusUnauthorized, err)
			return
		}
		if cookie == nil {
			writeErrorAndRespond(w, http.StatusUnauthorized, unauthorized)
			return
		}
		err = Service.Auth(cookie.Value)
		if err != nil {
			writeErrorAndRespond(w, http.StatusUnauthorized, err)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// PostPass обрабатывает запрос на создание токена после аутентификации.
func PostPass(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var pass Password
	buff := bytes.Buffer{}

	_, err := buff.ReadFrom(r.Body)
	if err != nil {
		writeErrorAndRespond(w, http.StatusUnauthorized, err)
		return
	}

	err = json.Unmarshal(buff.Bytes(), &pass)

	token, err := Service.Signin(pass)
	if err != nil {
		writeErrorAndRespond(w, http.StatusUnauthorized, err)
		return
	}

	ansBody, err := json.Marshal(
		struct {
			Token string `json:"token"`
		}{Token: token})

	w.Write(ansBody)
}

// writeErrorAndRespond пишет ошибку в ответ и устанавливает соответствующий код состояния.
func writeErrorAndRespond(w http.ResponseWriter, statusCode int, err error) {
	w.WriteHeader(statusCode)
	w.Write(errorutil.MarshalError(err))
}
