package stew_test

import (
	"testing"

	stew "github.com/marwanhawari/stew/lib"
	"github.com/marwanhawari/stew/lib/config"
	stewhttp "github.com/marwanhawari/stew/lib/http"
	"github.com/marwanhawari/stew/lib/testsupport"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testGithubRelease0 = stew.GithubRelease{
	TagName: "v0.0.3",
	Assets: []stew.GithubAsset{{
		Name:        "checksums.txt",
		DownloadURL: "https://github.com/marwanhawari/ppath/releases/download/v0.0.3/checksums.txt",
		Size:        394,
		ContentType: "text/plain; charset=utf-8",
	}, {
		Name:        "ppath-v0.0.3-darwin-amd64.tar.gz",
		DownloadURL: "https://github.com/marwanhawari/ppath/releases/download/v0.0.3/ppath-v0.0.3-darwin-amd64.tar.gz",
		Size:        626448,
		ContentType: "application/gzip",
	}, {
		Name:        "ppath-v0.0.3-darwin-arm64.tar.gz",
		DownloadURL: "https://github.com/marwanhawari/ppath/releases/download/v0.0.3/ppath-v0.0.3-darwin-arm64.tar.gz",
		Size:        625832,
		ContentType: "application/gzip",
	}, {
		Name:        "ppath-v0.0.3-linux-amd64.tar.gz",
		DownloadURL: "https://github.com/marwanhawari/ppath/releases/download/v0.0.3/ppath-v0.0.3-linux-amd64.tar.gz",
		Size:        567449,
		ContentType: "application/gzip",
	}, {
		Name:        "ppath-v0.0.3-linux-arm64.tar.gz",
		DownloadURL: "https://github.com/marwanhawari/ppath/releases/download/v0.0.3/ppath-v0.0.3-linux-arm64.tar.gz",
		Size:        529321,
		ContentType: "application/gzip",
	}},
}
var testGithubRelease1 = stew.GithubRelease{
	TagName: "v0.0.2",
	Assets: []stew.GithubAsset{{
		Name:        "ppath-v0.0.2-darwin-amd64.tar.gz",
		DownloadURL: "https://github.com/marwanhawari/ppath/releases/download/v0.0.2/ppath-v0.0.2-darwin-amd64.tar.gz",
		Size:        1139147,
		ContentType: "application/x-gtar",
	}, {
		Name:        "ppath-v0.0.2-darwin-arm64.tar.gz",
		DownloadURL: "https://github.com/marwanhawari/ppath/releases/download/v0.0.2/ppath-v0.0.2-darwin-arm64.tar.gz",
		Size:        1100483,
		ContentType: "application/x-gtar",
	}, {
		Name:        "ppath-v0.0.2-linux-amd64.tar.gz",
		DownloadURL: "https://github.com/marwanhawari/ppath/releases/download/v0.0.2/ppath-v0.0.2-linux-amd64.tar.gz",
		Size:        1092421,
		ContentType: "application/x-gtar",
	}, {
		Name:        "ppath-v0.0.2-linux-arm64.tar.gz",
		DownloadURL: "https://github.com/marwanhawari/ppath/releases/download/v0.0.2/ppath-v0.0.2-linux-arm64.tar.gz",
		Size:        1012554,
		ContentType: "application/x-gtar",
	}},
}
var testGithubRelease2 = stew.GithubRelease{
	TagName: "v0.0.1",
	Assets: []stew.GithubAsset{{
		Name:        "ppath-v0.0.1-darwin-amd64.tar.gz",
		DownloadURL: "https://github.com/marwanhawari/ppath/releases/download/v0.0.1/ppath-v0.0.1-darwin-amd64.tar.gz",
		Size:        1139366,
		ContentType: "application/x-gtar",
	}, {
		Name:        "ppath-v0.0.1-darwin-arm64.tar.gz",
		DownloadURL: "https://github.com/marwanhawari/ppath/releases/download/v0.0.1/ppath-v0.0.1-darwin-arm64.tar.gz",
		Size:        1101013,
		ContentType: "application/x-gtar",
	}, {
		Name:        "ppath-v0.0.1-linux-amd64.tar.gz",
		DownloadURL: "https://github.com/marwanhawari/ppath/releases/download/v0.0.1/ppath-v0.0.1-linux-amd64.tar.gz",
		Size:        1093728,
		ContentType: "application/x-gtar",
	}, {
		Name:        "ppath-v0.0.1-linux-arm64.tar.gz",
		DownloadURL: "https://github.com/marwanhawari/ppath/releases/download/v0.0.1/ppath-v0.0.1-linux-arm64.tar.gz",
		Size:        1014337,
		ContentType: "application/x-gtar",
	}},
}

var testGithubAPIResponse = stew.GithubAPIResponse{
	testGithubRelease0,
	testGithubRelease1,
	testGithubRelease2,
}

var testGithubProject = stew.GithubProject{
	Owner:    "marwanhawari",
	Repo:     "ppath",
	Releases: testGithubAPIResponse,
}

var testReleases = []string{
	"v0.0.3",
	"v0.0.2",
	"v0.0.1",
}

var testReleaseAssets = []string{
	"ppath-v0.0.1-darwin-amd64.tar.gz",
	"ppath-v0.0.1-darwin-arm64.tar.gz",
	"ppath-v0.0.1-linux-amd64.tar.gz",
	"ppath-v0.0.1-linux-arm64.tar.gz",
}

func TestNewGithubProject(t *testing.T) {
	testCases := []struct {
		name    string
		owner   string
		repo    string
		want    stew.GithubProject
		wantErr error
	}{{
		name:  "test1",
		owner: "marwanhawari",
		repo:  "ppath",
		want:  testGithubProject,
	}, {
		name:    "test2",
		owner:   "marwanhawari",
		repo:    "p",
		wantErr: stewhttp.NonZeroStatusCodeError{StatusCode: 404},
	}}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			rt := testsupport.NewDefaultRuntime(t)
			got, err := stew.NewGithubProject(rt.Config, tt.owner, tt.repo)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetGithubReleasesTags(t *testing.T) {
	type TestCase struct {
		name    string
		project stew.GithubProject
		want    []string
		wantErr error
	}
	testCases := []TestCase{{
		name:    "test1",
		project: testGithubProject,
		want:    testReleases,
	}, {
		name:    "test2",
		project: stew.GithubProject{},
		wantErr: stew.ReleasesNotFoundError{Owner: "", Repo: ""},
		want:    []string{},
	}}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := stew.GetGithubReleasesTags(tt.project)

			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetGithubReleasesAssets(t *testing.T) {
	type TestCase struct {
		name    string
		project stew.GithubProject
		tag     string
		want    []string
		wantErr error
	}
	cases := []TestCase{{
		name:    "test1",
		project: testGithubProject,
		tag:     "v0.0.1",
		want:    testReleaseAssets,
	}, {
		name:    "test2",
		project: testGithubProject,
		tag:     "v0.0.100",
		want:    []string{},
		wantErr: stew.AssetsNotFoundError{Tag: "v0.0.100"},
	}, {
		name:    "test3",
		project: stew.GithubProject{},
		tag:     "v0.0.1",
		want:    []string{},
		wantErr: stew.AssetsNotFoundError{Tag: "v0.0.1"},
	}}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := stew.GetGithubReleasesAssets(tt.project, tt.tag)

			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDetectAsset(t *testing.T) {
	type TestCase struct {
		name          string
		userOS        string
		userArch      string
		releaseAssets []string
		want          string
		wantErr       error
	}
	testCases := []TestCase{{
		name:          "test1",
		userOS:        "darwin",
		userArch:      "arm64",
		releaseAssets: testReleaseAssets,
		want:          "ppath-v0.0.1-darwin-arm64.tar.gz",
	}, {
		name:          "test2",
		userOS:        "darwin",
		userArch:      "amd64",
		releaseAssets: testReleaseAssets,
		want:          "ppath-v0.0.1-darwin-amd64.tar.gz",
	}, {
		name:          "test3",
		userOS:        "linux",
		userArch:      "arm64",
		releaseAssets: testReleaseAssets,
		want:          "ppath-v0.0.1-linux-arm64.tar.gz",
	}, {
		name:          "test4",
		userOS:        "linux",
		userArch:      "amd64",
		releaseAssets: testReleaseAssets,
		want:          "ppath-v0.0.1-linux-amd64.tar.gz",
	}, {
		name:          "test5",
		userOS:        "darwin",
		userArch:      "arm64",
		releaseAssets: testReleaseAssets,
		want:          "ppath-v0.0.1-darwin-arm64.tar.gz",
	}, {
		name:          "test6",
		userOS:        "windows",
		userArch:      "386",
		releaseAssets: append(testReleaseAssets, "ppath-v0.0.1-windows-386.tar.gz"),
		want:          "ppath-v0.0.1-windows-386.tar.gz",
	}, {
		name:          "test7",
		userOS:        "windows",
		userArch:      "unexpectedArch",
		releaseAssets: append(testReleaseAssets, "ppath-v0.0.1-windows-unexpectedArch.tar.gz"),
		want:          "ppath-v0.0.1-windows-unexpectedArch.tar.gz",
	}, {
		name:          "test8",
		userOS:        "darwin",
		userArch:      "arm64",
		releaseAssets: append(testReleaseAssets, "ppath-v0.0.1-darwin-x64.tar"),
		want:          "ppath-v0.0.1-darwin-arm64.tar.gz",
	}}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			rt := config.Runtime{
				OS:   tt.userOS,
				Arch: tt.userArch,
			}
			got, err := stew.DetectAsset(rt, tt.releaseAssets)

			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

var testGithubSearch = stew.GithubSearch{
	SearchQuery: "cardil/zsh",
	Count:       1,
	Items: []stew.GithubSearchResult{{
		FullName:    "cardil/zsh-defaults",
		Stars:       0,
		Language:    "Shell",
		Description: "Defaults to be used for all my ZSH installs",
	}},
}

var testFormattedSearchResults = []string{
	"cardil/zsh-defaults [⭐️0] Defaults to be used for all my ZSH installs",
}

func TestNewGithubSearch(t *testing.T) {
	testCases := []struct {
		query   string
		want    stew.GithubSearch
		wantErr error
	}{{
		query: "cardil/zsh",
		want:  testGithubSearch,
	}, {
		query:   "",
		want:    stew.GithubSearch{},
		wantErr: stew.InvalidGithubSearchQueryError{},
	}}
	for _, tt := range testCases {
		t.Run(tt.query, func(t *testing.T) {
			rt := testsupport.NewDefaultRuntime(t)
			got, err := stew.NewGithubSearch(rt.Config, tt.query)

			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFormatSearchResults(t *testing.T) {
	testCases := []struct {
		name   string
		search stew.GithubSearch
		want   []string
	}{{
		name:   "cardil/zsh",
		search: testGithubSearch,
		want:   testFormattedSearchResults,
	}}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got := stew.FormatSearchResults(tt.search)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestValidateGithubSearchQuery(t *testing.T) {
	type TestCase struct {
		name    string
		query   string
		wantErr error
	}
	testCases := []TestCase{{
		name:  "test1",
		query: "testQuery",
	}, {
		name:    "test2",
		query:   "^&*",
		wantErr: stew.InvalidGithubSearchQueryError{Query: "^&*"},
	}}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			err := stew.ValidateGithubSearchQuery(tt.query)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
