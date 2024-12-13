package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jub0bs/namecheck/bluesky"
	"github.com/jub0bs/namecheck/github"
)

type Checker interface {
	IsValid(string) bool
	IsAvailable(string) (bool, error)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stdout, "usage: %s <username>\n", os.Args[0])
		os.Exit(1)
	}
	username := os.Args[1]
	gh := github.GitHub{
		Client: http.DefaultClient,
	}
	var bs bluesky.Bluesky
	checkers := []Checker{&gh, &bs}
	for _, checker := range checkers {
		if !checker.IsValid(username) {
			fmt.Printf("%q is invalid on %s.\n", username, checker)
			return
		}
		avail, err := checker.IsAvailable(username)
		if err != nil {
			log.Fatal(err)
		}
		if !avail {
			fmt.Printf("%q is valid but unavailable on %s.\n", username, checker)
			return
		}
		fmt.Printf("%q is valid and available on %s.\n", username, checker)
	}
}
