package interceptors

import (
	"net/http"

	"github.com/vtopc/restclient/defaults"
)

// ReqUpdater updates HTTP request.
type ReqUpdater func(req *http.Request) error

// SetReqInterceptor updates every HTTP request with function `fn`.
func SetReqInterceptor(client *http.Client, fn ReqUpdater) {
	if client == nil {
		// TODO: return error instead?
		client = defaults.NewHTTPClient()
	}

	tr := client.Transport
	if tr == nil {
		tr = defaults.NewHTTPTransport()
	}

	client.Transport = reqInterceptor{
		transport: tr,
		fn:        fn,
	}
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
