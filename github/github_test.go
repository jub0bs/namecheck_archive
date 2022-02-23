package github_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/jub0bs/namecheck/github"
	"github.com/jub0bs/namecheck/stub"
)

func TestIsValid(t *testing.T) {
	cases := []struct {
		desc     string
		username string
		want     bool
	}{
		{desc: "too_short", username: "ab", want: false},
		{desc: "too_long", username: "1234567890123456789012345678901234567890", want: false},
		{desc: "contains_illegal_chars", username: "ju_bobs", want: false},
		{desc: "starts_with_hyphen", username: "-jubobs", want: false},
		{desc: "ends_with_hyphen", username: "jubobs-", want: false},
		{desc: "contains_hyphens_in_a_row", username: "ju--bobs", want: false},
		{desc: "valid", username: "jub0bs", want: true},
	}
	var gh github.GitHub
	for _, c := range cases {
		f := func(t *testing.T) {
			got := gh.IsValid(c.username)
			if got != c.want {
				t.Errorf("github.IsValid(%q): got %t; want %t", c.username, got, c.want)
			}
		}
		t.Run(c.desc, f)
	}
}

func TestIsAvailable(t *testing.T) {
	type TestCase struct {
		desc       string
		statusCode int
		want       bool
	}
	cases := []TestCase{
		{desc: "404", statusCode: http.StatusNotFound, want: true},
		{desc: "200", statusCode: http.StatusOK, want: false},
	}
	const username = "whatever"
	for _, c := range cases {
		gh := github.GitHub{
			Client: &stub.SuccessfulGetter{StatusCode: c.statusCode},
		}
		f := func(t *testing.T) {
			got, err := gh.IsAvailable(context.Background(), username)
			if err != nil || got != c.want {
				t.Errorf("github.IsAvailable(%q): got %t; want %t", username, got, c.want)
			}
		}
		t.Run(c.desc, f)
	}
}
