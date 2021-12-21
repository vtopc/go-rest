package defaults

import (
	"net/http"
	"time"
)

const (
	defaultTimeout = 30 * time.Second
)

func NewHTTPClient() *http.Client {
	return &http.Client{
		Timeout:   defaultTimeout,
		Transport: NewHTTPTransport(),
	}
}

func NewHTTPTransport() http.RoundTripper {
	return http.DefaultTransport
}
