package rest

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

// ReqError request error
type ReqError struct {
	Method string
	URL    *url.URL
	Err    error // wrapped error
}

// NewReqError builds ReqError
func NewReqError(req *http.Request, wrappedErr error) error {
	if req == nil {
		return errors.New("empty request")
	}

	return &ReqError{
		Method: req.Method,
		URL:    req.URL,
		Err:    wrappedErr,
	}
}

// Error implements error interface
func (e *ReqError) Error() string {
	return fmt.Sprintf("request %s %s failed: %s", e.Method, e.URL, e.Err)
}

// Unwrap provides compatibility for Go 1.13+ error chains.
func (e *ReqError) Unwrap() error {
	return e.Err
}

// APIError REST API error
type APIError struct {
	ResponseStatusCode  int // HTTP status code
	ExpectedStatusCodes []int
	Err                 error // wrapped error
}

// Error implements error interface
func (e *APIError) Error() string {
	return fmt.Sprintf("wrong status code (%d not in %v): %s",
		e.ResponseStatusCode, e.ExpectedStatusCodes, e.Err)
}

// Unwrap provides compatibility for Go 1.13+ error chains.
func (e *APIError) Unwrap() error {
	return e.Err
}

// TODO: uncomment or remove:
// // StatusCodeFromAPIError returns HTTP status code based on err from the (Client).Do(...).
// // `err` shouldn't be <nil>.
// func StatusCodeFromAPIError(err error) int {
// 	if err == nil {
// 		// It's not clear which status code to return in this case: 200, 201 or 204.
// 		// TODO: should it panic instead?
// 		return 0
// 	}
//
// 	var apiErr *APIError
// 	if errors.As(err, &apiErr) {
// 		return apiErr.ResponseStatusCode
// 	}
//
// 	// not an APIError type:
// 	return -1
// }
