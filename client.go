package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/vtopc/go-rest/defaults"
)

// HTTPClient a HTTP client
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client the REST API client
type Client struct {
	httpClient HTTPClient
}

// NewClient returns REST API Client.
// Use https://github.com/cristalhq/hedgedhttp for retries.
// Check https://github.com/vtopc/restclient/tree/master/interceptors for middlewares/interceptors.
// TODO: switch to HTTPClient interface?
func NewClient(client *http.Client) *Client {
	if client == nil {
		client = defaults.NewHTTPClient()
	}

	return &Client{httpClient: client}
}

// Do executes HTTP request.
//
// Stores the result in the value pointed to by v. If v is nil or not a pointer,
// Do returns an InvalidUnmarshalError.
// Use func `http.NewRequestWithContext` to create `req`.
func (c *Client) Do(req *http.Request, v interface{}, expectedStatusCodes ...int) error {
	// TODO: check that `v` is a pointer or nil
	if req == nil {
		return errors.New("empty request")
	}

	// TODO: add support for the multiple status codes
	if len(expectedStatusCodes) > 1 {
		return errors.New("support for multiple status codes is not implemented")
	}

	// Set defaults:
	if len(expectedStatusCodes) == 0 {
		expectedStatusCodes = []int{http.StatusOK}
	}

	err := c.do(req, v, expectedStatusCodes...)
	if err != nil {
		return NewReqError(req, err)
	}

	return nil
}

func (c *Client) do(req *http.Request, v interface{}, expectedStatusCodes ...int) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != expectedStatusCodes[0] {
		// Non expected status code.

		buf := &strings.Builder{}
		_, err := io.Copy(buf, resp.Body)
		if err != nil {
			err = fmt.Errorf("failed to read API error body: %w", err)
		} else {
			err = errors.New(buf.String())
		}

		return &APIError{
			ResponseStatusCode:  resp.StatusCode,
			ExpectedStatusCodes: expectedStatusCodes,
			// Err is either error body or io.Copy error.
			// TODO: be more specific?
			Err: err,
		}
	}

	if v == nil {
		// nothing to unmarshal
		return nil
	}

	err = json.NewDecoder(resp.Body).Decode(v)
	if err != nil {
		return fmt.Errorf("failed to unmarshal the response body: %w", err)
	}

	return nil
}
