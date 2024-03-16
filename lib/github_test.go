package stew

import (
	"reflect"
	"testing"
)

var testGithubRelease0 GithubRelease = GithubRelease{
	TagName: "v0.0.3",
	Assets: []GithubAsset{
		{
			Name:        "checksums.txt",
			DownloadURL: "https://api.github.com/repos/marwanhawari/ppath/releases/assets/52676090",
			Size:        394,
			ContentType: "text/plain; charset=utf-8",
		},
		{
			Name:        "ppath-v0.0.3-darwin-amd64.tar.gz",
			DownloadURL: "https://api.github.com/repos/marwanhawari/ppath/releases/assets/52676089",
			Size:        626448,
			ContentType: "application/gzip",
		},
		{
			Name:        "ppath-v0.0.3-darwin-arm64.tar.gz",
			DownloadURL: "https://api.github.com/repos/marwanhawari/ppath/releases/assets/52676091",
			Size:        625832,
			ContentType: "application/gzip",
		},
		{
			Name:        "ppath-v0.0.3-linux-amd64.tar.gz",
			DownloadURL: "https://api.github.com/repos/marwanhawari/ppath/releases/assets/52676093",
			Size:        567449,
			ContentType: "application/gzip",
		},
		{
			Name:        "ppath-v0.0.3-linux-arm64.tar.gz",
			DownloadURL: "https://api.github.com/repos/marwanhawari/ppath/releases/assets/52676092",
			Size:        529321,
			ContentType: "application/gzip",
		},
	},
}
var testGithubRelease1 GithubRelease = GithubRelease{
	TagName: "v0.0.2",
	Assets: []GithubAsset{
		{
			Name:        "ppath-v0.0.2-darwin-amd64.tar.gz",
			DownloadURL: "https://api.github.com/repos/marwanhawari/ppath/releases/assets/52647309",
			Size:        1139147,
			ContentType: "application/x-gtar",
		},
		{
			Name:        "ppath-v0.0.2-darwin-arm64.tar.gz",
			DownloadURL: "https://api.github.com/repos/marwanhawari/ppath/releases/assets/52647324",
			Size:        1100483,
			ContentType: "application/x-gtar",
		},
		{
			Name:        "ppath-v0.0.2-linux-amd64.tar.gz",
			DownloadURL: "https://api.github.com/repos/marwanhawari/ppath/releases/assets/52647307",
			Size:        1092421,
			ContentType: "application/x-gtar",
		},
		{
			Name:        "ppath-v0.0.2-linux-arm64.tar.gz",
			DownloadURL: "https://api.github.com/repos/marwanhawari/ppath/releases/assets/52647308",
			Size:        1012554,
			ContentType: "application/x-gtar",
		},
	},
}
var testGithubRelease2 GithubRelease = GithubRelease{
	TagName: "v0.0.1",
	Assets: []GithubAsset{
		{
			Name:        "ppath-v0.0.1-darwin-amd64.tar.gz",
			DownloadURL: "https://api.github.com/repos/marwanhawari/ppath/releases/assets/51111591",
			Size:        1139366,
			ContentType: "application/x-gtar",
		},
		{
			Name:        "ppath-v0.0.1-darwin-arm64.tar.gz",
			DownloadURL: "https://api.github.com/repos/marwanhawari/ppath/releases/assets/51111587",
			Size:        1101013,
			ContentType: "application/x-gtar",
		},
		{
			Name:        "ppath-v0.0.1-linux-amd64.tar.gz",
			DownloadURL: "https://api.github.com/repos/marwanhawari/ppath/releases/assets/51111593",
			Size:        1093728,
			ContentType: "application/x-gtar",
		},
		{
			Name:        "ppath-v0.0.1-linux-arm64.tar.gz",
			DownloadURL: "https://api.github.com/repos/marwanhawari/ppath/releases/assets/51111599",
			Size:        1014337,
			ContentType: "application/x-gtar",
		},
	},
}

var testGithubAPIResponse GithubAPIResponse = GithubAPIResponse{
	testGithubRelease0,
	testGithubRelease1,
	testGithubRelease2,
}

var testGithubProject GithubProject = GithubProject{
	Owner:    "marwanhawari",
	Repo:     "ppath",
	Releases: testGithubAPIResponse,
}

var testGithubJSON, _ = getGithubJSON("marwanhawari", "ppath")

var testReleases = []string{"v0.0.3", "v0.0.2", "v0.0.1"}

var testReleaseAssets = []string{"ppath-v0.0.1-darwin-amd64.tar.gz", "ppath-v0.0.1-darwin-arm64.tar.gz", "ppath-v0.0.1-linux-amd64.tar.gz", "ppath-v0.0.1-linux-arm64.tar.gz"}

