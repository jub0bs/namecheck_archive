package github

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/jub0bs/namecheck"
)

type GitHub struct {
	Client namecheck.Doer
}

var re = regexp.MustCompile("^[-0-9A-Za-z]{3,39}$")

func containsNoIllegalPattern(username string) bool {
	return !strings.HasPrefix(username, "-") &&
		!strings.HasSuffix(username, "-") &&
		!strings.Contains(username, "--")
}

func looksGood(username string) bool {
	return re.MatchString(username)
}

func (*GitHub) IsValid(username string) bool {
	return containsNoIllegalPattern(username) &&
		looksGood(username)
}

func (gh *GitHub) IsAvailable(ctx context.Context, username string) (bool, error) {
	endpoint := fmt.Sprintf("http://github.com/%s", username)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return false, err
	}
	resp, err := gh.Client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case http.StatusOK:
		return false, nil
	case http.StatusNotFound:
		return true, nil
	default:
		return false, ctx.Err()
	}
}

func (*GitHub) String() string {
	return "GitHub"
}
