package watcher

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

// 1- Fetch a page from the web
//2- Look for a specific word or phrase
//3- Share the result

func Fetch(url string, keyword string) (s string, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("the url was not fetched, %v", err)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("there was an issue reading the response %v", err)
	}
	sr := string(b)
	if strings.Contains(sr, keyword) {
		return fmt.Sprintf("The information about %v is available.", keyword), nil
	}
	return "Information not available", nil
}
