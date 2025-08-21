package api

import (
	"net/http"

	"github.com/ClownDUmpling/GoFinal/pkg/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	limit := 50
	search := r.URL.Query().Get("search")

	// Получаем задачи из базы
	tasks, err := db.Tasks(limit, search)
	if err != nil {
		writeJSON(w, map[string]string{"error": "Ошибка при получении задач"}, http.StatusInternalServerError)
		return
	}

	// Возвращаем ответ
	writeJSON(w, TasksResp{Tasks: tasks}, http.StatusOK)
}
