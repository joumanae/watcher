package main

import (
	"flag"
	"fmt"

	"github.com/joumanae/watcher"
)

func main() {
	//Runs Fetch, and provides a string ( probably will need a flag)
	url := flag.String("url", "", "Allows users to check the url")
	keyword := flag.String("keyword", "", "Allows users to check if one keyword is available")
	flag.Parse()
	fetchted, err := watcher.Fetch(*url, *keyword)
	if err != nil {
		panic(err)
	}
	fmt.Printf("This url was fetched %v", fetchted)
	// Later figure out the list
	// a server that runs all the time and alerts me then emails me

}
