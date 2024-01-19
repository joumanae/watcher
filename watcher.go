package watcher

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

// 1- Fetch a page from the web
//2- Look for a specific word or phrase
//3- Share the result via email

func Fetch(url string, keyword string) (matched bool, err error) {
	resp, err := http.Get(url)
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

func MagicFunction() map[string]string {

	checks := map[string]string{
		"https://innercitytennis.clubautomation.com/calendar/programs": "6U Red Rockets",
		"https://americankaratestudio.com/":                            "green",
	}

	return checks
}

// 1- Open the file
// 2 - Read it
// 3 - use the urls and keywords in the file to look up the useful information

// func EmailResult(message string) (s string, err error) {
// 	//https://mail.google.com/mail/u/0/#inbox
// 	emailConfig := gomail.NewMessage()
// 	emailConfig.SetHeader("From", "joumana.codes@gmail.com")
// 	emailConfig.SetHeader("To", "joumana.codes@gmail.com")
// 	emailConfig.SetHeader("Subject", "Summary of activity scanning")
// 	d := gomail.NewDialer("https://godoc.org/?q=smtp", 587, "joumana.codes@gmail.com", "123456")
// 	return "email sent", nil

// }
