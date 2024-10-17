package stew

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/marwanhawari/stew/constants"
)

// GithubProject contains information about the GitHub project including GitHub releases
type GithubProject struct {
	Owner    string
	Repo     string
	Releases GithubAPIResponse
}

// GithubAPIResponse is the response from the GitHub releases API
type GithubAPIResponse []GithubRelease

// GithubRelease contains information about a GitHub release, including the associated assets
type GithubRelease struct {
	TagName string        `json:"tag_name"`
	Assets  []GithubAsset `json:"assets"`
}

// GithubAsset contains information about a specific GitHub asset
type GithubAsset struct {
	Name        string `json:"name"`
	DownloadURL string `json:"url"`
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

func getGithubJSON(owner, repo string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%v/%v/releases?per_page=100", owner, repo)

	response, err := getHTTPResponseBody(url)
	if err != nil {
		return "", err
	}

	return response, nil
}

// NewGithubProject creates a new instance of the GithubProject struct
func NewGithubProject(owner, repo string) (GithubProject, error) {
	ghJSON, err := getGithubJSON(owner, repo)
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

// GetGithubReleasesTags gets a string slice of the releases for a GithubProject
func GetGithubReleasesTags(ghProject GithubProject) ([]string, error) {
	releasesTags := []string{}

	for _, release := range ghProject.Releases {
		releasesTags = append(releasesTags, release.TagName)
	}

	err := releasesFound(releasesTags, ghProject.Owner, ghProject.Repo)
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

// GetGithubReleasesAssets gets a string slice of the assets for a GithubRelease
func GetGithubReleasesAssets(ghProject GithubProject, tag string) ([]string, error) {

	releaseAssets := []string{}

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

func filterReleaseAssets(assets []string) []string {
	var filteredAssets []string
	re := regexp.MustCompile(`\.(sha(256|512)(sum)?)$`)

	for _, asset := range assets {
		if re.MatchString(asset) {
			continue
		}
		filteredAssets = append(filteredAssets, asset)
	}
	return filteredAssets
}

// DetectAsset will automatically detect a release asset matching your systems OS/arch or prompt you to manually select an asset
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

	filteredReleaseAssets := filterReleaseAssets(releaseAssets)

	for _, asset := range filteredReleaseAssets {
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
			finalAsset, err = WarningPromptSelect("Could not automatically detect the release asset matching your OS/Arch. Please select it manually:", detectedFinalAssets)
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

// GithubSearch contains information about the GitHub search including the GitHub search results
type GithubSearch struct {
	SearchQuery string
	Count       int                  `json:"total_count"`
	Items       []GithubSearchResult `json:"items"`
}

// GithubSearchResult contains information about the GitHub search result
type GithubSearchResult struct {
	FullName    string `json:"full_name"`
	Stars       int    `json:"stargazers_count"`
	Language    string `json:"language"`
	Description string `json:"description"`
}

func getGithubSearchJSON(searchQuery string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/search/repositories?q=%v%v", searchQuery, "+fork:true")

	response, err := getHTTPResponseBody(url)
	if err != nil {
		return "", err
	}

	return response, nil
}

func readGithubSearchJSON(jsonString string) (GithubSearch, error) {
	var ghSearch GithubSearch
	err := json.Unmarshal([]byte(jsonString), &ghSearch)
	if err != nil {
		return GithubSearch{}, err
	}
	return ghSearch, nil
}

// NewGithubSearch creates a new instance of the GithubSearch struct
func NewGithubSearch(searchQuery string) (GithubSearch, error) {
	ghJSON, err := getGithubSearchJSON(searchQuery)
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

// FormatSearchResults formats the GitHub search results for the terminal UI
func FormatSearchResults(ghSearch GithubSearch) []string {

	var formattedSearchResults []string
	for _, searchResult := range ghSearch.Items {
		formatted := fmt.Sprintf("%v [⭐️%v] %v", searchResult.FullName, searchResult.Stars, searchResult.Description)
		formattedSearchResults = append(formattedSearchResults, formatted)
	}

	return formattedSearchResults
}

// ValidateGithubSearchQuery makes sure the GitHub search query is valid
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
