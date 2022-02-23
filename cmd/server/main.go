package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/mux"
	"github.com/jub0bs/namecheck"
	"github.com/jub0bs/namecheck/github"
	"github.com/jub0bs/namecheck/twitter"
)

var (
	m  = make(map[string]uint)
	mu sync.Mutex
)

type Result struct {
	Username  string `json:"username"`
	Platform  string `json:"platform"`
	Valid     bool   `json:"valid"`
	Available bool   `json:"available"`
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/check", handleCheck)
	r.HandleFunc("/stats", handleStats)
	http.Handle("/", r)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func handleStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	mu.Lock()
	defer mu.Unlock()
	if err := enc.Encode(m); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func handleCheck(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	mu.Lock()
	m[username]++
	mu.Unlock()
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
	resultCh := make(chan Result, len(checkers))
	errCh := make(chan error, len(checkers)-1)
	var wg sync.WaitGroup
	wg.Add(len(checkers))
	for _, checker := range checkers {
		go check(r.Context(), checker, username, &wg, resultCh, errCh)
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
			w.WriteHeader(http.StatusInternalServerError)
			return
		case res, ok := <-resultCh:
			if !ok {
				done = true
				continue
			}
			results = append(results, res)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if err := enc.Encode(results); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
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
