package jsonify

import (
	"encoding/json"
	"net/http"
)

func Bind(r *http.Request) func(interface{}) error {
	return func(v interface{}) error {
		return json.NewDecoder(r.Body).Decode(v)
	}
}

func Json(w http.ResponseWriter) func(interface{}) error {
	w.Header().Set("Content-type", "application/json")
	return func(v interface{}) error {
		return json.NewEncoder(w).Encode(v)
	}
}
