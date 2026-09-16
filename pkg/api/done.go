package api

import (
	"net/http"
	"time"

	"github.com/Fro1ko/final/pkg/db"
)

func taskDoneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, taskResponse{Error: "метод не разрешен"})
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, taskResponse{Error: "не указан идентификатор"})
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

	if task.Repeat == "" {
		err = db.DeleteTask(taskID)
	} else {
		var next string
		next, err = NextDate(time.Now(), task.Date, task.Repeat)
		if err == nil {
			err = db.UpdateDate(next, taskID)
		}
	}
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, taskResponse{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, taskResponse{})
}
