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

func Json(w http.ResponseWriter) func(int, interface{}) error {
	return func(code int, v interface{}) error {
		w.Header().Add("Content-type", "application/json")
		w.WriteHeader(code)
		return json.NewEncoder(w).Encode(v)
	}
}
