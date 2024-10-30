package github_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/jub0bs/namecheck/github"
	"github.com/jub0bs/namecheck/stub"
)

func TestIsValid(t *testing.T) {
	type TestCase struct {
		desc     string
		username string
		want     bool
	}
	testCases := []TestCase{
		{"contains two consecutive hyphens", "jub0bs--on-GitHub", false},
		{"too short", "ab", false},
		{"too long", strings.Repeat("a", 40), false},
		{"starts with a hyphen", "-jub0bs", false},
		{"ends with a hyphen", "jub0bs-", false},
		{"contains illegal chars", "jub^0bs", false},
		{"all good", "jub0bs", true},
	}
	var gh github.GitHub
	for _, tc := range testCases {
		f := func(t *testing.T) {
			got := gh.IsValid(tc.username)
			if got != tc.want {
				t.Errorf("github.IsValid(%q): got %t; want %t", tc.username, got, tc.want)
			}
		}
		t.Run(tc.desc, f)
	}
}

func TestIsAvailable200(t *testing.T) {
	gh := github.GitHub{
		Client: &stub.SuccessfulGetter{StatusCode: http.StatusOK},
	}
	const username = "whatever"
	avail, err := gh.IsAvailable(username)
	if err != nil || avail {
		const tmpl = "IsAvailable(%q): got %t, %v; want false, nil"
		t.Errorf(tmpl, username, avail, err)
	}
}
