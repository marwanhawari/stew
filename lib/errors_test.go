package stew_test

import (
	"fmt"
	"testing"

	"github.com/marwanhawari/stew/constants"
	stew "github.com/marwanhawari/stew/lib"
)

func TestReleasesNotFoundError_Error(t *testing.T) {
	type fields struct {
		Owner string
		Repo  string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "test1",
			fields: fields{
				Owner: "testOwner",
				Repo:  "testRepo",
			},
			want: fmt.Sprintf("%v Could not find any releases for %v", constants.RedColor("Error:"), constants.RedColor("https://github.com/testOwner/testRepo")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := stew.ReleasesNotFoundError{
				Owner: tt.fields.Owner,
				Repo:  tt.fields.Repo,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("ReleasesNotFoundError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAssetsNotFoundError_Error(t *testing.T) {
	type fields struct {
		Tag string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "test1",
			fields: fields{
				Tag: "testTag",
			},
			want: fmt.Sprintf("%v Could not find any assets for release %v", constants.RedColor("Error:"), constants.RedColor("testTag")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := stew.AssetsNotFoundError{
				Tag: tt.fields.Tag,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("AssetsNotFoundError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNoPackagesInLockfileError_Error(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "test1",
			want: fmt.Sprintf("%v Cannot remove from an empty packages slice in the lockfile", constants.RedColor("Error:")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := stew.NoPackagesInLockfileError{}
			if got := e.Error(); got != tt.want {
				t.Errorf("NoPackagesInLockfileError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIndexOutOfBoundsInLockfileError_Error(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "test1",
			want: fmt.Sprintf("%v Index out of bounds in lockfile packages", constants.RedColor("Error:")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := stew.IndexOutOfBoundsInLockfileError{}
			if got := e.Error(); got != tt.want {
				t.Errorf("IndexOutOfBoundsInLockfileError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStewpathNotFoundError_Error(t *testing.T) {
	type fields struct {
		StewPath string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "test1",
			fields: fields{
				StewPath: "$HOME/.stew",
			},
			want: fmt.Sprintf("%v Could not find the stew path at %v", constants.RedColor("Error:"), constants.RedColor("$HOME/.stew")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := stew.StewpathNotFoundError{
				StewPath: tt.fields.StewPath,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("StewpathNotFoundError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNonZeroStatusCodeDownloadError_Error(t *testing.T) {
	type fields struct {
		StatusCode int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "test1",
			fields: fields{
				StatusCode: 404,
			},
			want: fmt.Sprintf("%v Received non-zero status code from HTTP request when attempting to download a file: %v", constants.RedColor("Error:"), constants.RedColor(404)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := stew.NonZeroStatusCodeDownloadError{
				StatusCode: tt.fields.StatusCode,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("NonZeroStatusCodeDownloadError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEmptyCLIInputError_Error(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "test1",
			want: fmt.Sprintf("%v Input cannot be empty. Use the --help flag for more info", constants.RedColor("Error:")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := stew.EmptyCLIInputError{}
			if got := e.Error(); got != tt.want {
				t.Errorf("EmptyCLIInputError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCLIFlagAndInputError_Error(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "test1",
			want: fmt.Sprintf("%v Cannot use the --all flag with a positional argument", constants.RedColor("Error:")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := stew.CLIFlagAndInputError{}
			if got := e.Error(); got != tt.want {
				t.Errorf("CLIFlagAndInputError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAssetAlreadyDownloadedError_Error(t *testing.T) {
	type fields struct {
		Asset string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "test1",
			fields: fields{
				Asset: "testAsset",
			},
			want: fmt.Sprintf("%v The %v asset has already been downloaded and installed", constants.RedColor("Error:"), constants.RedColor("testAsset")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := stew.AssetAlreadyDownloadedError{
				Asset: tt.fields.Asset,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("AssetAlreadyDownloadedError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAbortBinaryOverwriteError_Error(t *testing.T) {
	type fields struct {
		Binary string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "test1",
			fields: fields{
				Binary: "testBinary",
			},
			want: fmt.Sprintf("%v Overwrite of %v aborted", constants.RedColor("Error:"), constants.RedColor("testBinary")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := stew.AbortBinaryOverwriteError{
				Binary: tt.fields.Binary,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("AbortBinaryOverwriteError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBinaryNotInstalledError_Error(t *testing.T) {
	type fields struct {
		Binary string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "test1",
			fields: fields{
				Binary: "testBinary",
			},
			want: fmt.Sprintf("%v The binary %v is not currently installed", constants.RedColor("Error:"), constants.RedColor("testBinary")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := stew.BinaryNotInstalledError{
				Binary: tt.fields.Binary,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("BinaryNotInstalledError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNoBinariesInstalledError_Error(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "test1",
			want: fmt.Sprintf("%v No binaries are currently installed", constants.RedColor("Error:")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := stew.NoBinariesInstalledError{}
			if got := e.Error(); got != tt.want {
				t.Errorf("NoBinariesInstalledError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnrecognizedInputError_Error(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "test1",
			want: fmt.Sprintf("%v Input was not recognized as a URL or GitHub repo: %v",
				constants.RedColor("Error:"),
				constants.RedColor("")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := stew.UnrecognizedInputError{}
			if got := e.Error(); got != tt.want {
				t.Errorf("UnrecognizedInputError.Error()\n"+
					" got = %v\n"+
					"want = %v", got, tt.want)
			}
		})
	}
}

func TestInstalledFromURLError_Error(t *testing.T) {
	type fields struct {
		Binary string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "test1",
			fields: fields{
				Binary: "testBinary",
			},
			want: fmt.Sprintf("%v The %v binary was installed directly from a URL", constants.RedColor("Error:"), constants.RedColor("testBinary")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := stew.InstalledFromURLError{
				Binary: tt.fields.Binary,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("InstalledFromURLError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAlreadyInstalledLatestTagError_Error(t *testing.T) {
	type fields struct {
		Tag string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "test1",
			fields: fields{
				Tag: "testTag",
			},
			want: fmt.Sprintf("%v The latest tag %v is already installed", constants.RedColor("Error:"), constants.RedColor("testTag")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := stew.AlreadyInstalledLatestTagError{
				Tag: tt.fields.Tag,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("AlreadyInstalledLatestTagError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNoGithubSearchResultsError_Error(t *testing.T) {
	type fields struct {
		SearchQuery string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "test1",
			fields: fields{
				SearchQuery: "testQuery",
			},
			want: fmt.Sprintf("%v No GitHub search results found for search query %v", constants.RedColor("Error:"), constants.RedColor("testQuery")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := stew.NoGithubSearchResultsError{
				SearchQuery: tt.fields.SearchQuery,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("NoGithubSearchResultsError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInvalidGithubSearchQueryError_Error(t *testing.T) {
	type fields struct {
		SearchQuery string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "test1",
			fields: fields{
				SearchQuery: "testQuery",
			},
			want: fmt.Sprintf(`%v The search query "%v" contains invalid characters`,
				constants.RedColor("Error:"), constants.RedColor("testQuery")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := stew.InvalidGithubSearchQueryError{
				Query: tt.fields.SearchQuery,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("InvalidGithubSearchQueryError.Error()\n"+
					" got = %v\n"+
					"want = %v", got, tt.want)
			}
		})
	}
}
