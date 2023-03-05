package http

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/marwanhawari/stew/constants"
	"github.com/marwanhawari/stew/lib/config"
	"github.com/pkg/errors"
)

// ResponseBody returns the response body for given URL address.
func ResponseBody(cfg config.Config, url string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", errors.WithStack(err)
	}

	if strings.Contains(url, cfg.GithubAPI) {
		req.Header.Add("Accept", "application/vnd.github.v3+json")
		githubToken := os.Getenv("GITHUB_TOKEN")
		if githubToken == "" {
			githubToken = cfg.GithubToken
		}
		if githubToken != "" {
			req.Header.Add("Authorization", fmt.Sprintf("token %v", githubToken))
		}
	}

	res, err := client.Do(req)
	if err != nil {
		return "", errors.WithStack(err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", NonZeroStatusCodeError{res.StatusCode}
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return string(body), nil
}

// NonZeroStatusCodeError occurs if a non-zero status code is received from an HTTP request.
type NonZeroStatusCodeError struct {
	StatusCode int
}

func (e NonZeroStatusCodeError) Error() string {
	return fmt.Sprintf(
		"%v Received non-zero status code from HTTP request: %v",
		constants.RedColor("Error:"),
		constants.RedColor(e.StatusCode),
	)
}
