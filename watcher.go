package watcher

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type Check struct {
	url     string
	keyword string
}

// TODO: ask John why is this even a thing?
func NewCheck(url string, keyword string) *Check {
	return &Check{
		url:     url,
		keyword: keyword,
	}
}

// Start a list of urls.
func (c *Check) StartList() (string, error) {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	filePath := scanner.Text()
	scanner.Scan()
	content := scanner.Text()
	file, err := os.Create(filePath)
	if err != nil {
		return "no filepath", fmt.Errorf("there was an issue creating the file: %v", err)
	}
	defer file.Close()
	_, err = file.WriteString(content)
	if err != nil {
		return "no filepath", fmt.Errorf("there was an issue writing the file: %v", err)
	}
	return filePath, nil
}

// Read the file and create a map with urls and keywords.
func (c *Check) ReadFileToMap(filePath string) (map[string]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	checks := map[string]string{
		c.url:     c.url,
		c.keyword: c.keyword,
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		parts := strings.SplitAfterN(line, "keyword", 2)
		if len(parts) == 2 {
			url := strings.TrimSpace(parts[0])
			keyword := strings.TrimSpace(parts[1])
			checks[url] = keyword
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return checks, nil
}

// Fetch the url and verify if the keyword is on the page fetched.
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
