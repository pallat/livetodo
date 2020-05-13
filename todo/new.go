package todo

import (
	"database/sql"
	"fmt"
	"net/http"
	"withdom/jsonify"

	"github.com/pkg/errors"
)

type insertTaskFunc func(string) error

func InsertTask(db *sql.DB, tbl string) insertTaskFunc {
	return func(title string) error {
		stmt := fmt.Sprintf("INSERT INTO %s (title, done) VALUES (?,?)", tbl)
		result, err := db.Exec(stmt, title, false)
		if err != nil {
			return errors.Wrap(err, "insert todos")
		}

		if n, err := result.RowsAffected(); err != nil {
			return errors.Wrap(err, "insert "+title)
			if n < 1 {
				return errors.New("can not insert " + title)
			}
		}

		return nil
	}
}

func NewTaskHandler(insert insertTaskFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var task Task

		if err := jsonify.Bind(r)(&task); err != nil {
			jsonify.Json(w)(
				http.StatusBadRequest,
				map[string]string{
					"error": err.Error(),
				})
			return
		}

		if err := insert(task.Title); err != nil {
			jsonify.Json(w)(
				http.StatusInternalServerError,
				map[string]string{
					"error": err.Error(),
				})
			return
		}

		jsonify.Json(w)(http.StatusOK, struct{}{})
	}
}
