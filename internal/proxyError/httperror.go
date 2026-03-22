package proxyError

import (
	"errors"
	"fmt"
	"net/http"
)

/*
	TODO: http Error: Expand this struct to display customized http errors to user
*/
type HTTPError struct {
	Status  int
	Message string
	Err     error
}

func (e *HTTPError) Error() string {
	if e.Err == nil {
		return e.Message
	}
	return fmt.Sprintf("%s: %v", e.Message, e.Err)
}

func (e *HTTPError) Unwrap() error {
	return e.Err
}

type ErrorResponder interface {
	WriteError(w http.ResponseWriter, r *http.Request, err error)
}

type PlainTextErrorResponder struct{}

func (PlainTextErrorResponder) WriteError(w http.ResponseWriter, _ *http.Request, err error) {
	var httpErr *HTTPError
	if errors.As(err, &httpErr) {
		http.Error(w, httpErr.Message, httpErr.Status)
		return
	}

	http.Error(w, http.StatusText(http.StatusBadGateway), http.StatusBadGateway)
}