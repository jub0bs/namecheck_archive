package github_test

import (
	"strings"
	"testing"

	"github.com/jub0bs/namecheck/github"
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
