package todo

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewTaskHandler(t *testing.T) {
	t.Run("content-type should be application/json", func(t *testing.T) {
		var fakeInsert = func(string) error { return nil }

		handler := NewTaskHandler(fakeInsert)
		req := httptest.NewRequest("GET", "http://example.com/foo", nil)
		w := httptest.NewRecorder()
		handler(w, req)

		resp := w.Result()

		want := "application/json"
		get := resp.Header.Get("Content-Type")

		if want != get {
			t.Errorf("wants Content-Type %q but get %q", want, get)
		}
	})
	t.Run("binding request with nil payload should get BadRequest status", func(t *testing.T) {
		var fakeInsert = func(string) error { return nil }

		handler := NewTaskHandler(fakeInsert)
		req := httptest.NewRequest("GET", "http://example.com/foo", nil)
		w := httptest.NewRecorder()
		handler(w, req)

		resp := w.Result()

		want := http.StatusBadRequest
		get := resp.StatusCode

		if want != get {
			t.Errorf("nil body wants %d status but get %d", want, get)
		}
	})
}
