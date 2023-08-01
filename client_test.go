package rest

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Struct struct {
	Status string `json:"status"`
}

func TestClientDo(t *testing.T) {
	tests := map[string]struct {
		method              string
		urlPostfix          string
		statusCode          int
		expectedStatusCodes []int
		respBody            []byte
		v                   interface{}
		want                interface{}
		wantErr             error
		wantWrappedErr      error
	}{
		"positive_get": {
			method:              http.MethodGet,
			urlPostfix:          "/health",
			statusCode:          http.StatusOK,
			expectedStatusCodes: []int{http.StatusOK},
			respBody:            []byte(`{"status":"ok"}`),
			v:                   &Struct{},
			want:                &Struct{Status: "ok"},
		},

		"positive_multiple_statuses": {
			method:              http.MethodPost,
			urlPostfix:          "/health",
			statusCode:          http.StatusAccepted,
			expectedStatusCodes: []int{http.StatusOK, http.StatusAccepted},
			respBody:            []byte(`{"status":"ok"}`),
			v:                   &Struct{},
			want:                &Struct{Status: "ok"},
		},

		"positive_but_wrong_payload": {
			method:              http.MethodGet,
			urlPostfix:          "/health",
			statusCode:          http.StatusOK,
			expectedStatusCodes: []int{http.StatusOK},
			respBody:            []byte(`{"error":"some error"}`),
			v:                   &Struct{},
			want:                &Struct{}, // zero values
		},

		"positive_get_empty_body": {
			method:              http.MethodGet,
			urlPostfix:          "/health",
			statusCode:          http.StatusNoContent,
			expectedStatusCodes: []int{http.StatusNoContent},
			respBody:            nil,
			v:                   nil,
			want:                nil,
		},

		"negative_wrong_status_code": {
			method:              http.MethodGet,
			urlPostfix:          "/health",
			statusCode:          http.StatusInternalServerError,
			expectedStatusCodes: []int{http.StatusOK, http.StatusAccepted},
			respBody:            []byte(`{"error":"some error"}`),
			wantWrappedErr:      errors.New("wrong status code (500 not in [200, 202]): {\"error\":\"some error\"}"),
		},

		"negative_not_a_pointer": {
			method:              http.MethodGet,
			urlPostfix:          "/health",
			statusCode:          http.StatusOK,
			expectedStatusCodes: []int{http.StatusOK},
			respBody:            []byte(`{"status":"ok"}`),
			v:                   Struct{},
			wantErr:             errors.New("json: Unmarshal(non-pointer rest.Struct)"),
		},

		// TODO: add more test cases
	}

	for k, tt := range tests {
		tt := tt
		t.Run(k, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tt.method, r.Method)
				hasSuffix(t, r.URL.String(), tt.urlPostfix)

				// Send response to be tested
				w.WriteHeader(tt.statusCode)

				if tt.respBody != nil {
					_, err := w.Write(tt.respBody)
					require.NoError(t, err, "failed to write the body")
				}
			}))
			defer server.Close()

			c := Client{
				httpClient: server.Client(),
			}

			req, err := http.NewRequest(tt.method, server.URL+tt.urlPostfix, nil)
			require.NoError(t, err)

			// test:
			err = c.Do(req, tt.v, tt.expectedStatusCodes...)
			if tt.wantErr != nil {
				require.EqualError(t, err, tt.wantErr.Error())
				return
			}
			if tt.wantWrappedErr != nil {
				require.EqualError(t, errors.Unwrap(err), tt.wantWrappedErr.Error())
				return
			}

			require.NoError(t, err)
			if tt.v != nil {
				assert.EqualValues(t, tt.want, tt.v)
			}
		})
	}
}

func hasSuffix(t *testing.T, s, suffix string) {
	assert.Truef(t, strings.HasSuffix(s, suffix), "expected '%s' to ends with suffix '%s'", s, suffix)
}
