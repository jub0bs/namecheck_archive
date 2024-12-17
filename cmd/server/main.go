package main

import (
	"context"
	"encoding/json"
	"log"
	"maps"
	"net/http"
	"sync"

	"github.com/jub0bs/cors"
	"github.com/jub0bs/namecheck/github"
)

type Checker interface {
	IsValid(string) bool
	IsAvailable(string) (bool, error)
	String() string
}

type Result struct {
	Platform  string `json:"platform"`
	Valid     bool   `json:"valid"`
	Available bool   `json:"available"`
}

var (
	stats = make(map[string]uint)
	mu    sync.Mutex
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /check", handleCheck)
	mux.HandleFunc("GET /stats", handleStats)
	corsMw, err := cors.NewMiddleware(cors.Config{
		Origins: []string{"https://jub0bs.github.io"},
		ExtraConfig: cors.ExtraConfig{
			PrivateNetworkAccess: true,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	if err := http.ListenAndServe(":8080", corsMw.Wrap(mux)); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func handleStats(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	mu.Lock()
	stats := maps.Clone(stats)
	mu.Unlock()
	enc := json.NewEncoder(w)
	if err := enc.Encode(stats); err != nil {
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
	stats[username]++
	mu.Unlock()
	gh := github.GitHub{
		Client: http.DefaultClient,
	}
	const n = 20
	checkers := make([]Checker, 0, n)
	for range n {
		checkers = append(checkers, &gh)
	}
	resultCh := make(chan Result)
	errorCh := make(chan error)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var wg sync.WaitGroup
	wg.Add(len(checkers))
	for _, checker := range checkers {
		go check(ctx, checker, username, &wg, resultCh, errorCh)
	}
	go func() {
		wg.Wait()
		close(resultCh)
	}()
	var results []Result
	var finished bool
	for !finished {
		select {
		case <-ctx.Done():
			w.WriteHeader(http.StatusInternalServerError)
			return
		case <-errorCh:
			cancel()
			w.WriteHeader(http.StatusInternalServerError)
			return
		case res, ok := <-resultCh:
			if !ok {
				finished = true
				continue
			}
			results = append(results, res)
		}
	}
	type RespBpody struct {
		Username string   `json:"username"`
		Results  []Result `json:"results,omitempty"`
	}
	data := RespBpody{
		Username: username,
		Results:  results,
	}
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if err := enc.Encode(data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func check(
	ctx context.Context,
	checker Checker,
	username string,
	wg *sync.WaitGroup,
	resultCh chan<- Result,
	errorCh chan<- error,
) {
	defer wg.Done()
	res := Result{
		Platform: checker.String(),
		Valid:    checker.IsValid(username),
	}
	if !res.Valid {
		send(ctx, resultCh, res)
		return
	}
	avail, err := checker.IsAvailable(username)
	if err != nil {
		send(ctx, errorCh, err)
		return
	}
	res.Available = avail
	send(ctx, resultCh, res)
}

func send[T any](ctx context.Context, ch chan<- T, v T) {
	select {
	case <-ctx.Done():
	case ch <- v:
	}
}
