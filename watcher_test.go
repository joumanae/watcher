package watcher_test

import (
	"maps"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/joumanae/watcher"
)

func TestReadList(t *testing.T) {
	// some function that creates a checker from a config file
	// do you have the correct config
	checker, err := watcher.NewChecker("testdata/testfile.txt")
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]string{
		"https://innercitytennis.clubautomation.com/calendar/programs": "6U Rockets",
		"https://www.americankaratestudio.com/stlouispark":             "KIDS GREEN & UP",
	}
	got := checker.Checks
	if !maps.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
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
