package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Fro1ko/final/pkg/db"
)

const defaultLimit = 50

type tasksResponse struct {
	Tasks []taskResponseItem `json:"tasks"`
}

type taskResponseItem struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, taskResponse{Error: "метод не разрешен"})
		return
	}

	search := r.URL.Query().Get("search")
	date := ""
	if parsed, err := time.Parse("02.01.2006", search); err == nil {
		date = parsed.Format(dateFormat)
	}

	tasks, err := db.SearchTasks(search, date, defaultLimit)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, taskResponse{Error: err.Error()})
		return
	}

	response := tasksResponse{Tasks: make([]taskResponseItem, 0, len(tasks))}
	for _, task := range tasks {
		response.Tasks = append(response.Tasks, taskResponseItem{
			ID:      strconv.FormatInt(task.ID, 10),
			Date:    task.Date,
			Title:   task.Title,
			Comment: task.Comment,
			Repeat:  task.Repeat,
		})
	}
	writeJSON(w, http.StatusOK, response)
}
