package authorization

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/golang-jwt/jwt/v5"
)

var unauthorized = errors.New("authentication required")

type Password struct {
	Password string `json:"password"`
}

type SignService struct {
	initialPassHash string
	secretKey       []byte
}

func hash(s string) string {
	sha := sha256.Sum256([]byte(s))
	hash := hex.EncodeToString(sha[:])
	return hash
}

// InitSignService инициализирует SignService с начальным хешем пароля и секретным ключом.
func InitSignService(initialPass string, secretKey []byte) SignService {
	return SignService{
		initialPassHash: hash(initialPass),
		secretKey:       secretKey,
	}
}

// jwtToken генерирует JWT токен на основе начального хеша пароля.
func (service SignService) jwtToken() (string, error) {
	claims := jwt.MapClaims{
		"pass": service.initialPassHash,
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return jwtToken.SignedString(service.secretKey)
}

// Auth выполняет проверку JWT токена для аутентификации.
func (service SignService) Auth(token string) error {
	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return service.secretKey, nil
	})
	if err != nil {
		return err
	}

	if !jwtToken.Valid {
		return unauthorized
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return unauthorized
	}

	passHash, ok := claims["pass"].(string)
	if !ok {
		return unauthorized
	}
	if passHash != service.initialPassHash {
		return unauthorized
	}
	return nil
}

// Signin обрабатывает пароль для создания JWT токена.
func (service SignService) signIn(pass Password) (string, error) {
	// Проверяем, совпадает ли хеш введенного пароля с начальным хешем пароля.
	if service.initialPassHash == hash(pass.Password) {
		// Если хеши совпадают, создаем JWT токен.
		return service.jwtToken()
	}
	// Возвращаем ошибку "authentication required", если хеши не совпадают.
	return "", unauthorized
}
