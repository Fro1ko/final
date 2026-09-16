package api

import (
	"net/http"
	"os"
)

var configuredPassword string

func Init() {
	configuredPassword = os.Getenv("TODO_PASSWORD")
	http.HandleFunc("/api/nextdate", nextDateHandler)
	http.HandleFunc("/api/signin", signinHandler)
	http.HandleFunc("/api/task", auth(taskHandler))
	http.HandleFunc("/api/task/done", auth(taskDoneHandler))
	http.HandleFunc("/api/tasks", auth(tasksHandler))
}
