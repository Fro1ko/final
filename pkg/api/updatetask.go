package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Fro1ko/final/pkg/db"
)

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
