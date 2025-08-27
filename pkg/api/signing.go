package api

import (
	"encoding/json"
	"net/http"
	"os"
	"time"
)

type SignInRequest struct {
	Password string `json:"password"`
}

type SignInResponse struct {
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}

func signInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SignInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, SignInResponse{Error: "Ошибка декодирования JSON"}, http.StatusBadRequest)
		return
	}

	// Получаем пароль из переменных окружения
	expectedPassword := os.Getenv("TODO_PASSWORD")

	// Если пароль не установлен, пропускаем аутентификацию
	if expectedPassword == "" {
		writeJSON(w, SignInResponse{Error: "Аутентификация не требуется"}, http.StatusBadRequest)
		return
	}

	// Проверяем пароль
	if req.Password != expectedPassword {
		writeJSON(w, SignInResponse{Error: "Неверный пароль"}, http.StatusUnauthorized)
		return
	}

	// Генерируем JWT токен
	token, err := generateToken(expectedPassword)
	if err != nil {
		writeJSON(w, SignInResponse{Error: "Ошибка генерации токена"}, http.StatusInternalServerError)
		return
	}

	// Устанавливаем куку
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(8 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	})

	writeJSON(w, SignInResponse{Token: token}, http.StatusOK)
}
