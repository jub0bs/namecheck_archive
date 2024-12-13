package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jub0bs/namecheck/bluesky"
	"github.com/jub0bs/namecheck/github"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stdout, "usage: %s <username>\n", os.Args[0])
		os.Exit(1)
	}
	username := os.Args[1]
	gh := github.GitHub{
		Client: http.DefaultClient,
	}
	if !gh.IsValid(username) {
		fmt.Printf("%q is invalid on GitHub.\n", username)
		return
	}
	avail, err := gh.IsAvailable(username)
	if err != nil {
		log.Fatal(err)
	}
	if !avail {
		fmt.Printf("%q is valid but unavailable on GitHub.\n", username)
		return
	}
	fmt.Printf("%q is valid and available on GitHub.\n", username)

	var bs bluesky.Bluesky
	if !bs.IsValid(username) {
		fmt.Printf("%q is invalid on Bluesky.\n", username)
		return
	}
	avail, err = bs.IsAvailable(username)
	if err != nil {
		log.Fatal(err)
	}
	if !avail {
		fmt.Printf("%q is valid but unavailable on Bluesky.\n", username)
		return
	}
	fmt.Printf("%q is valid and available on Bluesky.\n", username)
}
