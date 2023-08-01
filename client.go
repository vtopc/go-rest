package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
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
// Stores the result in the value pointed to by v. If v is not a nil and not a pointer,
// Do returns a json.InvalidUnmarshalError.
// Use func `http.NewRequestWithContext` to create `req`.
func (c *Client) Do(req *http.Request, v interface{}, expectedStatusCodes ...int) error {
	// check that `v` is a pointer or nil before doing the request.
	if v != nil {
		rv := reflect.ValueOf(v)
		if rv.Kind() != reflect.Ptr {
			return &json.InvalidUnmarshalError{Type: reflect.TypeOf(v)}
		}
	}

	if req == nil {
		return errors.New("empty request")
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

	var expected bool
	for _, code := range expectedStatusCodes {
		if resp.StatusCode == code {
			expected = true
			break
		}
	}

	if !expected {
		// Non expected status code.

		buf := &strings.Builder{}
		_, err := io.Copy(buf, resp.Body)
		if err != nil {
			err = fmt.Errorf("failed to read API error body: %w", err)
		} else {
			err = errors.New(buf.String())
		}

		return &APIError{
			Resp:                resp,
			ExpectedStatusCodes: expectedStatusCodes,
			Err:                 err,
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
