package stew

import (
	"errors"
	"fmt"
	"testing"

	"github.com/marwanhawari/stew/constants"
	"github.com/stretchr/testify/assert"
)

func TestNonZeroStatusCodeError_Error(t *testing.T) {
	assert := assert.New(t)
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
				StatusCode: 1,
			},
			want: fmt.Sprintf("%v Received non-zero status code from HTTP request: %v", constants.RedColor("Error:"), constants.RedColor(1)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NonZeroStatusCodeError{
				StatusCode: tt.fields.StatusCode,
			}
			assert.Equal(tt.want, e.Error())
		})
	}
}

func TestReleasesNotFoundError_Error(t *testing.T) {
	assert := assert.New(t)
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
			e := ReleasesNotFoundError{
				Owner: tt.fields.Owner,
				Repo:  tt.fields.Repo,
			}
			assert.Equal(tt.want, e.Error())
		})
	}
}

func TestAssetsNotFoundError_Error(t *testing.T) {
	assert := assert.New(t)
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
			e := AssetsNotFoundError{
				Tag: tt.fields.Tag,
			}
			assert.Equal(tt.want, e.Error())
		})
	}
}

func TestNoPackagesInLockfileError_Error(t *testing.T) {
	assert := assert.New(t)
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
			e := NoPackagesInLockfileError{}
			assert.Equal(tt.want, e.Error())
		})
	}
}

func TestIndexOutOfBoundsInLockfileError_Error(t *testing.T) {
	assert := assert.New(t)
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
			e := IndexOutOfBoundsInLockfileError{}
			assert.Equal(tt.want, e.Error())
		})
	}
}

func TestExitUserSelectionError_Error(t *testing.T) {
	assert := assert.New(t)
	type fields struct {
		Err error
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "test1",
			fields: fields{
				Err: errors.New("testErr"),
			},
			want: fmt.Sprintf("%v Exited from user selection: %v", constants.RedColor("Error:"), constants.RedColor(errors.New("testErr"))),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := ExitUserSelectionError{
				Err: tt.fields.Err,
			}
			assert.Equal(tt.want, e.Error())
		})
	}
}

func TestStewpathNotFoundError_Error(t *testing.T) {
	assert := assert.New(t)
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
			e := StewpathNotFoundError{
				StewPath: tt.fields.StewPath,
			}
			assert.Equal(tt.want, e.Error())
		})
	}
}

func TestNonZeroStatusCodeDownloadError_Error(t *testing.T) {
	assert := assert.New(t)
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
			e := NonZeroStatusCodeDownloadError{
				StatusCode: tt.fields.StatusCode,
			}
			assert.Equal(tt.want, e.Error())
		})
	}
}

func TestEmptyCLIInputError_Error(t *testing.T) {
	assert := assert.New(t)
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
			e := EmptyCLIInputError{}
			assert.Equal(tt.want, e.Error())
		})
	}
}

func TestCLIFlagAndInputError_Error(t *testing.T) {
	assert := assert.New(t)
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
			e := CLIFlagAndInputError{}
			assert.Equal(tt.want, e.Error())
		})
	}
}

func TestAssetAlreadyDownloadedError_Error(t *testing.T) {
	assert := assert.New(t)
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
			e := AssetAlreadyDownloadedError{
				Asset: tt.fields.Asset,
			}
			assert.Equal(tt.want, e.Error())
		})
	}
}

func TestAbortBinaryOverwriteError_Error(t *testing.T) {
	assert := assert.New(t)
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
			e := AbortBinaryOverwriteError{
				Binary: tt.fields.Binary,
			}
			assert.Equal(tt.want, e.Error())
		})
	}
}

func TestBinaryNotInstalledError_Error(t *testing.T) {
	assert := assert.New(t)
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
			e := BinaryNotInstalledError{
				Binary: tt.fields.Binary,
			}
			assert.Equal(tt.want, e.Error())
		})
	}
}

func TestNoBinariesInstalledError_Error(t *testing.T) {
	assert := assert.New(t)
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
			e := NoBinariesInstalledError{}
			assert.Equal(tt.want, e.Error())
		})
	}
}

func TestUnrecognizedInputError_Error(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name string
		want string
	}{
		{
			name: "test1",
			want: fmt.Sprintf("%v Input was not recognized as a URL or GitHub repo", constants.RedColor("Error:")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := UnrecognizedInputError{}
			assert.Equal(tt.want, e.Error())
		})
	}
}

func TestInstalledFromURLError_Error(t *testing.T) {
	assert := assert.New(t)
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
			e := InstalledFromURLError{
				Binary: tt.fields.Binary,
			}
			assert.Equal(tt.want, e.Error())
		})
	}
}

func TestAlreadyInstalledLatestTagError_Error(t *testing.T) {
	assert := assert.New(t)
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
			e := AlreadyInstalledLatestTagError{
				Tag: tt.fields.Tag,
			}
			assert.Equal(tt.want, e.Error())
		})
	}
}

func TestNoGithubSearchResultsError_Error(t *testing.T) {
	assert := assert.New(t)
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
			e := NoGithubSearchResultsError{
				SearchQuery: tt.fields.SearchQuery,
			}
			assert.Equal(tt.want, e.Error())
		})
	}
}

func TestInvalidGithubSearchQueryError_Error(t *testing.T) {
	assert := assert.New(t)
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
			want: fmt.Sprintf("%v The search query %v contains invalid characters", constants.RedColor("Error:"), constants.RedColor("testQuery")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := InvalidGithubSearchQueryError{
				SearchQuery: tt.fields.SearchQuery,
			}
			assert.Equal(tt.want, e.Error())
		})
	}
}
