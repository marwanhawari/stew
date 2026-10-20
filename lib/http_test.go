package stew

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
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

func TestGetHTTPResponseBodyRateLimit(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		remaining  string
		reset      string
		want       error
	}{
		{
			name:       "403 with rate limit used up",
			statusCode: http.StatusForbidden,
			remaining:  "0",
			reset:      "1792513240",
			want:       GithubRateLimitError{Reset: time.Unix(1792513240, 0)},
		},
		{
			name:       "429 with rate limit used up",
			statusCode: http.StatusTooManyRequests,
			remaining:  "0",
			reset:      "1792513240",
			want:       GithubRateLimitError{Reset: time.Unix(1792513240, 0)},
		},
		{
			name:       "rate limit used up without a reset time",
			statusCode: http.StatusForbidden,
			remaining:  "0",
			reset:      "",
			want:       GithubRateLimitError{},
		},
		{
			name:       "403 with rate limit remaining",
			statusCode: http.StatusForbidden,
			remaining:  "10",
			reset:      "1792513240",
			want:       NonZeroStatusCodeError{http.StatusForbidden},
		},
		{
			name:       "403 without rate limit headers",
			statusCode: http.StatusForbidden,
			want:       NonZeroStatusCodeError{http.StatusForbidden},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.remaining != "" {
					w.Header().Set("X-RateLimit-Remaining", tt.remaining)
				}
				if tt.reset != "" {
					w.Header().Set("X-RateLimit-Reset", tt.reset)
				}
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			_, err := getHTTPResponseBody(server.URL)
			var rateLimitErr GithubRateLimitError
			if errors.As(tt.want, &rateLimitErr) {
				var gotErr GithubRateLimitError
				if !errors.As(err, &gotErr) {
					t.Fatalf("getHTTPResponseBody() error = %v, want %v", err, tt.want)
				}
				if !gotErr.Reset.Equal(rateLimitErr.Reset) || gotErr.HasToken != rateLimitErr.HasToken {
					t.Errorf("getHTTPResponseBody() error = %+v, want %+v", gotErr, rateLimitErr)
				}
				return
			}
			if err != tt.want {
				t.Errorf("getHTTPResponseBody() error = %v, want %v", err, tt.want)
			}
		})
	}
}
