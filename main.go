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
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: namecheck <username>")
		os.Exit(1)
	}
	username := os.Args[1]
	checkers := []Checker{
		&github.GitHub{Client: http.DefaultClient},
		&reddit.Reddit{Client: http.DefaultClient},
	}
	for _, checker := range checkers {
		valid := checker.IsValid(username)
		fmt.Printf("validity of %q on ???: %t\n", username, valid)
		if !valid {
			continue
		}
		avail, err := checker.IsAvailable(username)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("availability of %q on ???: %t\n", username, avail)
	}
}
