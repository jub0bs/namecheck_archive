package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

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
	const n = 20
	checkers := make([]Checker, 0, n)
	for range n {
		checkers = append(checkers, &gh)
	}
	var wg sync.WaitGroup
	wg.Add(len(checkers))
	for _, checker := range checkers {
		go check(checker, username, &wg)
	}
	wg.Wait()
}

func check(checker Checker, username string, wg *sync.WaitGroup) {
	defer wg.Done()
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
