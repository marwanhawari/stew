package stew

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
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
			req.Header.Add("Authorization", fmt.Sprintf("token %v", githubToken))
		}
	}

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	if isRateLimited(res) {
		return "", GithubRateLimitError{Reset: rateLimitReset(res), HasToken: req.Header.Get("Authorization") != ""}
	}

	if res.StatusCode != http.StatusOK {
		return "", NonZeroStatusCodeError{res.StatusCode}
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// isRateLimited checks if GitHub refused a request because the rate limit is used up
func isRateLimited(res *http.Response) bool {
	if res.StatusCode != http.StatusForbidden && res.StatusCode != http.StatusTooManyRequests {
		return false
	}
	return res.Header.Get("X-RateLimit-Remaining") == "0"
}

// rateLimitReset gets the time when the GitHub rate limit resets. It returns the zero time if the header is missing or not valid
func rateLimitReset(res *http.Response) time.Time {
	resetUnix, err := strconv.ParseInt(res.Header.Get("X-RateLimit-Reset"), 10, 64)
	if err != nil {
		return time.Time{}
	}
	return time.Unix(resetUnix, 0)
}
