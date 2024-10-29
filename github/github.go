package github

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

type GitHub struct {
	Client *http.Client
}

var re = regexp.MustCompile(`^[a-zA-Z0-9-]{3,39}$`)

func (*GitHub) IsValid(username string) bool {
	return !strings.HasPrefix(username, "-") &&
		!strings.HasSuffix(username, "-") &&
		!strings.Contains(username, "--") &&
		re.MatchString(username)
}

func (gh *GitHub) IsAvailable(username string) (bool, error) {
	addr := fmt.Sprintf("https://github.com/%s", username)
	resp, err := gh.Client.Get(addr)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case http.StatusNotFound:
		return true, nil
	case http.StatusOK:
		return false, nil
	default:
		return false, errors.New("unknown availability")
	}
}
