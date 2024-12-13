package stub

import "net/http"

type Client struct {
	StatusCode int
	Err        error
}

func (c *Client) Get(_ string) (*http.Response, error) {
	if c.Err != nil {
		return nil, c.Err
	}
	res := http.Response{
		StatusCode: c.StatusCode,
		Body:       http.NoBody,
	}
	return &res, nil
}
