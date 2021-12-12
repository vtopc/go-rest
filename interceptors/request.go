package interceptors

import (
	"errors"
	"net/http"

	"github.com/vtopc/restclient/defaults"
)

// ReqUpdater updates HTTP request.
type ReqUpdater func(req *http.Request) error

// SetReqInterceptor updates every HTTP request with function `fn`.
func SetReqInterceptor(client *http.Client, fn ReqUpdater) error {
	if client == nil {
		return errors.New("no client provided")
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