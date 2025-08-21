package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

var (
	db *sql.DB
	//Прописывание схемы создания таблицы
	schema = `
CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(255) NOT NULL,
    comment TEXT,
    repeat VARCHAR(128)
);
CREATE INDEX idx_date ON scheduler(date);`
)

// Init инициализирует БД и создаёт таблицы при необходимости
func Init() error {
	dbFile := getDBPath()
	// Проверка существования файла
	_, err := os.Stat(dbFile)
	install := os.IsNotExist(err)

	// Открытие БД
	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("ошибка открытия БД: %w", err)
	}

	// Создание таблицы при первом запуске
	if install {
		if _, err := db.Exec(schema); err != nil {
			return fmt.Errorf("ошибка создания таблицы: %w", err)
		}
	}

	// Проверка подключения
	return db.Ping()
}

func Close() error {
	if db != nil {
		return db.Close()
	}
	return nil
}

func getDBPath() string {
	if path := os.Getenv("TODO_DBFILE"); path != "" {
		return path
	}
	return "scheduler.db" // Значение по умолчанию
}
