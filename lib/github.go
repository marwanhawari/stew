package stew

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/marwanhawari/stew/constants"
)

type GithubProject struct {
	Owner    string
	Repo     string
	Releases GithubAPIResponse
}

type GithubAPIResponse []GithubRelease

type GithubRelease struct {
	TagName string        `json:"tag_name"`
	Assets  []GithubAsset `json:"assets"`
}
type GithubAsset struct {
	Name        string `json:"name"`
	DownloadURL string `json:"browser_download_url"`
	Size        int    `json:"size"`
	ContentType string `json:"content_type"`
}

func ReadGithubJSON(jsonString string) (GithubAPIResponse, error) {
	var ghProject GithubAPIResponse
	err := json.Unmarshal([]byte(jsonString), &ghProject)
	if err != nil {
		return GithubAPIResponse{}, err
	}
	return ghProject, nil
}

func getGithubJSON(owner, repo string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%v/%v/releases?per_page=100", owner, repo)

	response, err := GetHTTPResponseBody(url)
	if err != nil {
		return "", err
	}

	return response, nil
}

func NewGithubProject(owner, repo string) (GithubProject, error) {
	ghJSON, err := getGithubJSON(owner, repo)
	if err != nil {
		return GithubProject{}, err
	}

	ghAPIResponse, err := ReadGithubJSON(ghJSON)
	if err != nil {
		return GithubProject{}, err
	}

	ghProject := GithubProject{Owner: owner, Repo: repo, Releases: ghAPIResponse}

	return ghProject, nil
}

func GetGithubReleasesTags(ghProject GithubProject) ([]string, error) {
	releasesTags := []string{}

	for _, release := range ghProject.Releases {
		releasesTags = append(releasesTags, release.TagName)
	}

	err := ReleasesFound(releasesTags, ghProject.Owner, ghProject.Repo)
	if err != nil {
		return []string{}, err
	}

	return releasesTags, nil

}

func ReleasesFound(releaseTags []string, owner string, repo string) error {
	if len(releaseTags) == 0 {
		return ReleasesNotFoundError{Owner: owner, Repo: repo}
	}
	return nil
}

func GetGithubReleasesAssets(ghProject GithubProject, tag string) ([]string, error) {

	releaseAssets := []string{}

	for _, release := range ghProject.Releases {
		if release.TagName == tag {
			for _, asset := range release.Assets {
				releaseAssets = append(releaseAssets, asset.Name)
			}
		}
	}

	err := AssetsFound(releaseAssets, tag)
	if err != nil {
		return []string{}, err
	}

	return releaseAssets, nil

}

func AssetsFound(releaseAssets []string, releaseTag string) error {
	if len(releaseAssets) == 0 {
		return AssetsNotFoundError{Tag: releaseTag}
	}
	return nil
}

func DetectAsset(userOS string, userArch string, releaseAssets []string) (string, error) {
	var detectedOSAssets []string
	var reOS *regexp.Regexp
	var err error
	switch userOS {
	case "darwin":
		reOS, err = regexp.Compile(constants.RegexDarwin)
	case "windows":
		reOS, err = regexp.Compile(constants.RegexWindows)
	default:
		reOS, err = regexp.Compile(`(?i)` + userOS)
	}
	if err != nil {
		return "", err
	}

	for _, asset := range releaseAssets {
		if reOS.MatchString(asset) {
			detectedOSAssets = append(detectedOSAssets, asset)
		}
	}

	var detectedFinalAssets []string
	var reArch *regexp.Regexp
	switch userArch {
	case "arm64":
		reArch, err = regexp.Compile(constants.RegexArm64)
	case "amd64":
		reArch, err = regexp.Compile(constants.RegexAmd64)
	case "386":
		reArch, err = regexp.Compile(constants.Regex386)
	default:
		reArch, err = regexp.Compile(`(?i)` + userArch)
	}
	if err != nil {
		return "", err
	}

	for _, asset := range detectedOSAssets {
		if reArch.MatchString(asset) {
			detectedFinalAssets = append(detectedFinalAssets, asset)
		}
	}

	var finalAsset string
	if len(detectedFinalAssets) != 1 {
		if userOS == "darwin" && userArch == "arm64" {
			finalAsset, err = darwinARMFallback(detectedOSAssets)
			if err != nil {
				return "", err
			}
		}
		if finalAsset == "" {
			finalAsset, err = WarningPromptSelect("Could not automatically detect the release asset matching your OS/Arch. Please select it manually:", releaseAssets)
			if err != nil {
				return "", err
			}
		}
	} else {
		finalAsset = detectedFinalAssets[0]
	}

	return finalAsset, nil
}

func darwinARMFallback(darwinAssets []string) (string, error) {
	reArch, err := regexp.Compile(constants.RegexAmd64)
	if err != nil {
		return "", err
	}

	var altAssets []string
	for _, asset := range darwinAssets {
		if reArch.MatchString(asset) {
			altAssets = append(altAssets, asset)
		}
	}

	if len(altAssets) != 1 {
		return "", nil
	}

	return altAssets[0], nil
}

type GithubSearch struct {
	SearchQuery string
	Count       int                  `json:"total_count"`
	Items       []GithubSearchResult `json:"items"`
}

type GithubSearchResult struct {
	FullName    string `json:"full_name"`
	Stars       int    `json:"stargazers_count"`
	Language    string `json:"language"`
	Description string `json:"description"`
}

func getGithubSearchJSON(searchQuery string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/search/repositories?q=%v", searchQuery)

	response, err := GetHTTPResponseBody(url)
	if err != nil {
		return "", err
	}

	return response, nil
}

func ReadGithubSearchJSON(jsonString string) (GithubSearch, error) {
	var ghSearch GithubSearch
	err := json.Unmarshal([]byte(jsonString), &ghSearch)
	if err != nil {
		return GithubSearch{}, err
	}
	return ghSearch, nil
}

func NewGithubSearch(searchQuery string) (GithubSearch, error) {
	ghJSON, err := getGithubSearchJSON(searchQuery)
	if err != nil {
		return GithubSearch{}, err
	}

	ghSearch, err := ReadGithubSearchJSON(ghJSON)
	if err != nil {
		return GithubSearch{}, err
	}

	ghSearch.SearchQuery = searchQuery

	return ghSearch, nil
}

func FormatSearchResults(ghSearch GithubSearch) []string {

	var formattedSearchResults []string
	for _, searchResult := range ghSearch.Items {
		formatted := fmt.Sprintf("%v [⭐️%v] %v", searchResult.FullName, searchResult.Stars, searchResult.Description)
		formattedSearchResults = append(formattedSearchResults, formatted)
	}

	return formattedSearchResults
}

func ValidateGithubSearchQuery(searchQuery string) error {

	reSearch, err := regexp.Compile(constants.RegexGithubSearch)
	if err != nil {
		return err
	}

	if !reSearch.MatchString(searchQuery) {
		return InvalidGithubSearchQueryError{SearchQuery: searchQuery}
	}

	return nil
}
