package interceptors_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vtopc/restclient/interceptors"
)

func TestSetReqHeaderInterceptor(t *testing.T) {
	const (
		k = "foo"
		v = "bar"
	)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got := r.Header.Get(k)
		t.Log(got)
		assert.Equal(t, v, got)

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := server.Client()
	err := interceptors.SetReqHeaderInterceptor(client, k, v)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	require.NoError(t, err)

	// test:
	_, err = client.Do(req)
	require.NoError(t, err)
}
