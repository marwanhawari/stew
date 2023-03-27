package stew

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/marwanhawari/stew/constants"
	"github.com/marwanhawari/stew/lib/config"
	"github.com/marwanhawari/stew/lib/http"
	"github.com/marwanhawari/stew/lib/ui/prompt"
	"github.com/pkg/errors"
)

// GithubProject contains information about the GitHub project including GitHub releases.
type GithubProject struct {
	Owner    string
	Repo     string
	Releases GithubAPIResponse
}

// GithubAPIResponse is the response from the GitHub releases API.
type GithubAPIResponse []GithubRelease

// GithubRelease contains information about a GitHub release, including the associated assets.
type GithubRelease struct {
	TagName string        `json:"tag_name"`
	Assets  []GithubAsset `json:"assets"`
}

// GithubAsset contains information about a specific GitHub asset.
type GithubAsset struct {
	Name        string `json:"name"`
	DownloadURL string `json:"browser_download_url"`
	Size        int    `json:"size"`
	ContentType string `json:"content_type"`
}

func readGithubJSON(jsonString string) (GithubAPIResponse, error) {
	var ghProject GithubAPIResponse
	err := json.Unmarshal([]byte(jsonString), &ghProject)
	if err != nil {
		return GithubAPIResponse{}, err
	}
	return ghProject, nil
}

func getGithubJSON(config config.Config, owner, repo string) (string, error) {
	url := fmt.Sprintf("https://%s/repos/%v/%v/releases?per_page=100",
		config.GithubAPI, owner, repo)

	response, err := http.ResponseBody(config, url)
	if err != nil {
		return "", err
	}

	return response, nil
}

// NewGithubProject creates a new instance of the GithubProject struct.
func NewGithubProject(config config.Config, owner, repo string) (GithubProject, error) {
	ghJSON, err := getGithubJSON(config, owner, repo)
	if err != nil {
		return GithubProject{}, err
	}

	ghAPIResponse, err := readGithubJSON(ghJSON)
	if err != nil {
		return GithubProject{}, err
	}

	ghProject := GithubProject{Owner: owner, Repo: repo, Releases: ghAPIResponse}

	return ghProject, nil
}

// GetGithubReleasesTags gets a string slice of the releases for a GithubProject.
func GetGithubReleasesTags(project GithubProject) ([]string, error) {
	releasesTags := make([]string, 0, len(project.Releases))

	for _, release := range project.Releases {
		releasesTags = append(releasesTags, release.TagName)
	}

	err := releasesFound(releasesTags, project.Owner, project.Repo)
	if err != nil {
		return []string{}, err
	}

	return releasesTags, nil
}

func releasesFound(releaseTags []string, owner string, repo string) error {
	if len(releaseTags) == 0 {
		return ReleasesNotFoundError{Owner: owner, Repo: repo}
	}
	return nil
}

// GetGithubReleasesAssets gets a string slice of the assets for a GithubRelease.
func GetGithubReleasesAssets(ghProject GithubProject, tag string) ([]string, error) {
	var releaseAssets []string

	for _, release := range ghProject.Releases {
		if release.TagName == tag {
			for _, asset := range release.Assets {
				releaseAssets = append(releaseAssets, asset.Name)
			}
		}
	}

	err := assetsFound(releaseAssets, tag)
	if err != nil {
		return []string{}, err
	}

	return releaseAssets, nil
}

func assetsFound(releaseAssets []string, releaseTag string) error {
	if len(releaseAssets) == 0 {
		return AssetsNotFoundError{Tag: releaseTag}
	}
	return nil
}

