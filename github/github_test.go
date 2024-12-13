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
		username string
		want     bool
	}
	testCases := map[string]TestCase{
		"contains two consecutive hyphens": {"jub0bs--on-GitHub", false},
		"starts with a hyphen":             {"-jub0bs-on-GitHub", false},
		"ends with a hyphen":               {"jub0bs-on-GitHub-", false},
		"too short":                        {"ab", false},
		"too long":                         {strings.Repeat("a", 40), false},
		"contains invalid chars":           {"ju^b0bs", false},
		"all good":                         {"jub0bs", true},
	}
	var gh github.GitHub
	for desc, tc := range testCases {
		f := func(t *testing.T) {
			got := gh.IsValid(tc.username)
			if got != tc.want {
				t.Errorf("github.IsValid(%q): got %t; want %t", tc.username, got, tc.want)
			}
		}
		t.Run(desc, f)
	}
}

func TestIsAvailable200OK(t *testing.T) {
	gh := github.GitHub{
		Client: &stub.Client{StatusCode: http.StatusOK},
	}
	avail, err := gh.IsAvailable("whatever")
	if err != nil || avail {
		t.Errorf("got %t, %v; want false, nil", avail, err)
	}
}

func TestIsAvailable404NotFound(t *testing.T) {
	gh := github.GitHub{
		Client: &stub.Client{StatusCode: http.StatusNotFound},
	}
	avail, err := gh.IsAvailable("whatever")
	if err != nil || !avail {
		t.Errorf("got %t, %v; want true, nil", avail, err)
	}
}
