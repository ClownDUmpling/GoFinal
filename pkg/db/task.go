package db

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// Добавляем задачу
func AddTask(task *Task) (int64, error) {

	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// Получаем задачи
func Tasks(limit int, search string) ([]*Task, error) {
	var tasks []*Task
	var err error

	if search == "" {
		// Без поиска - простой запрос
		tasks, err = getTasksWithoutSearch(limit)
	} else {
		// С поиском - определяем тип поиска
		tasks, err = getTasksWithSearch(limit, search)
	}

	if err != nil {
		return nil, err
	}

	// Гарантируем, что возвращаем не nil
	if tasks == nil {
		tasks = []*Task{}
	}

	return tasks, nil
}

// getTasksWithoutSearch - задачи без поиска
func getTasksWithoutSearch(limit int) ([]*Task, error) {
	query := `SELECT id, date, title, comment, repeat FROM scheduler 
	          ORDER BY date ASC, id ASC 
	          LIMIT ?`

	rows, err := db.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer rows.Close()

	return scanTasks(rows)
}

// getTasksWithSearch - задачи с поиском
func getTasksWithSearch(limit int, search string) ([]*Task, error) {
	// Проверяем, является ли поиск датой в формате 02.01.2006
	if isDateSearch(search) {
		return getTasksByDate(limit, search)
	}

	// Текстовый поиск
	return getTasksByText(limit, search)
}

// isDateSearch - проверяет, является ли строка датой в формате 02.01.2006
func isDateSearch(search string) bool {
	_, err := time.Parse("02.01.2006", search)
	return err == nil
}

// getTasksByDate - поиск по дате
func getTasksByDate(limit int, dateStr string) ([]*Task, error) {
	// Преобразуем дату из 02.01.2006 в 20060102
	parsedDate, err := time.Parse("02.01.2006", dateStr)
	if err != nil {
		return nil, fmt.Errorf("неверный формат даты: %w", err)
	}

	formattedDate := parsedDate.Format("20060102")

	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? ORDER BY date ASC, id ASC LIMIT ?`

	rows, err := db.Query(query, formattedDate, limit)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer rows.Close()

	return scanTasks(rows)
}

// getTasksByText - текстовый поиск
func getTasksByText(limit int, search string) ([]*Task, error) {
	// Экранируем специальные символы для LIKE
	searchPattern := "%" + strings.ReplaceAll(search, "%", "\\%") + "%"

	query := `SELECT id, date, title, comment, repeat FROM scheduler 
	          WHERE (title LIKE ? ESCAPE '\' OR comment LIKE ? ESCAPE '\')
	          ORDER BY date ASC, id ASC 
	          LIMIT ?`

	rows, err := db.Query(query, searchPattern, searchPattern, limit)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer rows.Close()

	return scanTasks(rows)
}

// scanTasks - сканирует результаты запроса
func scanTasks(rows *sql.Rows) ([]*Task, error) {
	var tasks []*Task

	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования: %w", err)
		}
		tasks = append(tasks, &task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при обработке результатов: %w", err)
	}

	return tasks, nil
}

// GetTask возвращает задачу по ID
func GetTask(id string) (*Task, error) {
	var task Task

	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	err := db.QueryRow(query, id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("задача не найдена")
		}
		return nil, fmt.Errorf("ошибка получения задачи: %w", err)
	}

	return &task, nil
}

// UpdateTask обновляет существующую задачу
func UpdateTask(task *Task) error {
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`

	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return fmt.Errorf("ошибка обновления задачи: %w", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка получения количества измененных строк: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}

	return nil
}

// DeleteTask удаляет задачу по ID
func DeleteTask(id string) error {
	query := `DELETE FROM scheduler WHERE id = ?`

	res, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("ошибка удаления задачи: %w", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка получения количества удаленных строк: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}

	return nil
}

// UpdateDate обновляет дату задачи
func UpdateDate(id string, newDate string) error {
	query := `UPDATE scheduler SET date = ? WHERE id = ?`

	res, err := db.Exec(query, newDate, id)
	if err != nil {
		return fmt.Errorf("ошибка обновления даты: %w", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка получения количества измененных строк: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}

	return nil
}
