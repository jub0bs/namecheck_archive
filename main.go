package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/jub0bs/namecheck/github"
	"github.com/jub0bs/namecheck/reddit"
)

type Validator interface {
	IsValid(string) bool
}

type Availabler interface {
	IsAvailable(string) (bool, error)
}

type Checker interface {
	Validator
	Availabler
	fmt.Stringer
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: namecheck <username>")
		os.Exit(1)
	}
	username := os.Args[1]
	gh := github.GitHub{Client: http.DefaultClient}
	rd := reddit.Reddit{Client: http.DefaultClient}
	var checkers []Checker
	for range 20 {
		checkers = append(checkers, &gh, &rd)
	}
	for _, checker := range checkers {
		valid := checker.IsValid(username)
		fmt.Printf("validity of %q on %s: %t\n", username, checker, valid)
		if !valid {
			continue
		}
		avail, err := checker.IsAvailable(username)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("availability of %q on %s: %t\n", username, checker, avail)
	}
}
