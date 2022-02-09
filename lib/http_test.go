package stew

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetHTTPResponseBody(t *testing.T) {
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
			if (err != nil) != test.wantErr {
				t.Errorf("getHTTPResponseBody() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if got != test.want {
				t.Errorf("getHTTPResponseBody() = %v, want %v", got, test.want)
			}
		})
	}
}
