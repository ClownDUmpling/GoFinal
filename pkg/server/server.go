package server

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/ClownDUmpling/GoFinal/pkg/api"
)

const defaultPort = 7540

// Запуск сервера
func Run() error {
	// Получение порта
	port := getPort()
	// Инициализация до запуска
	api.Init()
	// Настройка обработчика для статических файлов
	fs := http.FileServer(http.Dir("./web"))
	http.Handle("/", fs)

	log.Printf("Сервер запущен на порту %d", port)
	return (http.ListenAndServe(":"+strconv.Itoa(port), nil))
}

func getPort() int {
	// Проверяем переменную окружения TODO_PORT
	if envPort := os.Getenv("TODO_PORT"); envPort != "" {
		if p, err := strconv.Atoi(envPort); err == nil {
			return p
		}
	}

	// Возвращаем порт по умолчанию
	return defaultPort
}
