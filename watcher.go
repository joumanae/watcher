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
func (c *Check) ReadFileSaveInput(filePath string) (map[string]string, error) {
	inputData := make(map[string]string)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		input := scanner.Text()
		parts := strings.SplitN(input, " ", 2)
		c.url = parts[0]
		c.keyword = parts[1]
		inputData[c.url] = c.keyword
	}
	return inputData, nil
}

// Fetch the url and verify if the keyword is on the page fetched.
// https://innercitytennis.clubautomation.com/calendar/programs, 6U Red Rockets, https://www.americankaratestudio.com/stlouispark KIDS GREEN-PURPLE
// When I fetch this, it does not look at images
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

	var c Check
	NewCheck(c.keyword, c.url)
	filepath, err := c.StartList()
	if err != nil {
		return 1
	}

	MapedFile, err2 := c.ReadFileSaveInput(filepath)
	if err2 != nil {
		return 1
	}
	fmt.Println("Here is the map", MapedFile)

	for c.url, c.keyword = range MapedFile {
		fmt.Println(Fetch(c.url, c.keyword))
	}
	return 0
}
