package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Fro1ko/final/pkg/db"
)

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
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
	if err := db.DeleteTask(taskID); err != nil {
		writeJSON(w, http.StatusNotFound, taskResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, taskResponse{})
}

func parseTaskID(value string) (int64, error) {
	id, err := strconv.ParseInt(value, 10, 64)
	if err != nil || id < 1 {
		return 0, fmt.Errorf("Некорректный идентификатор")
	}
	return id, nil
}
