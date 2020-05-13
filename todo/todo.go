package todo

import (
	"database/sql"
	"log"
	"net/http"
	"withdom/jsonify"

	"github.com/pkg/errors"
)

type Task struct {
	ID    uint64 `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

type findAllFunc func() ([]Task, error)

func FindAllTask(db *sql.DB) findAllFunc {
	return func() ([]Task, error) {
		rows, err := db.Query("SELECT rowid, title, done FROM todos;")
		if err != nil {
			return nil, errors.Wrap(err, "query all todos")
		}

		defer rows.Close()

		tasks := []Task{}
		for rows.Next() {
			var task Task
			if err := rows.Scan(&task.ID, &task.Title, &task.Done); err != nil {
				log.Println(err)
			}
			tasks = append(tasks, task)
		}

		return tasks, nil
	}
}

func TaskHandler(all findAllFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tasks, err := all()
		if err != nil {
			jsonify.Json(w)(
				http.StatusInternalServerError,
				map[string]string{
					"error": err.Error(),
				})
			return
		}

		jsonify.Json(w)(
			http.StatusOK,
			map[string]interface{}{
				"todos": tasks,
			})
	}
}
