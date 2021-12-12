package interceptors

import (
	"net/http"
)

const ContentTypeHeaderValue = "application/json; charset=utf-8"

const contentTypeHeaderName = "Content-Type"

// SetReqContentType sets Content-Type header to 'v' for requests with not empty body.
func SetReqContentType(client *http.Client, v string) error {
	return SetReq(client, setContentType(v))
}

func setContentType(v string) ReqUpdater {
	return func(req *http.Request) error {
		if req.Body != nil && req.Body != http.NoBody {
			req.Header.Set(contentTypeHeaderName, v)
		}

		return nil
	}
}