var testDarwinAssets = []string{"ppath-v0.0.1-darwin-amd64.tar.gz", "ppath-v0.0.1-darwin-arm64.tar.gz"}

func Test_readGithubJSON(t *testing.T) {
	type args struct {
		jsonString string
	}
	tests := []struct {
		name    string
		args    args
		want    GithubAPIResponse
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				jsonString: testGithubJSON,
			},
			want:    testGithubAPIResponse,
			wantErr: false,
		},
		{
			name: "test2",
			args: args{
				jsonString: "",
			},
			want:    GithubAPIResponse{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readGithubJSON(tt.args.jsonString)
			if (err != nil) != tt.wantErr {
				t.Errorf("readGithubJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readGithubJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getGithubJSON(t *testing.T) {
	type args struct {
		owner string
		repo  string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				owner: "marwanhawari",
				repo:  "ppath",
			},
			want:    testGithubJSON,
			wantErr: false,
		},
		{
			name: "test2",
			args: args{
				owner: "marwanhawari",
				repo:  "p",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getGithubJSON(tt.args.owner, tt.args.repo)
			if (err != nil) != tt.wantErr {
				t.Errorf("getGithubJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getGithubJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewGithubProject(t *testing.T) {
	type args struct {
		owner string
		repo  string
	}
	tests := []struct {
		name    string
		args    args
		want    GithubProject
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				owner: "marwanhawari",
				repo:  "ppath",
			},
			want:    testGithubProject,
			wantErr: false,
		},
		{
			name: "test2",
			args: args{
				owner: "marwanhawari",
				repo:  "p",
			},
			want:    GithubProject{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewGithubProject(tt.args.owner, tt.args.repo)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewGithubProject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGithubProject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetGithubReleasesTags(t *testing.T) {
	type args struct {
		ghProject GithubProject
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				ghProject: testGithubProject,
			},
			want:    testReleases,
			wantErr: false,
		},
		{
			name: "test2",
			args: args{
				ghProject: GithubProject{},
			},
			want:    []string{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetGithubReleasesTags(tt.args.ghProject)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetGithubReleasesTags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetGithubReleasesTags() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_releasesFound(t *testing.T) {
	type args struct {
		releaseTags []string
		owner       string
		repo        string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				releaseTags: testReleases,
				owner:       "marwanhawari",
				repo:        "ppath",
			},
			wantErr: false,
		},
		{
			name: "test2",
			args: args{
				releaseTags: []string{},
				owner:       "marwanhawari",
				repo:        "ppath",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := releasesFound(tt.args.releaseTags, tt.args.owner, tt.args.repo); (err != nil) != tt.wantErr {
				t.Errorf("releasesFound() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetGithubReleasesAssets(t *testing.T) {
	type args struct {
		ghProject GithubProject
		tag       string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				ghProject: testGithubProject,
				tag:       "v0.0.1",
			},
			want:    testReleaseAssets,
			wantErr: false,
		},
		{
			name: "test2",
			args: args{
				ghProject: testGithubProject,
				tag:       "v0.0.100",
			},
			want:    []string{},
			wantErr: true,
		},
		{
			name: "test3",
			args: args{
				ghProject: GithubProject{},
				tag:       "v0.0.1",
			},
			want:    []string{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetGithubReleasesAssets(tt.args.ghProject, tt.args.tag)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetGithubReleasesAssets() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetGithubReleasesAssets() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_assetsFound(t *testing.T) {
	type args struct {
		releaseAssets []string
		releaseTag    string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				releaseAssets: testReleaseAssets,
				releaseTag:    "v0.0.1",
			},
			wantErr: false,
		},
		{
			name: "test2",
			args: args{
				releaseAssets: []string{},
				releaseTag:    "v0.0.1",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := assetsFound(tt.args.releaseAssets, tt.args.releaseTag); (err != nil) != tt.wantErr {
				t.Errorf("assetsFound() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDetectAsset(t *testing.T) {
	type args struct {
		userOS        string
		userArch      string
		releaseAssets []string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				userOS:        "darwin",
				userArch:      "arm64",
				releaseAssets: testReleaseAssets,
			},
			want:    "ppath-v0.0.1-darwin-arm64.tar.gz",
			wantErr: false,
		},
		{
			name: "test2",
			args: args{
				userOS:        "darwin",
				userArch:      "amd64",
				releaseAssets: testReleaseAssets,
			},
			want:    "ppath-v0.0.1-darwin-amd64.tar.gz",
			wantErr: false,
		},
		{
			name: "test3",
			args: args{
				userOS:        "linux",
				userArch:      "arm64",
				releaseAssets: testReleaseAssets,
			},
			want:    "ppath-v0.0.1-linux-arm64.tar.gz",
			wantErr: false,
		},
		{
			name: "test4",
			args: args{
				userOS:        "linux",
				userArch:      "amd64",
				releaseAssets: testReleaseAssets,
			},
			want:    "ppath-v0.0.1-linux-amd64.tar.gz",
			wantErr: false,
		},
		{
			name: "test5",
			args: args{
				userOS:        "darwin",
				userArch:      "arm64",
				releaseAssets: []string{},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "test6",
			args: args{
				userOS:        "windows",
				userArch:      "386",
				releaseAssets: append(testReleaseAssets, "ppath-v0.0.1-windows-386.tar.gz"),
			},
			want:    "ppath-v0.0.1-windows-386.tar.gz",
			wantErr: false,
		},
		{
			name: "test7",
			args: args{
				userOS:        "windows",
				userArch:      "unexpectedArch",
				releaseAssets: append(testReleaseAssets, "ppath-v0.0.1-windows-unexpectedArch.tar.gz"),
			},
			want:    "ppath-v0.0.1-windows-unexpectedArch.tar.gz",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DetectAsset(tt.args.userOS, tt.args.userArch, tt.args.releaseAssets)
			if (err != nil) != tt.wantErr {
				t.Errorf("DetectAsset() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DetectAsset() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_darwinARMFallback(t *testing.T) {
	type args struct {
		darwinAssets []string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				darwinAssets: testDarwinAssets,
			},
			want:    "ppath-v0.0.1-darwin-amd64.tar.gz",
			wantErr: false,
		},
		{
			name: "test2",
			args: args{
				darwinAssets: []string{},
			},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := darwinARMFallback(tt.args.darwinAssets)
			if (err != nil) != tt.wantErr {
				t.Errorf("darwinARMFallback() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("darwinARMFallback() = %v, want %v", got, tt.want)
			}
		})
	}
}

var testGithubSearchJSON, _ = getGithubSearchJSON("marwanhawari/ppath")

var testGithubSearchReadJSON GithubSearch = GithubSearch{
	Count: 1,
	Items: []GithubSearchResult{
		{
			FullName:    "marwanhawari/ppath",
			Stars:       7,
			Language:    "Go",
			Description: "🌈 A command-line tool to pretty print your system's PATH environment variable.",
		},
	},
}

var testGithubSearch GithubSearch = GithubSearch{
	SearchQuery: "marwanhawari/ppath",
	Count:       1,
	Items: []GithubSearchResult{
		{
			FullName:    "marwanhawari/ppath",
			Stars:       7,
			Language:    "Go",
			Description: "🌈 A command-line tool to pretty print your system's PATH environment variable.",
		},
	},
}

var testFormattedSearchResults = []string{"marwanhawari/ppath [⭐️7]  🌈 A command-line tool to pretty print your system's PATH environment variable."}

func Test_getGithubSearchJSON(t *testing.T) {
	type args struct {
		searchQuery string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				searchQuery: "marwanhawari/ppath",
			},
			want:    testGithubSearchJSON,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getGithubSearchJSON(tt.args.searchQuery)
			if (err != nil) != tt.wantErr {
				t.Errorf("getGithubSearchJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getGithubSearchJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_readGithubSearchJSON(t *testing.T) {
	type args struct {
		jsonString string
	}
	tests := []struct {
		name    string
		args    args
		want    GithubSearch
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				jsonString: testGithubSearchJSON,
			},
			want:    testGithubSearchReadJSON,
			wantErr: false,
		},
		{
			name: "test2",
			args: args{
				jsonString: "",
			},
			want:    GithubSearch{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readGithubSearchJSON(tt.args.jsonString)
			if (err != nil) != tt.wantErr {
				t.Errorf("readGithubSearchJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readGithubSearchJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewGithubSearch(t *testing.T) {
	type args struct {
		searchQuery string
	}
	tests := []struct {
		name    string
		args    args
		want    GithubSearch
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				searchQuery: "marwanhawari/ppath",
			},
			want:    testGithubSearch,
			wantErr: false,
		},
		{
			name: "test2",
			args: args{
				searchQuery: "",
			},
			want:    GithubSearch{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewGithubSearch(tt.args.searchQuery)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewGithubSearch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGithubSearch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatSearchResults(t *testing.T) {
	type args struct {
		ghSearch GithubSearch
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "test1",
			args: args{
				ghSearch: testGithubSearch,
			},
			want: testFormattedSearchResults,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatSearchResults(tt.args.ghSearch); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FormatSearchResults() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateGithubSearchQuery(t *testing.T) {
	type args struct {
		searchQuery string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				searchQuery: "testQuery",
			},
			wantErr: false,
		},
		{
			name: "test2",
			args: args{
				searchQuery: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateGithubSearchQuery(tt.args.searchQuery); (err != nil) != tt.wantErr {
				t.Errorf("ValidateGithubSearchQuery() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
