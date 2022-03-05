package stew

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

func Test_isArchiveFile(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test1",
			args: args{
				filePath: "notArchive",
			},
			want: false,
		},
		{
			name: "test1",
			args: args{
				filePath: "Archive.tar.gz",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isArchiveFile(tt.args.filePath); got != tt.want {
				t.Errorf("isArchiveFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isExecutableFile(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				filePath: "testExecutableFile",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "test1",
			args: args{
				filePath: "notExecutableFile",
			},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			testExecutableFilePath := filepath.Join(tempDir, tt.args.filePath)
			if tt.want {
				ioutil.WriteFile(testExecutableFilePath, []byte("An executable file"), 0755)
			} else {
				ioutil.WriteFile(testExecutableFilePath, []byte("Not an executable file"), 0644)
			}

			got, err := isExecutableFile(testExecutableFilePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("isExecutableFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("isExecutableFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPathExists(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				path: "testFile",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "test2",
			args: args{
				path: "noFile",
			},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			testFilePath := filepath.Join(tempDir, tt.args.path)
			if tt.want {
				ioutil.WriteFile(testFilePath, []byte("A test file"), 0644)
			}

			got, err := PathExists(testFilePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("PathExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PathExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetStewPath(t *testing.T) {
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
			homeDir, _ := os.UserHomeDir()
			testStewPath := filepath.Join(homeDir, ".stew")
			stewPathExists, _ := PathExists(testStewPath)
			if !stewPathExists {
				os.MkdirAll(testStewPath, 0755)
			}

			got, err := GetStewPath()
			if !stewPathExists {
				os.RemoveAll(testStewPath)
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("GetStewPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != testStewPath {
				t.Errorf("GetStewPath() = %v, want %v", got, testStewPath)
			}
		})
	}
}

func TestDownloadFile(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				url: "https://github.com/marwanhawari/ppath/releases/download/v0.0.3/ppath-v0.0.3-darwin-arm64.tar.gz",
			},
			wantErr: false,
		},
		{
			name: "test1",
			args: args{
				url: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			testDownloadPath := filepath.Join(tempDir, filepath.Base(tt.args.url))
			if err := DownloadFile(testDownloadPath, tt.args.url); (err != nil) != tt.wantErr {
				t.Errorf("DownloadFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if fileExists, _ := PathExists(testDownloadPath); !fileExists {
				t.Errorf("The file %v was not found", testDownloadPath)
			}

		})
	}
}

func Test_copyFile(t *testing.T) {
	type args struct {
		srcFile  string
		destFile string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				srcFile:  "sourceFile.txt",
				destFile: "destFile.txt",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			srcFilePath := filepath.Join(tempDir, tt.args.srcFile)
			destFilePath := filepath.Join(tempDir, tt.args.destFile)

			ioutil.WriteFile(srcFilePath, []byte("A test file"), 0644)

			srcExists, _ := PathExists(srcFilePath)
			destExists, _ := PathExists(destFilePath)

			if !srcExists {
				t.Errorf("Source file %v not found", srcFilePath)
			}

			if destExists {
				t.Errorf("Dest file %v already exists", destFilePath)
			}

			if err := copyFile(srcFilePath, destFilePath); (err != nil) != tt.wantErr {
				t.Errorf("copyFile() error = %v, wantErr %v", err, tt.wantErr)
			}

			srcExists, _ = PathExists(srcFilePath)
			destExists, _ = PathExists(destFilePath)

			if !srcExists || !destExists {
				t.Errorf("Copy failed - src or dest file does not exist")
			}

		})
	}
}

func Test_walkDir(t *testing.T) {
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

			ioutil.WriteFile(filepath.Join(tempDir, "testFile.txt"), []byte("A test file"), 0644)
			os.MkdirAll(filepath.Join(tempDir, "bin"), 0755)
			ioutil.WriteFile(filepath.Join(tempDir, "bin", "binDirTestFile.txt"), []byte("Another test file"), 0644)

			want := []string{filepath.Join(tempDir, "bin", "binDirTestFile.txt"), filepath.Join(tempDir, "testFile.txt")}

			got, err := walkDir(tempDir)
			if (err != nil) != tt.wantErr {
				t.Errorf("walkDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("walkDir() = %v, want %v", got, want)
			}
		})
	}
}

