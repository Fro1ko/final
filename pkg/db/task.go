package db

import (
	"database/sql"
	"errors"
	"fmt"
)

type Task struct {
	ID      int64  `db:"id" json:"id"`
	Date    string `db:"date" json:"date"`
	Title   string `db:"title" json:"title"`
	Comment string `db:"comment" json:"comment"`
	Repeat  string `db:"repeat" json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	result, err := database.Exec(
		`INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`,
		task.Date, task.Title, task.Comment, task.Repeat,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func GetTask(id int64) (*Task, error) {
	task := new(Task)
	err := database.QueryRow(
		`SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`,
		id,
	).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("task not found")
		}
		return nil, err
	}
	return task, nil
}

func UpdateTask(task *Task) error {
	result, err := database.Exec(
		`UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`,
		task.Date, task.Title, task.Comment, task.Repeat, task.ID,
	)
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}

func UpdateDate(next string, id int64) error {
	result, err := database.Exec(
		`UPDATE scheduler SET date = ? WHERE id = ?`,
		next, id,
	)
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}

func DeleteTask(id int64) error {
	result, err := database.Exec(`DELETE FROM scheduler WHERE id = ?`, id)
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}

func Tasks(limit int) ([]*Task, error) {
	return SearchTasks("", "", limit)
}

func SearchTasks(search string, date string, limit int) ([]*Task, error) {
	if limit < 1 {
		return nil, fmt.Errorf("task limit must be positive")
	}

	query := `SELECT id, date, title, comment, repeat FROM scheduler`
	args := make([]any, 0, 3)
	switch {
	case date != "":
		query += ` WHERE date = ?`
		args = append(args, date)
	case search != "":
		query += ` WHERE title LIKE ? OR comment LIKE ?`
		pattern := "%" + search + "%"
		args = append(args, pattern, pattern)
	}
	query += ` ORDER BY date ASC, id ASC LIMIT ?`
	args = append(args, limit)

	rows, err := database.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]*Task, 0)
	for rows.Next() {
		task := new(Task)
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}
