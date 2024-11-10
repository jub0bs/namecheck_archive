package main

import (
	"context"
	"encoding/json"
	"log"
	"maps"
	"net/http"
	"sync"

	"github.com/jub0bs/cors"
	"github.com/jub0bs/namecheck"
	"github.com/jub0bs/namecheck/github"
	"github.com/jub0bs/namecheck/reddit"
	"golang.org/x/sync/errgroup"
)

type Result struct {
	Platform  string `json:"platform"`
	Valid     bool   `json:"valid"`
	Available bool   `json:"available"`
}

var (
	m  = make(map[string]uint)
	mu sync.Mutex
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /check", handleCheck)
	mux.HandleFunc("GET /stats", handleStats)

	// instantiate a CORS middleware whose config suits your needs
	corsMw, err := cors.NewMiddleware(cors.Config{
		Origins: []string{"https://jub0bs.github.io"},
		ExtraConfig: cors.ExtraConfig{
			PrivateNetworkAccess: true,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// apply the CORS middleware
	handler := corsMw.Wrap(mux)

	// start the server on port 8080; make sure to pass your custom handler
	log.Fatal(http.ListenAndServe(":8080", handler))
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
	gh := github.GitHub{Client: http.DefaultClient}
	rd := reddit.Reddit{Client: http.DefaultClient}
	var checkers []namecheck.Checker
	const n = 10
	for range n {
		checkers = append(checkers, &gh, &rd)
	}
	resultCh := make(chan Result)
	g, ctx := errgroup.WithContext(context.Background())
	var wg sync.WaitGroup
	for _, checker := range checkers {
		wg.Add(1)
		g.Go(func() error { return check(ctx, checker, username, &wg, resultCh) })
	}
	var err error
	go func() {
		err = g.Wait()
		close(resultCh)
	}()
	var results []Result
	for res := range resultCh {
		results = append(results, res)
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	data := struct { // note: anonymous struct type
		Username string   `json:"username"`
		Results  []Result `json:"results,omitempty"`
	}{
		Username: username,
		Results:  results,
	}
	enc := json.NewEncoder(w)
	if err := enc.Encode(data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func handleStats(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	mu.Lock()
	m := maps.Clone(m) // create a "deep" copy of the map
	mu.Unlock()
	enc := json.NewEncoder(w)
	if err := enc.Encode(m); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func check(
	ctx context.Context,
	checker namecheck.Checker,
	username string,
	wg *sync.WaitGroup,
	resultCh chan<- Result,
) error {
	defer wg.Done()
	res := Result{
		Platform: checker.String(),
		Valid:    checker.IsValid(username),
	}
	if !res.Valid {
		select {
		case <-ctx.Done():
		case resultCh <- res:
		}
		return nil
	}
	avail, err := checker.IsAvailable(username)
	if err != nil {
		return err
	}
	res.Available = avail
	select {
	case <-ctx.Done():
	case resultCh <- res:
	}
	return nil
}
