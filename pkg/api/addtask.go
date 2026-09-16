package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Fro1ko/final/pkg/db"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, http.StatusBadRequest, taskResponse{Error: err.Error()})
		return
	}
	if err := normalizeTask(&task); err != nil {
		writeJSON(w, http.StatusBadRequest, taskResponse{Error: err.Error()})
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, taskResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, taskResponse{ID: id})
}

func normalizeTask(task *db.Task) error {
	if task.Title == "" {
		return fmt.Errorf("не указан заголовок задачи")
	}

	now := time.Now()
	today := now.Format(dateFormat)
	if task.Date == "" {
		task.Date = today
	}

	date, err := time.Parse(dateFormat, task.Date)
	if err != nil {
		return err
	}

	var next string
	if task.Repeat != "" {
		next, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}
	if date.Format(dateFormat) < today {
		if task.Repeat == "" {
			task.Date = today
		} else {
			task.Date = next
		}
	}
	return nil
}
