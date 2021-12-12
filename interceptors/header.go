package interceptors

import (
	"net/http"

	"github.com/vtopc/restclient/defaults"
)

// SetReqHeaderInterceptor sets `value` to header `key` for all requests.
// Could be used for API key auth flow.
func SetReqHeaderInterceptor(client *http.Client, key, value string) {
	if client == nil {
		// TODO: return error instead?
		client = defaults.NewHTTPClient()
	}

	tr := client.Transport
	if tr == nil {
		tr = defaults.NewHTTPTransport()
	}

	client.Transport = headerInterceptor{
		transport: tr,
		key:       key,
		value:     value,
	}
}

type headerInterceptor struct {
	transport http.RoundTripper
	key       string
	value     string
}

// RoundTrip implements http.RoundTripper
func (i headerInterceptor) RoundTrip(req *http.Request) (*http.Response, error) {
	// Do "before sending requests" actions here.

	req.Header.Set(i.key, i.value)

	return i.transport.RoundTrip(req)
}
