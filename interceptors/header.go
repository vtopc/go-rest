package interceptors

import (
	"net/http"
)

// SetReqHeaderInterceptor sets `value` to header `key` for all requests.
// Could be used for API key auth flow.
func SetReqHeaderInterceptor(client *http.Client, key, value string) {
	SetReqInterceptor(client, setHeader(key, value))
}

func setHeader(key, value string) ReqUpdater {
	return func(req *http.Request) error {
		req.Header.Set(key, value)

		return nil
	}
}
