package watcher_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joumanae/watcher"
)

func TestStartList(t *testing.T) {

}

func FuzzReadAndSaveInput(f *testing.F) {
	f.Fuzz(func(t *testing.T, input string) {
		var c watcher.Checker
		c.ReadFileSaveInput(input)
	})
}

func TestFetch(t *testing.T) {
	want := "Hello"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Write([]byte(want))
	}))
	defer server.Close()
	s, err := watcher.Fetch(server.URL, want)
	if err != nil {
		t.Error("the server failed")
	}
	if s != true {
		t.Fatalf("the wrong text was fetched, here is the text that was fetched %v", s)
	}
}
