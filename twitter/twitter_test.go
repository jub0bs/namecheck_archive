package twitter_test

import (
	"testing"

	"github.com/jub0bs/namecheck/twitter"
)

func TestIsValid(t *testing.T) {
	type TestCase struct {
		desc     string
		username string
		want     bool
	}
	cases := []TestCase{
		{desc: "contains_twitter", username: "jub0bsOnTwitter", want: false},
		{desc: "too_short", username: "foo", want: false},
		{desc: "too_long", username: "tooloonnnnnnnnnnnnnnng", want: false},
		{desc: "contains_illegal_chars", username: "ju-bobs", want: false},
		{desc: "valid", username: "jub0bs", want: true},
	}
	var tw twitter.Twitter
	for _, c := range cases {
		f := func(t *testing.T) {
			got := tw.IsValid(c.username)
			if got != c.want {
				t.Errorf("twitter.IsValid(%q): got %t; want %t", c.username, got, c.want)
			}
		}
		t.Run(c.desc, f)
	}
}
