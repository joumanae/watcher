package watcher

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/fatih/color"
)

type Checker struct {
	Checks map[string]string
	Output io.Writer
}

// NewChecker starts the program.
func NewChecker(path string) (*Checker, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	Checks := map[string]string{}

	for scanner.Scan() {
		line := scanner.Text()
		matches := regexp.MustCompile(`(https?://[^\s]+)\s+([^\r\n]+)`).FindStringSubmatch(line)
		if len(matches) == 3 {
			url := strings.TrimSpace(matches[2])
			keyword := strings.TrimSpace(matches[1])
			Checks[keyword] = url
		}

	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return &Checker{
		Checks: Checks,
		Output: os.Stdout,
	}, nil
}

// Fetch fetches the urls and verifies that a typed keyword is on a page.
func Fetch(url string, keyword string) (matched bool, err error) {

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

	c, err := NewChecker("Checks.txt")
	if err != nil {
		fmt.Printf("There was an issue with the file %v", err)
		os.Exit(1)
	}

	for url, keyword := range c.Checks {
		matched, err := Fetch(url, keyword)
		if err != nil {
			fmt.Fprintf(c.Output, "There was an issue fetching the url %s \n", err)
		}
		if !matched {
			fmt.Fprintf(c.Output, "[%s]: No additional information about %s is available on the page %s \n", color.RedString("CHECKED-NO INFO"), keyword, url)

		} else {
			fmt.Fprintf(c.Output, "[%s]: There is information about %s. on the page %s\n", color.GreenString("CHECKED"), keyword, url)
		}
	}
	return 0
}
