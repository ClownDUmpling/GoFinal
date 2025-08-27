package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ClownDUmpling/GoFinal/pkg/db"
)

// GET /api/task?id=<id>
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, map[string]string{"error": "не указан идентификатор"}, http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()}, http.StatusNotFound)
		return
	}

	writeJSON(w, task, http.StatusOK)
}

// PUT /api/task
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, map[string]string{"error": "ошибка декодирования JSON"}, http.StatusBadRequest)
		return
	}

	// Проверяем обязательное поле
	if task.Title == "" {
		writeJSON(w, map[string]string{"error": "не указан заголовок задачи"}, http.StatusBadRequest)
		return
	}

	// Проверяем и корректируем дату
	if err := checkDate(&task); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	// Дополнительная проверка правила повторения
	if task.Repeat != "" {
		now := time.Now()
		_, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeJSON(w, map[string]string{"error": "Неверное правило повторения"}, http.StatusBadRequest)
			return
		}
	}

	// Обновляем задачу
	if err := db.UpdateTask(&task); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}

	// Успешное обновление - возвращаем пустой JSON
	writeJSON(w, map[string]interface{}{}, http.StatusOK)
}
