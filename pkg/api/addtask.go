package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ClownDUmpling/GoFinal/pkg/db"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req TaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, TaskResponse{Error: "Ошибка декодирования JSON"}, http.StatusBadRequest)
		return
	}

	// Валидация обязательного поля
	if req.Title == "" {
		writeJSON(w, TaskResponse{Error: "Не указан заголовок задачи"}, http.StatusBadRequest)
		return
	}

	// Создаем задачу
	task := &db.Task{
		Date:    req.Date,
		Title:   req.Title,
		Comment: req.Comment,
		Repeat:  req.Repeat,
	}

	// Проверяем и корректируем дату
	if err := checkDate(task); err != nil {
		writeJSON(w, TaskResponse{Error: err.Error()}, http.StatusBadRequest)
		return
	}

	// Добавляем задачу в базу
	id, err := db.AddTask(task)
	if err != nil {
		writeJSON(w, TaskResponse{Error: "Ошибка при добавлении задачи в базу"}, http.StatusInternalServerError)
		return
	}

	writeJSON(w, TaskResponse{ID: id}, http.StatusOK)
}

func writeJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// checkDate проверяет и корректирует дату задачи согласно рекомендациям
func checkDate(task *db.Task) error {
	now := time.Now()

	// Если дата пустая, используем сегодняшнюю
	if task.Date == "" {
		task.Date = now.Format(DateFormat)

		return nil
	}

	// Проверяем корректность формата даты
	t, err := time.Parse(DateFormat, task.Date)
	if err != nil {

		return err
	}

	// если дата в прошлом или сегодня
	if !afterNow(t, now) {
		if t.Format(DateFormat) == now.Format(DateFormat) {
			return nil
		}

		if task.Repeat == "" {
			// если правила повторения нет, то берём сегодняшнее число
			task.Date = now.Format(DateFormat)

		} else {
			// в противном случае, вычисляем следующую дату
			next, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return err
			}
			task.Date = next
		}
	}
	// если дата в будущем, оставляем как есть
	return nil
}
