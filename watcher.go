package watcher

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strings"

	"github.com/fatih/color"
)

type Checker struct {
	url     string
	keyword string
	Output  io.Writer
}

func NewChecker(url string, keyword string) *Checker {
	return &Checker{
		url:     url,
		keyword: keyword,
		Output:  os.Stdout,
	}
}

// Fetch the url and verify if the keyword is on the page fetched.
// https://innercitytennis.clubautomation.com/calendar/programs 6U Red Rockets https://www.americankaratestudio.com/stlouispark KIDS GREEN-PURPLE
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

func (c *Checker) PrintInformation(ctx context.Context) {

	matched, err := Fetch(c.url, c.keyword)
	if err != nil {
		fmt.Fprintf(c.Output, "There was an issue fetching the url %s \n", err)
	}
	if !matched {
		fmt.Fprintf(c.Output, "[%s]: No additional information about %s is available on the page %s \n", color.RedString("CHECKED-NO INFO"), c.keyword, c.url)

	} else {
		fmt.Fprintf(c.Output, "[%s]: There is information about %s. on the page %s\n", color.GreenString("CHECKED"), c.keyword, c.url)
	}
}

// Start a list of urls.
func (c *Checker) StartList() (string, error) {

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	filePath := scanner.Text()

	file, err := os.Create(filePath)
	if err != nil {
		return "no filepath", fmt.Errorf("there was an issue creating the file: %v", err)
	}
	defer file.Close()
	for scanner.Scan() {
		line := scanner.Text()
		_, err = file.WriteString(line + "\n ")
		if err != nil {
			return "no filepath", fmt.Errorf("there was an issue writing the file: %v", err)
		}
		if line == "exit" {
			break
		}

	}
	return filePath, nil
}

// Read the file and create a map with urls and keywords.
func (c *Checker) ReadFileSaveInput(filePath string) (map[string]string, error) {
	inputData := make(map[string]string)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		matches := regexp.MustCompile(`(https?://[^\s]+)\s+([^\r\n]+)`).FindStringSubmatch(line)
		if len(matches) == 3 {
			c.url = matches[2]
			c.keyword = matches[1]
			inputData[c.keyword] = c.url
		}

	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return inputData, nil
}

// Run the program
func Main() int {
	var c Checker

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	filepath, err := c.StartList()
	if err != nil {
		return 1
	}

	MapedFile, err2 := c.ReadFileSaveInput(filepath)
	if err2 != nil {
		return 1
	}

	resultCh := make(chan int)
	for url, keyword := range MapedFile {
		go func(u, k string) {
			defer func() {
				resultCh <- 0
			}()
			checker := NewChecker(u, k)
			checker.PrintInformation(ctx)
		}(url, keyword)
	}

	// Wait for all goroutines to complete
	for range MapedFile {
		<-resultCh
	}

	return 0
}
