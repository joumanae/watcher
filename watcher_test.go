package watcher_test

import (
	"cmp"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/joumanae/watcher"
	"github.com/rogpeppe/go-internal/testscript"
)

func TestMatch(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
	s := watcher.ServerFile{}
	go func() {
		address := ":8080"
		filename := "Checks.txt"
		err := s.StartServerFile(address, filename)
		if err != nil {
			panic(err)
		}
	}()
	r := helperGet("8080")

	if r.StatusCode != http.StatusOK {
		t.Fatalf("Exepected status %d, got %d", http.StatusOK, r.StatusCode)
	}
}

func helperGet(port string) *http.Response {
	r, err := http.Get("http://127.0.0.1:" + port)
	for err != nil {
		time.Sleep(time.Millisecond * 10)
		r, err = http.Get("http://127.0.0.1:" + port)
	}
	return r
}

func TestThatHandlerServesHTML(t *testing.T) {
	t.Parallel()
	s := watcher.ServerFile{}
	go func() {
		address := ":8081"
		filename := "Checks.txt"
		err := s.StartServerFile(address, filename)
		if err != nil {
			panic(err)
		}
	}()
	r := helperGet("8081")

	expectedContentType := "text/html"
	if contentType := r.Header.Get("content-type"); contentType != expectedContentType {
		t.Errorf("handler returned wrong content type: got %v want %v",
			contentType, expectedContentType)
	}
}

func TestCheck(t *testing.T) {
	t.Parallel()
	var c watcher.Checker
	want := 2
	got := len(c.Check("checkstest.txt"))
	if cmp.Compare(want, got) != 0 {
		t.Errorf("Incorrect length. Want %v, got %v,", want, got)
	}
}

func TestScript(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir: "testdata/script",
	})
}

func TestMain(m *testing.M) {
	os.Exit(testscript.RunMain(m, map[string]func() int{
		"watcher": watcher.Main,
	}))
}
