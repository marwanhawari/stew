package stew

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/marwanhawari/stew/lib/config"
	"github.com/marwanhawari/stew/lib/pathutil"
	"github.com/marwanhawari/stew/lib/ui/terminal"
	"github.com/stretchr/testify/require"
)

var testLockfile = LockFile{
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

var testStewfileContents = `junegunn/fzf@0.29.0
https://github.com/sharkdp/hyperfine/releases/download/v1.12.0/hyperfine-v1.12.0-x86_64-apple-darwin.tar.gz
marwanhawari/ppath@v0.0.3
`

var testStewfileSlice = []string{
	"junegunn/fzf@0.29.0",
	"https://github.com/sharkdp/hyperfine/releases/download/v1.12.0/hyperfine-v1.12.0-x86_64-apple-darwin.tar.gz",
	"marwanhawari/ppath@v0.0.3",
}

var testStewLockFileContents = `{
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

var testStewLockFileSlice = []string{
	"cli/cli@v2.4.0::gh_2.4.0_macOS_amd64.tar.gz!!gh",
	"junegunn/fzf@0.29.0::fzf-0.29.0-darwin_arm64.zip!!fzf",
	"https://github.com/sharkdp/hyperfine/releases/download/v1.12.0/hyperfine-v1.12.0-x86_64-apple-darwin.tar.gz",
}

func Test_readLockFileJSON(t *testing.T) {
	tests := []struct {
		name    string
		want    LockFile
		wantErr bool
	}{{
		name:    "test1",
		want:    testLockfile,
		wantErr: false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			ptr := terminal.TestTerminal{TestingT: t}
			lockFilePath := filepath.Join(tempDir, "Stewfile.lock.json")
			require.NoError(t, WriteLockFileJSON(ptr, testLockfile, lockFilePath))

			got, err := readLockFileJSON(lockFilePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("readLockFileJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readLockFileJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWriteLockFileJSON(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{{
		name:    "test1",
		wantErr: false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			ptr := terminal.TestTerminal{TestingT: t}
			lockFilePath := filepath.Join(tempDir, "Stewfile.lock.json")

			if err := WriteLockFileJSON(ptr, testLockfile, lockFilePath); (err != nil) != tt.wantErr {
				t.Errorf("WriteLockFileJSON() error = %v, wantErr %v", err, tt.wantErr)
			}

			got, _ := readLockFileJSON(lockFilePath)

			if !reflect.DeepEqual(got, testLockfile) {
				t.Errorf("WriteLockFileJSON() = %v, want %v", got, testLockfile)
			}

		})
	}
}

func TestRemovePackage(t *testing.T) {
	tests := []struct {
		name    string
		index   int
		want    []PackageData
		wantErr bool
	}{{
		name:    "test1",
		index:   0,
		want:    testLockfile.Packages[1:],
		wantErr: false,
	}, {
		name:    "test2",
		index:   1,
		want:    []PackageData{testLockfile.Packages[0], testLockfile.Packages[2]},
		wantErr: false,
	}, {
		name:    "test3",
		index:   2,
		want:    testLockfile.Packages[:2],
		wantErr: false,
	}, {
		name:    "test4",
		index:   -1,
		want:    []PackageData{},
		wantErr: true,
	}, {
		name:    "test5",
		index:   0,
		want:    []PackageData{},
		wantErr: true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var testLockfilePackages []PackageData

			if tt.name != "test5" {
				testLockfilePackages = make([]PackageData, len(testLockfile.Packages))
				copy(testLockfilePackages, testLockfile.Packages)
			}

			got, err := RemovePackage(testLockfilePackages, tt.index)
			if (err != nil) != tt.wantErr {
				t.Errorf("RemovePackage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RemovePackage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadStewfileContents(t *testing.T) {
	tests := []struct {
		name    string
		want    []string
		wantErr bool
	}{{
		name:    "test1",
		want:    testStewfileSlice,
		wantErr: false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tempDir := t.TempDir()
			testStewfilePath := filepath.Join(tempDir, "Stewfile")
			require.NoError(t, os.WriteFile(testStewfilePath, []byte(testStewfileContents), 0o644))

			got, err := ReadStewfileContents(testStewfilePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadStewfileContents() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadStewfileContents() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadStewLockFileContents(t *testing.T) {
	tests := []struct {
		name string
		want []string
	}{{
		name: "test1",
		want: testStewLockFileSlice,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			testStewLockFilePath := filepath.Join(tempDir, "Stewfile.lock.json")
			if err := os.WriteFile(testStewLockFilePath, []byte(testStewLockFileContents), 0o644); err != nil {
				t.Fatalf("unexpected error: %+v", err)
			}

			got, err := ReadStewLockFileContents(testStewLockFilePath)
			if err != nil {
				t.Fatalf("unexpected error = %+v", err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("result\n"+
					" got = %+v\n"+
					"want = %+v", got, tt.want)
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
	}{{
		name: "test1",
		args: args{
			userOS:   testLockfile.Os,
			userArch: testLockfile.Arch,
		},
		want:    testLockfile,
		wantErr: false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			ptr := terminal.TestTerminal{TestingT: t}
			lockFilePath := filepath.Join(tempDir, "Stewfile.lock.json")
			if err := WriteLockFileJSON(ptr, testLockfile, lockFilePath); err != nil {
				t.Fatalf("unexpected error: %+v", err)
			}

			run := config.Runtime{
				OS:   tt.args.userOS,
				Arch: tt.args.userArch,
				SystemInfo: config.SystemInfo{
					LockPath: lockFilePath,
				},
			}
			got, err := NewLockFile(run)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewLockFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewLockFile() = %v, want %v", got, tt.want)
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
	}{{
		name: "test1",
		args: args{
			userOS:   testLockfile.Os,
			userArch: testLockfile.Arch,
		},
		wantErr: false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tempDir := t.TempDir()
			lockFilePath := filepath.Join(tempDir, "Stewfile.lock.json")

			run := config.Runtime{
				OS:   tt.args.userOS,
				Arch: tt.args.userArch,
				SystemInfo: config.SystemInfo{
					LockPath: lockFilePath,
				},
			}
			got, err := NewLockFile(run)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewLockFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			want := LockFile{Os: tt.args.userOS, Arch: tt.args.userArch, Packages: []PackageData{}}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("NewLockFile() = %v, want %v", got, want)
			}
		})
	}
}

func TestNewSystemInfo(t *testing.T) {
	tests := []struct {
		name string
	}{{
		name: "test1",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()

			testStewConfig := config.Config{
				StewPath:    tempDir,
				StewBinPath: filepath.Join(tempDir, "bin"),
			}

			testSystemInfo := config.SystemInfo{
				PkgPath:  filepath.Join(tempDir, "pkg"),
				LockPath: filepath.Join(tempDir, "Stewfile.lock.json"),
				TmpPath:  filepath.Join(tempDir, "tmp"),
			}

			got := config.NewSystemInfo(testStewConfig)
			if !reflect.DeepEqual(got, testSystemInfo) {
				t.Errorf("NewSystemInfo() = %v, want %v", got, testSystemInfo)
			}
		})
	}
}

func TestDeleteAssetAndBinary(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{{
		name:    "test1",
		wantErr: false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			require.NoError(t, os.MkdirAll(filepath.Join(tempDir, "pkg"), 0o755))
			testStewAssetPath := filepath.Join(tempDir, "pkg", "testAsset.tar.gz")
			require.NoError(t, os.MkdirAll(filepath.Join(tempDir, "bin"), 0o755))
			testStewBinaryPath := filepath.Join(tempDir, "bin", "testBinary")

			require.NoError(t, os.WriteFile(testStewAssetPath, []byte("This is a test asset"), 0o644))
			require.NoError(t, os.WriteFile(testStewBinaryPath, []byte("This is a test binary"), 0o644))

			assetExists, _ := pathutil.Exists(testStewAssetPath)
			binaryExists, _ := pathutil.Exists(testStewBinaryPath)

			if !assetExists || !binaryExists {
				t.Errorf("Either the asset or the binary does not exist yet")
			}

			if err := DeleteAssetAndBinary(filepath.Dir(testStewAssetPath), filepath.Dir(testStewBinaryPath), filepath.Base(testStewAssetPath), filepath.Base(testStewBinaryPath)); (err != nil) != tt.wantErr {
				t.Errorf("DeleteAssetAndBinary() error = %v, wantErr %v", err, tt.wantErr)
			}

			assetExists, _ = pathutil.Exists(testStewAssetPath)
			binaryExists, _ = pathutil.Exists(testStewBinaryPath)

			if assetExists || binaryExists {
				t.Errorf("Either the binary or the asset still exists")
			}

		})
	}
}
