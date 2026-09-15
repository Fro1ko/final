package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Fro1ko/final/pkg/db"
)

type taskResponse struct {
	ID    int64  `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

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

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, taskResponse{Error: "Не указан идентификатор"})
		return
	}

	task, err := db.GetTask(id)
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

type taskUpdateRequest struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var request taskUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeJSON(w, http.StatusBadRequest, taskResponse{Error: err.Error()})
		return
	}

	id, err := strconv.ParseInt(request.ID, 10, 64)
	if err != nil || id < 1 {
		writeJSON(w, http.StatusBadRequest, taskResponse{Error: "Некорректный идентификатор"})
		return
	}
	task := db.Task{
		ID:      id,
		Date:    request.Date,
		Title:   request.Title,
		Comment: request.Comment,
		Repeat:  request.Repeat,
	}
	if err := normalizeTask(&task); err != nil {
		writeJSON(w, http.StatusBadRequest, taskResponse{Error: err.Error()})
		return
	}
	if err := db.UpdateTask(&task); err != nil {
		writeJSON(w, http.StatusNotFound, taskResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, taskResponse{})
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, taskResponse{Error: "Не указан идентификатор"})
		return
	}

	if err := db.DeleteTask(id); err != nil {
		writeJSON(w, http.StatusNotFound, taskResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, taskResponse{})
}

func normalizeTask(task *db.Task) error {
	if task.Title == "" {
		return fmt.Errorf("Не указан заголовок задачи")
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

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		return
	}
}
