package stew

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func GetHTTPResponseBody(url string) (string, error) {
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

	if res.StatusCode != http.StatusOK {
		return "", NonZeroStatusCodeError{res.StatusCode}
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
