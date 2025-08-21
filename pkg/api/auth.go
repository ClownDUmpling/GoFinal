package api

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		password := os.Getenv("TODO_PASSWORD")
		if len(password) == 0 {
			next(w, r)
			return
		}
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		// здесь код для валидации и проверки JWT-токена
		// ...
		valid := validateToken(cookie.Value, password)
		if !valid {
			// возвращаем ошибку авторизации 401
			http.Error(w, "Authentification required", http.StatusUnauthorized)
			return
		}

		next(w, r)
	})
}

func validateToken(tokenString, password string) bool {
	hashedPassword := getPasswordHash(password)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(hashedPassword), nil // Используем хеш пароля как секретный ключ
	})

	if err != nil || !token.Valid {
		return false
	}

	return true
}

// generateToken создает JWT токен с хешем пароля в payload
func generateToken(password string) (string, error) {
	// Создаем хеш пароля для payload
	hashedPassword := getPasswordHash(password)

	// Создаем claims с хешем пароля
	claims := jwt.MapClaims{
		"password_hash": hashedPassword,
	}

	// Создаем токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Подписываем токен с использованием хеша пароля как секретного ключа
	return token.SignedString([]byte(hashedPassword))
}

// getPasswordHash возвращает SHA256 хеш пароля
func getPasswordHash(password string) string {
	hash := sha256.Sum256([]byte(password))
	return fmt.Sprintf("%x", hash)
}
