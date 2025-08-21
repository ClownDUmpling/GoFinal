package api

import (
	"log"
	"net/http"
)

func Init() {
	http.HandleFunc("/api/nextdate", nextDayHandler)
	http.HandleFunc("/api/task", auth(taskHandler))
	http.HandleFunc("/api/tasks", auth(tasksHandler))
	http.HandleFunc("/api/task/done", auth(doneHandler))

	http.HandleFunc("/api/signin", signInHandler)

	log.Printf("API endpoints registered")
}
