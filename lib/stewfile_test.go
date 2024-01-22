package stew

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testLockfile LockFile = LockFile{
	Os:   "darwin",
	Arch: "arm64",
	Packages: []PackageData{
		{
			Source: "github",
			Owner:  "junegunn",
			Repo:   "fzf",
			Tag:    "0.29.0",
			Asset:  "fzf-0.29.0-darwin_arm64.zip",
			Binary: "fzf",
			URL:    "https://github.com/junegunn/fzf/releases/download/0.29.0/fzf-0.29.0-darwin_arm64.zip",
		},
		{
			Source: "other",
			Owner:  "",
			Repo:   "",
			Tag:    "",
			Asset:  "hyperfine-v1.12.0-x86_64-apple-darwin.tar.gz",
			Binary: "hyperfine",
			URL:    "https://github.com/sharkdp/hyperfine/releases/download/v1.12.0/hyperfine-v1.12.0-x86_64-apple-darwin.tar.gz",
		},
		{
			Source: "github",
			Owner:  "marwanhawari",
			Repo:   "ppath",
			Tag:    "v0.0.3",
			Asset:  "ppath-v0.0.3-darwin-arm64.tar.gz",
			Binary: "ppath",
			URL:    "https://github.com/marwanhawari/ppath/releases/download/v0.0.3/ppath-v0.0.3-darwin-arm64.tar.gz",
		},
	},
}

var testStewfileContents string = `junegunn/fzf@0.29.0
https://github.com/sharkdp/hyperfine/releases/download/v1.12.0/hyperfine-v1.12.0-x86_64-apple-darwin.tar.gz
marwanhawari/ppath@v0.0.3
`

var testStewfileSlice []string = []string{
	"junegunn/fzf@0.29.0",
	"https://github.com/sharkdp/hyperfine/releases/download/v1.12.0/hyperfine-v1.12.0-x86_64-apple-darwin.tar.gz",
	"marwanhawari/ppath@v0.0.3",
}

var testStewLockFileContents string = `{
	"os": "darwin",
	"arch": "arm64",
	"packages": [
	{
		"source": "github",
		"owner": "cli",
		"repo": "cli",
		"tag": "v2.4.0",
		"asset": "gh_2.4.0_macOS_amd64.tar.gz",
		"binary": "gh",
		"url": "https://github.com/cli/cli/releases/download/v2.4.0/gh_2.4.0_macOS_amd64.tar.gz"
	},
	{
		"source": "github",
		"owner": "junegunn",
		"repo": "fzf",
		"tag": "0.29.0",
		"asset": "fzf-0.29.0-darwin_arm64.zip",
		"binary": "fzf",
		"url": "https://github.com/junegunn/fzf/releases/download/0.29.0/fzf-0.29.0-darwin_arm64.zip"
	},
	{
		"source": "other",
		"owner": "",
		"repo": "",
		"tag": "",
		"asset": "hyperfine-v1.12.0-x86_64-apple-darwin.tar.gz",
		"binary": "hyperfine",
		"url": "https://github.com/sharkdp/hyperfine/releases/download/v1.12.0/hyperfine-v1.12.0-x86_64-apple-darwin.tar.gz"
	}
	]
}
`

var testStewLockFileSlice []string = []string{
	"cli/cli@v2.4.0::gh_2.4.0_macOS_amd64.tar.gz",
	"junegunn/fzf@0.29.0::fzf-0.29.0-darwin_arm64.zip",
	"https://github.com/sharkdp/hyperfine/releases/download/v1.12.0/hyperfine-v1.12.0-x86_64-apple-darwin.tar.gz",
}

func Test_readLockFileJSON(t *testing.T) {
	tests := []struct {
		name    string
		want    LockFile
		wantErr bool
	}{
		{
			name:    "test1",
			want:    testLockfile,
			wantErr: false,
		},
	}
	assert := assert.New(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			lockFilePath := filepath.Join(tempDir, "Stewfile.lock.json")
			WriteLockFileJSON(testLockfile, lockFilePath)

			got, err := readLockFileJSON(lockFilePath)
			if tt.wantErr {
				assert.Error(err, "Expected an error")
			} else {
				assert.NoError(err, "Did not expect an error")
				assert.Equal(tt.want, got, "Mismatch in lock file content")
			}
		})
	}
}

func TestWriteLockFileJSON(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "test1",
			wantErr: false,
		},
	}
	assert := assert.New(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			lockFilePath := filepath.Join(tempDir, "Stewfile.lock.json")

			err := WriteLockFileJSON(testLockfile, lockFilePath)
			if tt.wantErr {
				assert.Error(err, "Expected an error")
			} else {
				assert.NoError(err, "Did not expect an error")

				got, _ := readLockFileJSON(lockFilePath)
				assert.Equal(testLockfile, got, "Mismatch in written and read lock file")
			}
		})
	}
}

func TestRemovePackage(t *testing.T) {
	type args struct {
		index int
	}
	tests := []struct {
		name    string
		args    args
		want    []PackageData
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				index: 0,
			},
			want:    testLockfile.Packages[1:],
			wantErr: false,
		},
		{
			name: "test2",
			args: args{
				index: 1,
			},
			want:    []PackageData{testLockfile.Packages[0], testLockfile.Packages[2]},
			wantErr: false,
		},
		{
			name: "test3",
			args: args{
				index: 2,
			},
			want:    testLockfile.Packages[:2],
			wantErr: false,
		},
		{
			name: "test4",
			args: args{
				index: -1,
			},
			want:    []PackageData{},
			wantErr: true,
		},
		{
			name: "test5",
			args: args{
				index: 0,
			},
			want:    []PackageData{},
			wantErr: true,
		},
	}
	assert := assert.New(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var testLockfilePackages []PackageData

			if tt.name != "test5" {
				testLockfilePackages = make([]PackageData, len(testLockfile.Packages))
				copy(testLockfilePackages, testLockfile.Packages)
			}

			got, err := RemovePackage(testLockfilePackages, tt.args.index)
			if tt.wantErr {
				assert.Error(err, "Expected an error")
			} else {
				assert.NoError(err, "Did not expect an error")
				assert.Equal(tt.want, got, "Mismatch in package data")
			}
		})
	}
}

