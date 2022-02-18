package rest

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vtopc/go-rest/defaults"
)

func TestAPIError_Error(t *testing.T) {
	tests := map[string]struct {
		err  *APIError
		want string
	}{
		"all_fields": {
			err: &APIError{
				Resp:                &http.Response{StatusCode: 400},
				ExpectedStatusCodes: []int{200},
				Err:                 errors.New("internal fail"),
			},
			want: "wrong status code (400 not in [200]): internal fail",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.EqualError(t, tt.err, tt.want)
		})
	}
}

func TestStatusCodeFromAPIError(t *testing.T) {
	client := NewClient(defaults.NewHTTPClient())

	type S struct {
		Foo string `json:"foo"`
	}

	tests := map[string]struct {
		statusCode         int // status code to return
		expectedStatusCode int // for .Do(...)
		body               []byte

		v interface{} // for .Do(...)

		wantErr        bool
		wantStatusCode int
	}{
		"no_errors": {
			statusCode:         200,
			expectedStatusCode: 200,
			body:               []byte(`{"foo":"bar"}`),
			v:                  new(S),
			wantStatusCode:     0,
		},
		"bad_request": {
			statusCode:         400,
			expectedStatusCode: 201,
			body:               []byte(`{"errors":[{"message":"test error"}]}`),
			wantErr:            true,
			wantStatusCode:     400,
		},
		"not_found": {
			statusCode:         404,
			expectedStatusCode: 200,
			body:               []byte(`{"errors":[{"message":"the entity not found"}]}`),
			wantErr:            true,
			wantStatusCode:     404,
		},
		"not_APIError": {
			statusCode:         200,
			expectedStatusCode: 200,
			body:               []byte(`{"foo":"bar"}`),
			wantErr:            true,
			v:                  make(chan struct{}), // a channel just to fail Unmarshal
			wantStatusCode:     500,
		},
		"not_API_error_or_misconfiguration": {
			statusCode:         201,
			expectedStatusCode: 200,
			body:               []byte(`{"foo":"bar"}`),
			v:                  new(S),
			wantErr:            true,
			wantStatusCode:     500,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			handler := func(w http.ResponseWriter, req *http.Request) {
				w.WriteHeader(tt.statusCode)

				if tt.body != nil {
					_, err := w.Write(tt.body)
					require.NoError(t, err, "failed to write the body")
				}
			}

			ts := httptest.NewServer(http.HandlerFunc(handler))
			defer ts.Close()

			req, err := http.NewRequest(http.MethodGet, ts.URL+"/foo", nil)
			require.NoError(t, err)

			err = client.Do(req, tt.v, tt.expectedStatusCode)
			t.Logf("got error: %v", err)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			// got := StatusCodeFromAPIError(err)
			// assert.Equal(t, tt.wantStatusCode, got)
		})
	}
}
