package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/jub0bs/namecheck"
	"github.com/jub0bs/namecheck/github"
	"github.com/jub0bs/namecheck/twitter"
)

type Result struct {
	Username  string
	Platform  string
	Valid     bool
	Available bool
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: namecheck [username]")
		os.Exit(1)
	}
	username := os.Args[1]
	tw := twitter.Twitter{
		Client: http.DefaultClient,
	}
	gh := github.GitHub{
		Client: http.DefaultClient,
	}
	var checkers []namecheck.Checker
	for i := 0; i < 20; i++ {
		checkers = append(checkers, &tw, &gh)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()
	resultCh := make(chan Result, len(checkers))
	errCh := make(chan error, len(checkers)-1)
	var wg sync.WaitGroup
	wg.Add(len(checkers))
	for _, checker := range checkers {
		go check(ctx, checker, username, &wg, resultCh, errCh)
	}
	go func() {
		wg.Wait()
		close(resultCh)
	}()
	results := make([]Result, 0, len(checkers))
	var done bool
	for !done {
		select {
		case err := <-errCh:
			fmt.Fprintln(os.Stderr, err)
			return
		case res, ok := <-resultCh:
			if !ok {
				done = true
				continue
			}
			results = append(results, res)
		}
	}
	fmt.Println(results)
}

func check(
	ctx context.Context,
	checker namecheck.Checker,
	username string,
	wg *sync.WaitGroup,
	resultCh chan<- Result,
	errCh chan<- error,
) {
	defer wg.Done()
	res := Result{
		Username: username,
		Platform: checker.String(),
	}
	res.Valid = checker.IsValid(username)
	if !res.Valid {
		resultCh <- res
		return
	}
	var err error
	res.Available, err = checker.IsAvailable(ctx, username)
	if err != nil {
		errCh <- err
		return
	}
	resultCh <- res
}
