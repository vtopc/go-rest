package interceptors

import (
	"net/http"
)

// SetReqHeader sets `value` to header `key` for all requests.
// Could be used for API key auth flow.
func SetReqHeader(client *http.Client, key, value string) error {
	return SetReq(client, setHeader(key, value))
}

func setHeader(key, value string) ReqUpdater {
	return func(req *http.Request) error {
		req.Header.Set(key, value)

		return nil
	}
}