// DetectAsset will automatically detect a release asset matching your
// systems OS/arch or prompt you to manually select an asset.
func DetectAsset(rt config.Runtime, releaseAssets []string) (string, error) {
	var detectedOSAssets []string
	var reOS *regexp.Regexp
	switch rt.OS {
	case "darwin":
		reOS = constants.RegexDarwin
	case "windows":
		reOS = constants.RegexWindows
	default:
		reOS = regexp.MustCompile(`(?i)` + rt.OS)
	}

	for _, asset := range releaseAssets {
		if reOS.MatchString(asset) {
			detectedOSAssets = append(detectedOSAssets, asset)
		}
	}

	var detectedFinalAssets []string
	var reArch *regexp.Regexp
	switch rt.Arch {
	case "arm64":
		reArch = constants.RegexArm64
	case "amd64":
		reArch = constants.RegexAmd64
	case "386":
		reArch = constants.Regex386
	default:
		reArch = regexp.MustCompile(`(?i)` + rt.Arch)
	}

	for _, asset := range detectedOSAssets {
		if reArch.MatchString(asset) {
			detectedFinalAssets = append(detectedFinalAssets, asset)
		}
	}

	var finalAsset string
	var err error
	if len(detectedFinalAssets) != 1 {
		if rt.OS == "darwin" && rt.Arch == "arm64" {
			finalAsset = darwinARMFallback(detectedOSAssets)
		}
		if finalAsset == "" {
			finalAsset, err = prompt.SelectWarn(rt,
				"Could not automatically detect the release asset matching "+
					"your OS/Arch. Please select it manually:", releaseAssets)
			if err != nil {
				return "", err
			}
		}
	} else {
		finalAsset = detectedFinalAssets[0]
	}

	return finalAsset, nil
}

func darwinARMFallback(darwinAssets []string) string {
	reArch := constants.RegexAmd64

	var altAssets []string
	for _, asset := range darwinAssets {
		if reArch.MatchString(asset) {
			altAssets = append(altAssets, asset)
		}
	}

	if len(altAssets) != 1 {
		return ""
	}

	return altAssets[0]
}

// GithubSearch contains information about the GitHub search including the GitHub search results.
type GithubSearch struct {
	SearchQuery string
	Count       int                  `json:"total_count"`
	Items       []GithubSearchResult `json:"items"`
}

// GithubSearchResult contains information about the GitHub search result.
type GithubSearchResult struct {
	FullName    string `json:"full_name"`
	Stars       int    `json:"stargazers_count"`
	Language    string `json:"language"`
	Description string `json:"description"`
}

func getGithubSearchJSON(config config.Config, searchQuery string) (string, error) {
	url := fmt.Sprintf("https://%s/search/repositories?q=%v%v",
		config.GithubAPI, searchQuery, "+fork:true")

	response, err := http.ResponseBody(config, url)
	if err != nil {
		return "", err
	}

	return response, nil
}

func readGithubSearchJSON(jsonString string) (GithubSearch, error) {
	var ghSearch GithubSearch
	err := json.Unmarshal([]byte(jsonString), &ghSearch)
	if err != nil {
		return GithubSearch{}, errors.WithStack(err)
	}
	return ghSearch, nil
}

// NewGithubSearch creates a new instance of the GithubSearch struct.
func NewGithubSearch(config config.Config, searchQuery string) (GithubSearch, error) {
	if err := ValidateGithubSearchQuery(searchQuery); err != nil {
		return GithubSearch{}, err
	}
	ghJSON, err := getGithubSearchJSON(config, searchQuery)
	if err != nil {
		return GithubSearch{}, err
	}

	ghSearch, err := readGithubSearchJSON(ghJSON)
	if err != nil {
		return GithubSearch{}, err
	}

	ghSearch.SearchQuery = searchQuery

	return ghSearch, nil
}

// FormatSearchResults formats the GitHub search results for the terminal UI.
func FormatSearchResults(ghSearch GithubSearch) []string {
	var formattedSearchResults []string
	for _, searchResult := range ghSearch.Items {
		formatted := fmt.Sprintf("%v [⭐️%v] %v", searchResult.FullName, searchResult.Stars, searchResult.Description)
		formattedSearchResults = append(formattedSearchResults, formatted)
	}

	return formattedSearchResults
}

// ValidateGithubSearchQuery makes sure the GitHub search query is valid.
func ValidateGithubSearchQuery(query string) error {
	if !constants.RegexGithubSearch.MatchString(query) {
		return InvalidGithubSearchQueryError{Query: query}
	}

	return nil
}
