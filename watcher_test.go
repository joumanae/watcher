package watcher_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joumanae/watcher"
)

func TestMatch(t *testing.T) {
	want := "Hello"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Write([]byte(want))
	}))
	defer server.Close()
	var c watcher.Check
	s, err := c.Match(server.URL, want)
	if err != nil {
		t.Error("the server failed")
	}
	if s != true {
		t.Fatalf("the wrong text was fetched, here is the text that was fetched %v", s)
	}
}

func TestStartServerFile(t *testing.T) {
	s := watcher.ServerFile{}
	go func() {
		address := ":8081"
		filename := "Checks.txt"
		err := s.StartServerFile(address, filename)
		if err != nil {
			panic(err)
		}
	}()
	r, err := http.Get("http://127.0.0.1:8081")
	if err != nil {
		t.Fatal(err)
	}

	if r.StatusCode != http.StatusOK {
		t.Fatalf("Exepected status %d, got %d", http.StatusOK, r.StatusCode)
	}
}

func TestThatHandlerServesHTML(t *testing.T) {
	s := watcher.ServerFile{}
	go func() {
		address := ":8081"
		filename := "Checks.txt"
		err := s.StartServerFile(address, filename)
		if err != nil {
			panic(err)
		}
	}()
	_, err := http.Get("http://127.0.0.1:8081")
	if err != nil {
		t.Fatal(err)
	}
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	s.Handler(w, r)

	expectedContentType := "text/html"
	if contentType := w.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("handler returned wrong content type: got %v want %v",
			contentType, expectedContentType)
	}
}

func TestCheck(t *testing.T) {
	var checker watcher.Checker
	checker.Check("Checks.txt")

}
