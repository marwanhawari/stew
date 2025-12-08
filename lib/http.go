package stew

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
)

func getHTTPResponseBody(url string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	if strings.Contains(url, "api.github.com") {
		req.Header.Add("Accept", "application/vnd.github.v3+json")
		githubToken := os.Getenv("GITHUB_TOKEN")
		if githubToken != "" {
			slog.Debug("using gitub token", "token", githubToken)
			req.Header.Add("Authorization", fmt.Sprintf("token %v", githubToken))
		}
	}

	slog.Debug("HTTP request", "url", req.URL)

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	slog.Debug("HTTP response", "code", res.StatusCode, "status", res.Status)

	if res.StatusCode != http.StatusOK {
		return "", NonZeroStatusCodeError{res.StatusCode}
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