func TestReadStewfileContents(t *testing.T) {
	tests := []struct {
		name    string
		want    []string
		wantErr bool
	}{
		{
			name:    "test1",
			want:    testStewfileSlice,
			wantErr: false,
		},
	}
	assert := assert.New(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			testStewfilePath := filepath.Join(tempDir, "Stewfile")
			os.WriteFile(testStewfilePath, []byte(testStewfileContents), 0644)

			got, err := ReadStewfileContents(testStewfilePath)
			if tt.wantErr {
				assert.Error(err, "Expected an error")
			} else {
				assert.NoError(err, "Did not expect an error")
				assert.Equal(tt.want, got, "Mismatch in stewfile contents")
			}
		})
	}
}

func TestReadStewLockFileContents(t *testing.T) {
	tests := []struct {
		name    string
		want    []string
		wantErr bool
	}{
		{
			name:    "test1",
			want:    testStewLockFileSlice,
			wantErr: false,
		},
	}
	assert := assert.New(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			testStewLockFilePath := filepath.Join(tempDir, "Stewfile.lock.json")
			os.WriteFile(testStewLockFilePath, []byte(testStewLockFileContents), 0644)

			got, err := ReadStewLockFileContents(testStewLockFilePath)
			if tt.wantErr {
				assert.Error(err, "Expected an error")
			} else {
				assert.NoError(err, "Did not expect an error")
				assert.Equal(tt.want, got, "Mismatch in stewlock file contents")
			}
		})
	}
}

func TestNewLockFile(t *testing.T) {
	type args struct {
		userOS   string
		userArch string
	}
	tests := []struct {
		name    string
		args    args
		want    LockFile
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				userOS:   testLockfile.Os,
				userArch: testLockfile.Arch,
			},
			want:    testLockfile,
			wantErr: false,
		},
	}
	assert := assert.New(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			lockFilePath := filepath.Join(tempDir, "Stewfile.lock.json")
			WriteLockFileJSON(testLockfile, lockFilePath)

			got, err := NewLockFile(lockFilePath, tt.args.userOS, tt.args.userArch)
			if tt.wantErr {
				assert.Error(err, "Expected an error")
			} else {
				assert.NoError(err, "Did not expect an error")
				assert.Equal(tt.want, got, "Mismatch in lock file creation")
			}
		})
	}
}

func TestNewLockFileDoesntExist(t *testing.T) {
	type args struct {
		userOS   string
		userArch string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				userOS:   testLockfile.Os,
				userArch: testLockfile.Arch,
			},
			wantErr: false,
		},
	}
	assert := assert.New(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			lockFilePath := filepath.Join(tempDir, "Stewfile.lock.json")

			got, err := NewLockFile(lockFilePath, tt.args.userOS, tt.args.userArch)
			if tt.wantErr {
				assert.Error(err, "Expected an error")
			} else {
				assert.NoError(err, "Did not expect an error")
				want := LockFile{Os: tt.args.userOS, Arch: tt.args.userArch, Packages: []PackageData{}}
				assert.Equal(want, got, "Mismatch in new lock file creation")
			}
		})
	}
}

func TestNewSystemInfo(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "test1",
		},
	}
	assert := assert.New(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()

			testStewConfig := StewConfig{
				StewPath:    tempDir,
				StewBinPath: filepath.Join(tempDir, "bin"),
			}

			testSystemInfo := SystemInfo{
				StewPath:         tempDir,
				StewBinPath:      filepath.Join(tempDir, "bin"),
				StewPkgPath:      filepath.Join(tempDir, "pkg"),
				StewLockFilePath: filepath.Join(tempDir, "Stewfile.lock.json"),
				StewTmpPath:      filepath.Join(tempDir, "tmp"),
			}

			got := NewSystemInfo(testStewConfig)
			assert.Equal(testSystemInfo, got, "Mismatch in system info creation")
		})
	}
}

func TestDeleteAssetAndBinary(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "test1",
			wantErr: false,
		},
	}
	assert := assert.New(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			os.MkdirAll(filepath.Join(tempDir, "pkg"), 0755)
			testStewAssetPath := filepath.Join(tempDir, "pkg", "testAsset.tar.gz")
			os.MkdirAll(filepath.Join(tempDir, "bin"), 0755)
			testStewBinaryPath := filepath.Join(tempDir, "bin", "testBinary")

			os.WriteFile(testStewAssetPath, []byte("This is a test asset"), 0644)
			os.WriteFile(testStewBinaryPath, []byte("This is a test binary"), 0644)

			assetExists, _ := PathExists(testStewAssetPath)
			binaryExists, _ := PathExists(testStewBinaryPath)
			assert.True(assetExists && binaryExists, "Either the asset or the binary does not exist yet")

			err := DeleteAssetAndBinary(filepath.Dir(testStewAssetPath), filepath.Dir(testStewBinaryPath), filepath.Base(testStewAssetPath), filepath.Base(testStewBinaryPath))
			if tt.wantErr {
				assert.Error(err, "Expected an error")
			} else {
				assert.NoError(err, "Did not expect an error")

				assetExists, _ = PathExists(testStewAssetPath)
				binaryExists, _ = PathExists(testStewBinaryPath)
				assert.False(assetExists || binaryExists, "Either the binary or the asset still exists")
			}
		})
	}
}
