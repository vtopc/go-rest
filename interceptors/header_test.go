package interceptors_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vtopc/go-rest/interceptors"
)

func TestSetReqHeader(t *testing.T) {
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
	err := interceptors.SetReqHeader(client, k, v)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	require.NoError(t, err)

	// test:
	_, err = client.Do(req)
	require.NoError(t, err)
}
