package stew

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetHTTPResponseBody(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	type Test struct {
		name    string
		server  *httptest.Server
		want    string
		wantErr bool
		err     error
	}

	tests := []Test{
		{
			name: "test1",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"test":"ok"}`))
			})),
			want:    `{"test":"ok"}`,
			wantErr: false,
			err:     nil,
		},
		{
			name: "test2",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte(``))
			})),
			want:    "",
			wantErr: true,
			err:     NonZeroStatusCodeError{403},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer test.server.Close()

			got, err := getHTTPResponseBody(test.server.URL)
			if test.wantErr {
				require.Error(err, "Expected an error but got none")
				assert.IsType(test.err, err, "Error type mismatch")
			} else {
				require.NoError(err, "Unexpected error")
				assert.Equal(test.want, got, "Response body mismatch")
			}
		})
	}
}
