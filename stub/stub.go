package stub

import "net/http"

type SuccessfulGetter struct {
	StatusCode int
}

func (g *SuccessfulGetter) Get(_ string) (*http.Response, error) {
	resp := http.Response{
		StatusCode: g.StatusCode,
		Body:       http.NoBody,
	}
	return &resp, nil
}
