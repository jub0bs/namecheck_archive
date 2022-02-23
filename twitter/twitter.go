package twitter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/jub0bs/namecheck"
)

type Twitter struct {
	Client namecheck.Doer
}

var re = regexp.MustCompile("^[0-9A-Z_a-z]{4,15}$")

func containsNoIllegalPattern(username string) bool {
	return !strings.Contains(strings.ToLower(username), "twitter")
}

func looksGood(username string) bool {
	return re.MatchString(username)
}

func (*Twitter) IsValid(username string) bool {
	return containsNoIllegalPattern(username) &&
		looksGood(username)
}

func (tw *Twitter) IsAvailable(ctx context.Context, username string) (bool, error) {
	const tmpl = "https://europe-west6-namechecker-api.cloudfunctions.net/userlookup?username=%s"
	endpoint := fmt.Sprintf(tmpl, username)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return false, err
	}
	resp, err := tw.Client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return false, errors.New("unknown availability")
	}
	var dto struct {
		Data interface{} `json:"data"`
	}
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&dto); err != nil {
		return false, err
	}
	return dto.Data == nil, nil
}

func (*Twitter) String() string {
	return "Twitter"
}
