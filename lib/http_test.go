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

func Test_isGithubAPIURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want bool
	}{
		{
			name: "releases endpoint",
			url:  "https://api.github.com/repos/marwanhawari/ppath/releases?per_page=100",
			want: true,
		},
		{
			name: "asset endpoint",
			url:  "https://api.github.com/repos/marwanhawari/ppath/releases/assets/1",
			want: true,
		},
		{
			name: "http scheme",
			url:  "http://api.github.com/repos/marwanhawari/ppath/releases",
			want: false,
		},
		{
			name: "host name in the path",
			url:  "https://example.com/api.github.com/binary",
			want: false,
		},
		{
			name: "host name in the query",
			url:  "https://example.com/binary?mirror=api.github.com",
			want: false,
		},
		{
			name: "host name as a subdomain of a different host",
			url:  "https://api.github.com.example.com/binary",
			want: false,
		},
		{
			name: "host name as user info",
			url:  "https://api.github.com@example.com/binary",
			want: false,
		},
		{
			name: "github release download",
			url:  "https://github.com/marwanhawari/ppath/releases/download/v0.0.3/ppath-v0.0.3-darwin-arm64.tar.gz",
			want: false,
		},
		{
			name: "empty",
			url:  "",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isGithubAPIURL(tt.url); got != tt.want {
				t.Errorf("isGithubAPIURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
