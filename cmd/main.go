package main

import (
	"github.com/joumanae/watcher"
)

func main() {
	// Call MagicFunction
	// get an output that says if there is available information per keyword
	// Access the file
	// Range over urls and keywords of the file
	// type check struct {url string, keyword string}
	var c watcher.Check
	filepath, err := c.StartList()
	if err != nil {
		return
	}

	MapedFile, err2 := c.ReadFileToMap(filepath)
	if err2 != nil {
		return
	}

	for url, keyword := range MapedFile {
		watcher.Fetch(url, keyword)
	}

	//TODO: once the list is started, how do I run the rest of the methods, what do I call for filepath
}
