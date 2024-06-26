package watcher_test

import (
	"bytes"
	"cmp"
	"errors"
	"fmt"
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
	if err != nil {
		t.Error("there was an issue reading the response")
	}
	if s != true {
		t.Fatalf("the wrong text was fetched, here is the text that was fetched %v", s)
	}
	if !s {
		t.Fatalf("expected true, got false; the keyword was not found in the response.")
	}
}

func TestStartServerFile_DataPrinted(t *testing.T) {

	s := watcher.ServerFile{}
	go func() {
		var buf bytes.Buffer
		address := ":8080"
		filename := "checkstest.txt"
		err := s.StartServerFile(address, filename)
		expectedError := "The server did not start"
		if err != nil && err.Error() != expectedError {
			t.Errorf("Expected error: %s, got: %v", expectedError, err)
		}

		output := buf.String()
		expectedOutput := "Starting server on localhost:8080\n"
		if output != expectedOutput {
			t.Errorf("Unexpected output. Expected: %s, Got: %s", expectedOutput, output)
		}
	}()
	r := helperGet("8080")
	if r.StatusCode != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			r.StatusCode, http.StatusOK)
	}

}

//     err := s.StartServerFile("localhost:8080", "filename")

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

func TestHandlerReturnsErrorOpeningFile(t *testing.T) {
	t.Parallel()
	s := watcher.ServerFile{}
	w := httptest.NewRecorder()
	err := s.Handler(w, nil, "doesnotexist.txt")
	if err == nil {
		t.Fatal("expected error from non existent file")
	}
}

func TestCheck(t *testing.T) {
	t.Parallel()
	var c watcher.Checker
	want := 2
	_, err := c.Check("inexistantfile.txt")
	if err == nil {
		t.Error("expected the error that the file does not exist.")
	}
	checks, err := c.Check("checkstest.txt")
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}
	got := len(checks)
	if got != 2 && err != nil {
		t.Errorf(" An unexpected error occurred %v", err)
	}
	if cmp.Compare(want, got) != 0 {
		t.Errorf("Incorrect length. Want %v, got %v,%v", want, got, checks)
	}
}

func TestRecordResult(t *testing.T) {
	t.Parallel()
	var c watcher.Check

	_, err := c.Match("https://wizardzines.com/", "rr")

	record := c.RecordResult()
	if err != nil {
		t.Fatal(err)
	}
	if record != "<p><span style='color:red;'>[ERROR] </span> For keyword <span style='color:black;'></span></p>" {
		t.Errorf("there was an unexpected record %v", record)
	}

}

func TestChecker_Check_ErrorScanningSlice(t *testing.T) {

	tempFile, err := os.CreateTemp("", "example")
	if err != nil {
		t.Fatal(err)
	}

	_, err = tempFile.WriteString("invalid content")
	if err != nil {
		t.Fatal(err)
	}
	tempFile.Close()

	defer os.Remove(tempFile.Name())

	c := watcher.Checker{}

	result, err := c.Check(tempFile.Name())
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	if len(result) != 1 {
		t.Errorf("expected empty slice, got %v", result)
	}
}

func TestStartServerFile_FileNotFoundError(t *testing.T) {
	s := &watcher.ServerFile{}

	err := s.StartServerFile("localhost:8080", "/path/to/nonexistent/file")

	if err == nil || err.Error() != "issues with file info" {
		t.Errorf("Expected error: 'issues with file info', got: %v", err)
	}
}

func TestStartServerFile_NewCheckerError(t *testing.T) {

	expectedTasksError := errors.New("There were no files found")

	NewChecker := func(filename string) (*watcher.Checker, error) {
		return nil, expectedTasksError
	}
	c, err := NewChecker("emptytestfile.txt")

	if err != expectedTasksError {
		t.Errorf("expected error: %v, got: %v. checks length is %v", expectedTasksError, err, len(c.Checks))
	}
}

func TestStartServerFile_NoFilesFound(t *testing.T) {

	tempFile := createTempFile(t, "")
	defer os.Remove(tempFile)

	s := &watcher.ServerFile{}

	err := s.StartServerFile("localhost:8080", tempFile)

	if err == nil || err.Error() != "there were no files found" {
		t.Errorf("Expected the error there were no files found, got %v", err)
	}
}

// Helper function to create a temporary file for testing
func createTempFile(t *testing.T, content string) string {

	tempFile, err := os.CreateTemp("", "example")
	if err != nil {
		t.Fatalf("failed to create temporary file: %v", err)
	}
	defer tempFile.Close()

	_, err = tempFile.WriteString(content)
	if err != nil {
		t.Fatalf("failed to write to temporary file: %v", err)
	}

	return tempFile.Name()
}

func ExampleCheck() {
	var c watcher.Checker
	fmt.Println(c.Check("checkstest.txt"))
	// Output ["https://wizardzines.com/", "How DNS Works!", "https://betterexplained.com/cheatsheet/", "Intuitive Learning"]
}

func TestMain(m *testing.M) {
	os.Exit(testscript.RunMain(m, map[string]func() int{
		"watcher": watcher.Main,
	}))
}

func TestScript(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir: "testdata/script",
	})
}
