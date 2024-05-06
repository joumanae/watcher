package watcher_test

import (
	"cmp"
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

func TestRecordResult(t *testing.T) {
	t.Parallel()
	var c watcher.Check
	var s watcher.State
	_, err := c.Match("https://wizardzines.com/", "rr")

	record := c.RecordResult()
	if err != nil && s != "StateError" {
		t.Errorf("the wrong error was given %v", err)
	}
	if record != "<p><span style='color:red;'>[ERROR] </span> For keyword <span style='color:black;'></span></p>" {
		t.Errorf("there was an unexpected record %v", record)
	}

}

func TestChecker_Check_ErrorScanningSlice(t *testing.T) {
	// Create a temporary file with contents that will cause an error during scanning
	tempFile, err := os.CreateTemp("", "example")
	if err != nil {
		t.Fatal(err)
	}
	// Write invalid content to the temporary file
	_, err = tempFile.WriteString("invalid content")
	if err != nil {
		t.Fatal(err)
	}
	tempFile.Close()
	// Remove the temporary file when the test finishes
	defer os.Remove(tempFile.Name())

	// Create a Checker instance
	c := watcher.Checker{}

	// Call the Check method with the path to the temporary file
	result := c.Check(tempFile.Name())

	// Verify that the result is empty
	if len(result) != 0 {
		t.Errorf("Expected empty result, got: %v", result)
	}
}

func TestStartServerFile_FileNotFoundError(t *testing.T) {
	s := &watcher.ServerFile{}

	err := s.StartServerFile("localhost:8080", "/path/to/nonexistent/file")

	if err == nil || err.Error() != "issues with file info" {
		t.Errorf("Expected error: 'issues with file info', got: %v", err)
	}
}

func TestStartServerFile_EmptyFileError(t *testing.T) {
	// Create a temporary empty file for testing
	tempFile := createTempFile(t, "")
	defer os.Remove(tempFile)

	s := &watcher.ServerFile{}

	err := s.StartServerFile("localhost:8080", tempFile)

	if err == nil || err.Error() != "there were no file found" {
		t.Errorf("Expected error: 'there were no file found', got: %v", err)
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