func Test_getBinary(t *testing.T) {
	type args struct {
		repo string
	}
	tests := []struct {
		name       string
		args       args
		binaryName string
		wantErr    bool
	}{
		{
			name: "test1",
			args: args{
				repo: "testBinary",
			},
			binaryName: "testBinary",
			wantErr:    false,
		},
		{
			name: "test2",
			args: args{
				repo: "someRepo",
			},
			binaryName: "testBinary",
			wantErr:    false,
		},
		{
			name: "test3",
			args: args{
				repo: "someRepo",
			},
			binaryName: "testBinary.exe",
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()

			testBinaryFilePath := filepath.Join(tempDir, tt.binaryName)
			ioutil.WriteFile(testBinaryFilePath, []byte("An executable file"), 0755)
			testNonBinaryFilePath := filepath.Join(tempDir, "testNonBinary")
			ioutil.WriteFile(testNonBinaryFilePath, []byte("Not an executable file"), 0644)

			testFilePaths := []string{testBinaryFilePath, testNonBinaryFilePath}

			wantBinaryFile := filepath.Join(tempDir, tt.binaryName)
			wantBinaryName := filepath.Base(wantBinaryFile)

			got, got1, err := getBinary(testFilePaths, tt.args.repo)
			if (err != nil) != tt.wantErr {
				t.Errorf("getBinary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != wantBinaryFile {
				t.Errorf("getBinary() got = %v, want %v", got, wantBinaryFile)
			}
			if got1 != wantBinaryName {
				t.Errorf("getBinary() got1 = %v, want %v", got1, wantBinaryName)
			}
		})
	}
}

func Test_getBinaryError(t *testing.T) {
	type args struct {
		repo string
	}
	tests := []struct {
		name       string
		args       args
		binaryName string
		wantErr    bool
	}{
		{
			name: "test1",
			args: args{
				repo: "testBinary",
			},
			binaryName: "testBinary",
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()

			testNonBinaryFilePath := filepath.Join(tempDir, "testNonBinary")
			ioutil.WriteFile(testNonBinaryFilePath, []byte("Not an executable file"), 0644)

			testFilePaths := []string{testNonBinaryFilePath}

			wantBinaryFile := ""
			wantBinaryName := ""

			got, got1, err := getBinary(testFilePaths, tt.args.repo)
			if (err != nil) != tt.wantErr {
				t.Errorf("getBinary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != wantBinaryFile {
				t.Errorf("getBinary() got = %v, want %v", got, wantBinaryFile)
			}
			if got1 != wantBinaryName {
				t.Errorf("getBinary() got1 = %v, want %v", got1, wantBinaryName)
			}
		})
	}
}

