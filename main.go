package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/labstack/gommon/log"
	_ "github.com/mattn/go-sqlite3"

	"withdom/todo"
)

func main() {
	r := mux.NewRouter()

	r.Use(MaskedMiddleware)

	r.HandleFunc("/versions", versionHandler)
	r.HandleFunc("/echo", echoHandler).Methods(http.MethodPost)

	db, err := newDB("./todos.db")
	if err != nil {
		log.Panicf("db file %s: %s", "todos.db", err)
	}
	defer db.Close()

	r.HandleFunc("/todos", todo.TaskHandler(todo.FindAllTask(db))).Methods(http.MethodGet)
	r.HandleFunc("/todos", todo.NewTaskHandler(todo.InsertTask(db, "todos"))).Methods(http.MethodPost)

	fmt.Println("serve on :1323")
	http.ListenAndServe(":1323", r)
}

func newDB(filedb string) (*sql.DB, error) {
	return sql.Open("sqlite3", filedb)
}

func jsonizer(r *http.Request) func(interface{}) error {
	return func(v interface{}) error {
		return json.NewDecoder(r.Body).Decode(v)
	}
}

func jsonify(w http.ResponseWriter) func(interface{}) error {
	w.Header().Set("Content-type", "application/json")
	return func(v interface{}) error {
		return json.NewEncoder(w).Encode(v)
	}
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	m := map[string]int{
		"verison": 1,
	}

	jsonify(w)(m)
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	var m map[string]interface{}

	jsonizer(r)(&m)
	m["id"] = "0123456789012"
	jsonify(w)(&m)
}

func MaskedMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mw := &maskedWriter{w}
		next.ServeHTTP(mw, r)
	})
}

type maskedWriter struct {
	http.ResponseWriter
}

func (w *maskedWriter) Write(b []byte) (int, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return 0, err
	}

	if v, ok := m["id"]; ok {
		if s, ok := v.(string); ok {
			m["id"] = "xxxxxxxxx" + s[9:]
		}
	}

	b, err := json.Marshal(&m)
	if err != nil {
		return 0, err
	}
	return w.ResponseWriter.Write(b)
}
