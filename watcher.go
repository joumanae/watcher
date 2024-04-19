package watcher

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

const (
	Red   = "#ff0000"
	Blue  = "#0000ff"
	Green = "#00ff00 "
)

type Checker struct {
	Output io.Writer
	Checks []Check
}

type Check struct {
	url     string
	keyword string
	state   string
}

type ServerFile struct {
	Srv *http.Server
	C   Checker
}

func (s *ServerFile) StartServerFile(address, filename string) error {
	fmt.Printf("serving the file %v\n", filename)
	c, err := NewChecker(filename)
	if err != nil {
		fmt.Printf("There was an issue with the file %v", err)
		os.Exit(1)
	}

	fileInfo, err := os.Stat(filename)
	if err != nil {
		return errors.New("issues with file info")
	}

	if fileInfo.Size() == 0 {
		return errors.New("there were no checks found")
	}

	s.C = *c
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		s.Handler(w, r, filename)
	})

	s.Srv = &http.Server{
		Addr:    address,
		Handler: mux,
	}

	err = s.Srv.ListenAndServe()
	if err != nil {
		return fmt.Errorf("the server did not start %v", err)
	}

	return nil
}

func (s *ServerFile) Handler(w http.ResponseWriter, r *http.Request, filename string) {

	// Set the content type to HTML
	w.Header().Set("Content-Type", "text/html")

	// Start the HTML response
	htmlContent := "<html><head><title>Checker Results</title></head><body>"
	// Concatenate HTML content for all checks with proper HTML formatting
	var c Checker
	checks := c.Check(filename)

	for _, check := range checks {
		htmlContent += check.RecordResult()

	}
	// End the HTML response
	htmlContent += "</body></html>"

	// Write the HTML content to the response
	_, err := w.Write([]byte(htmlContent))
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)

	}

}

func (c *Checker) Check(path string) []Check {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("The file could not open")
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		matches := regexp.MustCompile(`(https?://[^\s]+)\s+([^\r\n]+)`).FindStringSubmatch(line)
		if len(matches) == 3 {
			keyword := strings.TrimSpace(matches[2])
			url := strings.TrimSpace(matches[1])
			c.Checks = append(c.Checks, Check{
				keyword: keyword,
				url:     url,
			})
		}

	}
	var s ServerFile
	s.C.Checks = c.Checks

	if err := scanner.Err(); err != nil {
		fmt.Println("Error scanning the slice")
		os.Exit(0)
	}

	return s.C.Checks
}

func (c *Check) RecordResult() string {
	var s State

	m, err := c.Match(c.url, c.keyword)

	if err != nil {

		s = StateError
		c.state = s.HtmlString()
		return fmt.Sprintf("<p><span style='color:red;'>[%s] </span> For keyword <span style='color:black;'>%s</span></p>",
			c.state,
			c.keyword,
		)
	}

	if m {
		s = StateFound
		c.state = s.HtmlString()
		return fmt.Sprintf("<p><span style='color:green;'>[%s] </span> For keyword <span style='color:black;'>%s</span></p>",
			c.state,
			c.keyword,
		)
	}

	s = StateChecked
	c.state = s.HtmlString()
	return fmt.Sprintf("<p><span style='color:blue;'>[%s] </span> For keyword <span style='color:black;'>%s</span></p>",
		c.state,
		c.keyword,
	)
}

func (s State) HtmlString() string {

	msg := string(s)
	switch s {
	case StateError:
		return msg
	case StateChecked:
		return msg
	case StateFound:
		return msg
	default:
		return msg
	}
}

// NewChecker starts the program.
func NewChecker(path string) (*Checker, error) {

	return &Checker{
		Checks: []Check{},
		Output: os.Stdout,
	}, nil
}

type State string

const (
	StateError   State = "ERROR"
	StateFound   State = "FOUND"
	StateChecked State = "CHECKED"
)

// Check just needs to check itself
// Fetch fetches the urls and verifies that a typed keyword is on a page.
func (c *Check) Match(url string, keyword string) (matched bool, err error) {

	resp, err := http.Get(string(url))
	if err != nil {
		return false, fmt.Errorf("the url was not fetched, %v", err)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("there was an issue reading the response %v", err)
	}
	sr := string(b)
	return strings.Contains(sr, keyword), nil
}

// Run the program
func Main() int {
	s := ServerFile{}
	f := "checks.txt"
	// Check that the file exists
	if len(os.Args) > 1 {
		f = os.Args[1]
	}

	_, err := os.Stat(f)
	if os.IsNotExist(err) {
		fmt.Println("checks does not exit. The program will create it for you.")
		_, err := os.Create(f)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("File created %v. Please add data.", f)
		os.Exit(0)
	} else {
		fmt.Println("File exist, moving on to the next phase.")
	}

	if err != nil {
		fmt.Println("Error checking file:", err)
		return 1 // Return error status
	}

	//Start the server
	err = s.StartServerFile(":8080", f)
	if err != nil {
		fmt.Printf("The server did not start %v", err)
		os.Exit(1)
	}
	return 0
}