func TestValidateCLIInput(t *testing.T) {
	type args struct {
		cliInput string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				cliInput: "",
			},
			wantErr: true,
		},
		{
			name: "test1",
			args: args{
				cliInput: "marwanhawari/ppath",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateCLIInput(tt.args.cliInput); (err != nil) != tt.wantErr {
				t.Errorf("ValidateCLIInput() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParseCLIInput(t *testing.T) {
	type args struct {
		cliInput string
	}
	tests := []struct {
		name    string
		args    args
		want    CLIInput
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				cliInput: "",
			},
			want:    CLIInput{},
			wantErr: true,
		},
		{
			name: "test2",
			args: args{
				cliInput: "marwanhawari/ppath",
			},
			want: CLIInput{
				IsGithubInput: true,
				Owner:         "marwanhawari",
				Repo:          "ppath",
			},
			wantErr: false,
		},
		{
			name: "test3",
			args: args{
				cliInput: "https://github.com/marwanhawari/ppath/releases/download/v0.0.3/ppath-v0.0.3-darwin-arm64.tar.gz",
			},
			want: CLIInput{
				IsGithubInput: false,
				Asset:         "ppath-v0.0.3-darwin-arm64.tar.gz",
				DownloadURL:   "https://github.com/marwanhawari/ppath/releases/download/v0.0.3/ppath-v0.0.3-darwin-arm64.tar.gz",
			},
			wantErr: false,
		},
		{
			name: "test4",
			args: args{
				cliInput: "marwan",
			},
			want:    CLIInput{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseCLIInput(tt.args.cliInput)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseCLIInput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseCLIInput() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseGithubInput(t *testing.T) {
	type args struct {
		cliInput string
	}
	tests := []struct {
		name    string
		args    args
		want    CLIInput
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				cliInput: "marwanhawari/ppath",
			},
			want: CLIInput{
				IsGithubInput: true,
				Owner:         "marwanhawari",
				Repo:          "ppath",
			},
			wantErr: false,
		},
		{
			name: "test2",
			args: args{
				cliInput: "marwanhawari/ppath@v0.0.3",
			},
			want: CLIInput{
				IsGithubInput: true,
				Owner:         "marwanhawari",
				Repo:          "ppath",
				Tag:           "v0.0.3",
			},
			wantErr: false,
		},
		{
			name: "test3",
			args: args{
				cliInput: "marwanhawari/ppath@v0.0.3::ppath-v0.0.3-linux-amd64.tar.gz",
			},
			want: CLIInput{
				IsGithubInput: true,
				Owner:         "marwanhawari",
				Repo:          "ppath",
				Tag:           "v0.0.3",
				Asset:         "ppath-v0.0.3-linux-amd64.tar.gz",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseGithubInput(tt.args.cliInput)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseGithubInput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseGithubInput() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseURLInput(t *testing.T) {
	type args struct {
		cliInput string
	}
	tests := []struct {
		name    string
		args    args
		want    CLIInput
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				cliInput: "https://github.com/marwanhawari/ppath/releases/download/v0.0.3/ppath-v0.0.3-darwin-arm64.tar.gz",
			},
			want: CLIInput{
				IsGithubInput: false,
				Asset:         "ppath-v0.0.3-darwin-arm64.tar.gz",
				DownloadURL:   "https://github.com/marwanhawari/ppath/releases/download/v0.0.3/ppath-v0.0.3-darwin-arm64.tar.gz",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseURLInput(tt.args.cliInput)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseURLInput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseURLInput() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContains(t *testing.T) {
	type args struct {
		slice  []string
		target string
	}
	tests := []struct {
		name  string
		args  args
		want  int
		want1 bool
	}{
		{
			name: "test1",
			args: args{
				slice:  []string{"a", "b", "c"},
				target: "a",
			},
			want:  0,
			want1: true,
		},
		{
			name: "test2",
			args: args{
				slice:  []string{"a", "b", "c"},
				target: "z",
			},
			want:  -1,
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := Contains(tt.args.slice, tt.args.target)
			if got != tt.want {
				t.Errorf("Contains() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Contains() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_getOS(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "test1",
			want: runtime.GOOS,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getOS(); got != tt.want {
				t.Errorf("getOS() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getArch(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "test1",
			want: runtime.GOARCH,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getArch(); got != tt.want {
				t.Errorf("getArch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_extractBinary(t *testing.T) {
	type args struct {
		downloadedFilePath string
		tmpExtractionPath  string
	}
	tests := []struct {
		name    string
		args    args
		url     string
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				downloadedFilePath: filepath.Join(t.TempDir(), "ppath-v0.0.3-darwin-arm64.tar.gz"),
				tmpExtractionPath:  filepath.Join(t.TempDir(), "tmp"),
			},
			url:     "https://github.com/marwanhawari/ppath/releases/download/v0.0.3/ppath-v0.0.3-darwin-arm64.tar.gz",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			DownloadFile(tt.args.downloadedFilePath, tt.url)

			if err := extractBinary(tt.args.downloadedFilePath, tt.args.tmpExtractionPath); (err != nil) != tt.wantErr {
				t.Errorf("extractBinary() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInstallBinary(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{
			name:    "test1",
			want:    "ppath",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			stewPath := filepath.Join(tempDir, ".stew")

			repo := "ppath"
			systemInfo := NewSystemInfo(stewPath)
			os.MkdirAll(systemInfo.StewBinPath, 0755)
			os.MkdirAll(systemInfo.StewPkgPath, 0755)
			os.MkdirAll(systemInfo.StewTmpPath, 0755)

			lockFile := LockFile{
				Os:   "darwin",
				Arch: "arm64",
				Packages: []PackageData{
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

			downloadedFilePath := filepath.Join(systemInfo.StewPkgPath, "ppath-v0.0.3-darwin-arm64.tar.gz")
			err := DownloadFile(downloadedFilePath, "https://github.com/marwanhawari/ppath/releases/download/v0.0.3/ppath-v0.0.3-darwin-arm64.tar.gz")

			if err != nil {
				t.Errorf("Could not download file to %v", downloadedFilePath)
			}

			got, err := InstallBinary(downloadedFilePath, repo, systemInfo, &lockFile, true)
			if (err != nil) != tt.wantErr {
				t.Errorf("InstallBinary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("InstallBinary() = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestInstallBinary_Fail(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{
			name:    "test1",
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			stewPath := filepath.Join(tempDir, ".stew")

			repo := "ppath"
			systemInfo := NewSystemInfo(stewPath)
			os.MkdirAll(systemInfo.StewBinPath, 0755)
			os.MkdirAll(systemInfo.StewPkgPath, 0755)
			os.MkdirAll(systemInfo.StewTmpPath, 0755)

			lockFile := LockFile{
				Os:   "darwin",
				Arch: "arm64",
				Packages: []PackageData{
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

			downloadedFilePath := filepath.Join(systemInfo.StewPkgPath, "ppath-v0.0.3-darwin-arm64.tar.gz")
			err := DownloadFile(downloadedFilePath, "https://github.com/marwanhawari/ppath/releases/download/v0.0.3/ppath-v0.0.3-darwin-arm64.tar.gz")

			if err != nil {
				t.Errorf("Could not download file to %v", downloadedFilePath)
			}

			got, err := InstallBinary(downloadedFilePath, repo, systemInfo, &lockFile, false)
			if (err != nil) != tt.wantErr {
				t.Errorf("InstallBinary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("InstallBinary() = %v, want %v", got, tt.want)
			}

		})
	}
}
