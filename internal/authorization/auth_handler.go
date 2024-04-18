package authorization

import (
	"encoding/json"
	"github.com/ZnNr/go-todo/internal/errorutil"
	"net/http"
)

var Service SignService

// AuthService описывает интерфейс сервиса аутентификации.
type AuthService interface {
	Auth(token string) error
	Signin(pass Password) (string, error)
}

// AuthMiddleware обеспечивает проверку аутентификации.
type AuthMiddleware struct {
	Service AuthService
}

// Auth проверяет аутентификацию и переходит к следующему обработчику в цепочке.
func (auth *AuthMiddleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		cookie, err := r.Cookie("token")
		if err != nil || cookie == nil {
			writeErrorAndRespond(w, http.StatusUnauthorized, unauthorized)
			return
		}

		err = auth.Service.Auth(cookie.Value)
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
	err := json.NewDecoder(r.Body).Decode(&pass)
	if err != nil {
		writeErrorAndRespond(w, http.StatusBadRequest, err)
		return
	}

	token, err := Service.Signin(pass)
	if err != nil {
		writeErrorAndRespond(w, http.StatusUnauthorized, err)
		return
	}

	ansBody, err := json.Marshal(map[string]string{"token": token})
	if err != nil {
		writeErrorAndRespond(w, http.StatusInternalServerError, err)
		return
	}

	w.Write(ansBody)
}

// writeErrorAndRespond пишет ошибку в ответ и устанавливает соответствующий код состояния.
func writeErrorAndRespond(w http.ResponseWriter, statusCode int, err error) {
	w.WriteHeader(statusCode)
	w.Write(errorutil.MarshalError(err))
}
