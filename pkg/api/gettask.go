package api

import (
	"net/http"
	"strconv"

	"github.com/Fro1ko/final/pkg/db"
)

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, taskResponse{Error: "Не указан идентификатор"})
		return
	}

	taskID, err := parseTaskID(id)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, taskResponse{Error: err.Error()})
		return
	}
	task, err := db.GetTask(taskID)
	if err != nil {
		writeJSON(w, http.StatusNotFound, taskResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, taskResponseItem{
		ID:      strconv.FormatInt(task.ID, 10),
		Date:    task.Date,
		Title:   task.Title,
		Comment: task.Comment,
		Repeat:  task.Repeat,
	})
}
