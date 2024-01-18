package watcher_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joumanae/watcher"
)

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
	if s != "The information about Hello is available." {
		t.Fatalf("the wrong text was fetched, here is the text that was fetched %v", s)
	}
}
