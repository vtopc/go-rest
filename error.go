package rest

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// ReqError request error.
type ReqError struct {
	// Req HTTP request
	Req *http.Request

	// Err is either an (*APIError) or error from (*http.Client).Do()
	Err error
}

// NewReqError builds ReqError
func NewReqError(req *http.Request, wrappedErr error) error {
	if req == nil {
		return errors.New("empty request")
	}

	return &ReqError{
		Req: req,
		Err: wrappedErr,
	}
}

// Error implements error interface
func (e *ReqError) Error() string {
	return fmt.Sprintf("request %s %s failed: %s", e.Req.Method, e.Req.URL, e.Err)
}

// Unwrap provides compatibility for Go 1.13+ error chains.
func (e *ReqError) Unwrap() error {
	return e.Err
}

// APIError REST API error
type APIError struct {
	ExpectedStatusCodes []int

	// Resp HTTP response
	Resp *http.Response

	// Err is either error body or error from io.Copy.
	// TODO: be more specific?
	Err error
}

// Error implements error interface
func (e *APIError) Error() string {
	codes := make([]string, 0, len(e.ExpectedStatusCodes))
	for _, code := range e.ExpectedStatusCodes {
		codes = append(codes, strconv.Itoa(code))
	}

	return fmt.Sprintf("wrong status code (%d not in [%s]): %s",
		e.Resp.StatusCode, strings.Join(codes, ", "), e.Err)
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
