package stew

import (
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"runtime"
	"testing"
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			lockFilePath := path.Join(tempDir, "Stewfile.lock.json")
			WriteLockFileJSON(testLockfile, lockFilePath)

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
	}{
		{
			name:    "test1",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			lockFilePath := path.Join(tempDir, "Stewfile.lock.json")

			if err := WriteLockFileJSON(testLockfile, lockFilePath); (err != nil) != tt.wantErr {
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var testLockfilePackages []PackageData

			if tt.name != "test5" {
				testLockfilePackages = make([]PackageData, len(testLockfile.Packages))
				copy(testLockfilePackages, testLockfile.Packages)
			}

			got, err := RemovePackage(testLockfilePackages, tt.args.index)
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
	}{
		{
			name:    "test1",
			want:    testStewfileSlice,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tempDir := t.TempDir()
			testStewfilePath := path.Join(tempDir, "Stewfile")
			ioutil.WriteFile(testStewfilePath, []byte(testStewfileContents), 0644)

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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tempDir := t.TempDir()
			lockFilePath := path.Join(tempDir, "Stewfile.lock.json")
			WriteLockFileJSON(testLockfile, lockFilePath)

			got, err := NewLockFile(lockFilePath, tt.args.userOS, tt.args.userArch)
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tempDir := t.TempDir()
			lockFilePath := path.Join(tempDir, "Stewfile.lock.json")

			got, err := NewLockFile(lockFilePath, tt.args.userOS, tt.args.userArch)
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
	}{
		{
			name: "test1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			testSystemInfo := SystemInfo{
				Os:               runtime.GOOS,
				Arch:             runtime.GOARCH,
				StewPath:         tempDir,
				StewBinPath:      path.Join(tempDir, "bin"),
				StewPkgPath:      path.Join(tempDir, "pkg"),
				StewLockFilePath: path.Join(tempDir, "Stewfile.lock.json"),
				StewTmpPath:      path.Join(tempDir, "tmp"),
			}

			got := NewSystemInfo(tempDir)
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
	}{
		{
			name:    "test1",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			os.MkdirAll(path.Join(tempDir, "pkg"), 0755)
			testStewAssetPath := path.Join(tempDir, "pkg", "testAsset.tar.gz")
			os.MkdirAll(path.Join(tempDir, "bin"), 0755)
			testStewBinaryPath := path.Join(tempDir, "bin", "testBinary")

			ioutil.WriteFile(testStewAssetPath, []byte("This is a test asset"), 0644)
			ioutil.WriteFile(testStewBinaryPath, []byte("This is a test binary"), 0644)

			assetExists, _ := PathExists(testStewAssetPath)
			binaryExists, _ := PathExists(testStewBinaryPath)

			if !assetExists || !binaryExists {
				t.Errorf("Either the asset or the binary does not exist yet")
			}

			if err := DeleteAssetAndBinary(path.Dir(testStewAssetPath), path.Dir(testStewBinaryPath), path.Base(testStewAssetPath), path.Base(testStewBinaryPath)); (err != nil) != tt.wantErr {
				t.Errorf("DeleteAssetAndBinary() error = %v, wantErr %v", err, tt.wantErr)
			}

			assetExists, _ = PathExists(testStewAssetPath)
			binaryExists, _ = PathExists(testStewBinaryPath)

			if assetExists || binaryExists {
				t.Errorf("Either the binary or the asset still exists")
			}

		})
	}
}
