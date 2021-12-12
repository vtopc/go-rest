package interceptors

import (
	"errors"
	"net/http"

	"github.com/vtopc/go-rest/defaults"
)

// ReqUpdater updates HTTP request.
type ReqUpdater func(req *http.Request) error

// SetReq updates every HTTP request with function `fn`.
func SetReq(client *http.Client, fn ReqUpdater) error {
	if client == nil {
		return errors.New("no client provided")
	}

	if fn == nil {
		return errors.New("no ReqUpdater provided")
	}

	tr := client.Transport
	if tr == nil {
		tr = defaults.NewHTTPTransport()
	}

	client.Transport = reqInterceptor{
		transport: tr,
		fn:        fn,
	}

	return nil
}

type reqInterceptor struct {
	transport http.RoundTripper
	fn        ReqUpdater
}

// RoundTrip implements http.RoundTripper
func (i reqInterceptor) RoundTrip(req *http.Request) (*http.Response, error) {
	// Do "before sending requests" actions here.

	err := i.fn(req)
	if err != nil {
		return nil, err
	}

	return i.transport.RoundTrip(req)
}
