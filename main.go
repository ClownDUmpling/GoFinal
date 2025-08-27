package main

import (
	"log"

	"github.com/ClownDUmpling/GoFinal/pkg/db"
	"github.com/ClownDUmpling/GoFinal/pkg/server"
)

func main() {
	// Инициализация БД
	if err := db.Init(); err != nil {
		log.Fatalf("Ошибка инициализации БД: %v", err)
	}
	defer db.Close()

	//Старт сервера
	if err := server.Run(); err != nil {
		log.Fatal("Ошибка при запуске сервера:", err)
	}
}
