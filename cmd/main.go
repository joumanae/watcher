package main

import (
	"fmt"

	"github.com/joumanae/watcher"
)

func main() {
	// Call MagicFunction
	// get an output that says if there is available information per keyword
	// Access the file
	// Range over urls and keywords of the file
	// type check struct {url string, keyword string}

	checks := watcher.MagicFunction()
	for key, check := range checks {
		fmt.Println(watcher.Fetch(key, check))
	}
}
