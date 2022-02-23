package stub

import "net/http"

type SuccessfulGetter struct {
	StatusCode int
}

func (sg *SuccessfulGetter) Do(_ *http.Request) (*http.Response, error) {
	res := http.Response{
		StatusCode: sg.StatusCode,
		Body:       http.NoBody,
	}
	return &res, nil
}
