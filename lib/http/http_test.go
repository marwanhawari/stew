package http_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/marwanhawari/stew/constants"
	stewhttp "github.com/marwanhawari/stew/lib/http"
	"github.com/marwanhawari/stew/lib/testsupport"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetHTTPResponseBody(t *testing.T) {
	type TestCase struct {
		name    string
		server  *httptest.Server
		want    string
		wantErr bool
		err     error
	}

	testCases := []TestCase{{
		name: "test1",
		server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"test":"ok"}`))
		})),
		want:    `{"test":"ok"}`,
		wantErr: false,
		err:     nil,
	}, {
		name: "test2",
		server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte(``))
		})),
		want:    "",
		wantErr: true,
		err:     stewhttp.NonZeroStatusCodeError{StatusCode: 403},
	}}
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			defer test.server.Close()

			rt := testsupport.NewDefaultRuntime(t)
			got, err := stewhttp.ResponseBody(rt.Config, test.server.URL)
			require.ErrorIs(t, err, test.err)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestNonZeroStatusCodeError_Error(t *testing.T) {
	type fields struct {
		StatusCode int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "test1",
			fields: fields{
				StatusCode: 1,
			},
			want: fmt.Sprintf("%v Received non-zero status code from HTTP request: %v", constants.RedColor("Error:"), constants.RedColor(1)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := stewhttp.NonZeroStatusCodeError{
				StatusCode: tt.fields.StatusCode,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("NonZeroStatusCodeError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}
