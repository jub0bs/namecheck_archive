package github

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

var re = regexp.MustCompile("^[a-zA-Z0-9-]{3,39}$")

type GitHub struct {
	Client *http.Client
}

func (*GitHub) IsValid(username string) bool {
	return !strings.HasPrefix(username, "-") &&
		!strings.HasSuffix(username, "-") &&
		!strings.Contains(username, "--") &&
		re.MatchString(username)
}

func (gh *GitHub) IsAvailable(username string) (bool, error) {
	addr := "https://github.com/" + username
	res, err := gh.Client.Get(addr)
	if err != nil {
		return false, err
	}
	res.Body.Close()
	switch res.StatusCode {
	case http.StatusNotFound:
		return true, nil
	case http.StatusOK:
		return false, nil
	default:
		return false, fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
}

func (gh *GitHub) String() string {
	return "GitHub"
}
