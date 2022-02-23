package namecheck

import (
	"context"
	"fmt"
	"net/http"
)

type Validator interface {
	IsValid(string) bool
}

type Availabler interface {
	IsAvailable(context.Context, string) (bool, error)
}

type Checker interface {
	Validator
	Availabler
	fmt.Stringer
}

type Doer interface {
	Do(*http.Request) (*http.Response, error)
}
