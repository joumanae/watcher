package watcher

import (
	"bufio"
	"context"
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

	c, err := NewChecker(filename)
	if err != nil {
		fmt.Printf("There was an issue with the file %v", err)
		os.Exit(1)
	}
	s.C = *c
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.Handler)
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

func (s *ServerFile) Handler(w http.ResponseWriter, r *http.Request) {
	// Set the content type to HTML
	w.Header().Set("Content-Type", "text/html")

	// Concatenate HTML content for all checks without line breaks
	htmlContent := ""
	for _, check := range s.C.Checks {
		htmlContent += check.RecordResult()
	}

	// Write the concatenated HTML content to the response
	_, err := w.Write([]byte(htmlContent))
	if err != nil {

		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func (s *ServerFile) Shutdown() error {

	ctx := context.Background()
	err := s.Srv.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("the server did not shut down %v", err)
	}
	return nil
}

func (c *Checker) Check(path string) error {

	var cs Check
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Please create a file with the correct name.")
		return err
	}
	defer file.Close()
	data := make([]byte, 50)
	_, err = file.Read(data)
	if err != nil {
		fmt.Println("File empty, add data:", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {

		line := scanner.Text()
		matches := regexp.MustCompile(`(https?://[^\s]+)\s+([^\r\n]+)`).FindStringSubmatch(line)
		if len(matches) == 3 {
			cs.keyword = strings.TrimSpace(matches[2])
			cs.url = strings.TrimSpace(matches[1])
			c.Checks = append(c.Checks, Check{
				keyword: cs.keyword,
				url:     cs.url,
			})
		}
		for _, check := range c.Checks {
			check.RecordResult()
			fmt.Println(check.RecordResult())
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func (c *Check) RecordResult() string {
	var s State
	fmt.Println("In record result")
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
	fmt.Println("in Match")
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
	// Check that the file exists
	f := "checks.txt"
	_, err := os.Stat(f)
	if os.IsNotExist(err) {
		fmt.Println("checks does not exit. The program will create it for you. Please add data to your file.")
		_, err := os.Create("checks.txt")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Print("File created")
		os.Exit(1)
	} else {
		fmt.Println("File exist, moving on to the next phase.")
	}
	var c Checker
	c.Check("checks.txt")

	//Start the server

	s.StartServerFile(":8080", "checks.txt")

	return 0
}
