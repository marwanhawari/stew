package stew

import (
	"errors"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/marwanhawari/stew/lib/pathutil"
	"github.com/marwanhawari/stew/lib/testsupport"
	"github.com/marwanhawari/stew/lib/ui/prompt"
	"github.com/marwanhawari/stew/lib/ui/terminal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	RegularFilePerm    = 0o644
	ExecutableFilePerm = 0o755
)

func Test_isArchiveFile(t *testing.T) {
	type TestCase struct {
		name     string
		filePath string
		want     bool
	}
	cases := []TestCase{{
		name:     "test1",
		filePath: "notArchive",
	}, {
		name:     "test1",
		filePath: "Archive.tar.gz",
		want:     true,
	}}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if got := isArchiveFile(tt.filePath); got != tt.want {
				t.Errorf("isArchiveFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isExecutableFile(t *testing.T) {
	type TestCase struct {
		name     string
		filePath string
		want     bool
		wantErr  bool
	}
	cases := []TestCase{{
		name:     "test1",
		filePath: "testExecutableFile",
		want:     true,
	}, {
		name:     "test1",
		filePath: "notExecutableFile",
	}}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			testExecutableFilePath := filepath.Join(tempDir, tt.filePath)
			if tt.want {
				require.NoError(t, os.WriteFile(testExecutableFilePath, []byte("An executable file"), ExecutableFilePerm))
			} else {
				require.NoError(t, os.WriteFile(testExecutableFilePath, []byte("Not an executable file"), RegularFilePerm))
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

func TestDownloadFile(t *testing.T) {
	type TestCase struct {
		name    string
		url     string
		wantErr error
	}
	cases := []TestCase{{
		name: "test1",
		url:  "https://github.com/marwanhawari/ppath/releases/download/v0.0.3/ppath-v0.0.3-darwin-arm64.tar.gz",
	}, {
		name: "test2",
		url:  "",
		wantErr: &url.Error{
			Op:  "Get",
			URL: "",
			Err: errors.New(`unsupported protocol scheme ""`),
		},
	}}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			rt := testsupport.NewDefaultRuntime(t)
			testDownloadPath := filepath.Join(tempDir, filepath.Base(tt.url))
			if err := DownloadFile(rt, testDownloadPath, tt.url); tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Fatalf("DownloadFile() error\n"+
					" got = %+v\n"+
					"want = %+v", err, tt.wantErr)
			}
			if fileExists, _ := pathutil.Exists(testDownloadPath); !fileExists {
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
	}{{
		name: "test1",
		args: args{
			srcFile:  "sourceFile.txt",
			destFile: "destFile.txt",
		},
		wantErr: false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			srcFilePath := filepath.Join(tempDir, tt.args.srcFile)
			destFilePath := filepath.Join(tempDir, tt.args.destFile)

			require.NoError(t, os.WriteFile(srcFilePath, []byte("A test file"), RegularFilePerm))

			srcExists, _ := pathutil.Exists(srcFilePath)
			destExists, _ := pathutil.Exists(destFilePath)

			if !srcExists {
				t.Errorf("Source file %v not found", srcFilePath)
			}

			if destExists {
				t.Errorf("Dest file %v already exists", destFilePath)
			}

			if err := copyFile(srcFilePath, destFilePath); (err != nil) != tt.wantErr {
				t.Errorf("copyFile() error = %v, wantErr %v", err, tt.wantErr)
			}

			srcExists, _ = pathutil.Exists(srcFilePath)
			destExists, _ = pathutil.Exists(destFilePath)

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
	}{{
		name:    "test1",
		wantErr: false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()

			require.NoError(t, os.WriteFile(filepath.Join(tempDir, "testFile.txt"),
				[]byte("A test file"), RegularFilePerm))
			require.NoError(t, os.MkdirAll(filepath.Join(tempDir, "bin"), ExecutableFilePerm))
			require.NoError(t, os.WriteFile(filepath.Join(tempDir, "bin",
				"binDirTestFile.txt"), []byte("Another test file"), RegularFilePerm))

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
	tests := []struct {
		name       string
		repo       string
		binaryName string
		wantErr    bool
	}{{
		name:       "test1",
		repo:       "testBinary",
		binaryName: "testBinary",
		wantErr:    false,
	}, {
		name:       "test2",
		repo:       "someRepo",
		binaryName: "testBinary",
		wantErr:    false,
	}, {
		name:       "test3",
		repo:       "someRepo",
		binaryName: "testBinary.exe",
		wantErr:    false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			prt := terminal.TestTerminal{TestingT: t}

			testBinaryFilePath := filepath.Join(tempDir, tt.binaryName)
			require.NoError(t, os.WriteFile(testBinaryFilePath, []byte("An executable file"), ExecutableFilePerm))
			testNonBinaryFilePath := filepath.Join(tempDir, "testNonBinary")
			require.NoError(t, os.WriteFile(testNonBinaryFilePath, []byte("Not an executable file"), RegularFilePerm))

			testFilePaths := []string{testBinaryFilePath, testNonBinaryFilePath}

			wantBinaryFile := filepath.Join(tempDir, tt.binaryName)
			wantBinaryName := filepath.Base(wantBinaryFile)

			install := Installation{
				Repo: tt.repo,
			}
			got, got1, err := getBinary(prt, testFilePaths, install)
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
	tests := []struct {
		name       string
		repo       string
		binaryName string
		wantErr    bool
	}{{
		name:       "test1",
		repo:       "testBinary",
		binaryName: "testBinary",
		wantErr:    true,
	}}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			io := terminal.TestTerminal{TestingT: t}

			testNonBinaryFilePath := filepath.Join(tempDir, "testNonBinary")
			require.NoError(t, os.WriteFile(testNonBinaryFilePath,
				[]byte("Not an executable file"), RegularFilePerm))

			testFilePaths := []string{testNonBinaryFilePath}

			wantBinaryFile := ""
			wantBinaryName := ""

			install := Installation{
				Repo:      tt.repo,
				BatchMode: true,
			}
			got, got1, err := getBinary(io, testFilePaths, install)
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

func TestParseGithubInput(t *testing.T) {
	t.Parallel()
	tests := []struct {
		args    string
		want    CLIInput
		wantErr error
	}{{
		args: "marwanhawari/ppath",
		want: CLIInput{
			IsGithubInput: true,
			Owner:         "marwanhawari",
			Repo:          "ppath",
		},
	}, {
		args: "marwanhawari/ppath@v0.0.3",
		want: CLIInput{
			IsGithubInput: true,
			Owner:         "marwanhawari",
			Repo:          "ppath",
			Tag:           "v0.0.3",
		},
	}, {
		args: "marwanhawari/ppath@v0.0.3::ppath-v0.0.3-linux-amd64.tar.gz",
		want: CLIInput{
			IsGithubInput: true,
			Owner:         "marwanhawari",
			Repo:          "ppath",
			Tag:           "v0.0.3",
			Asset:         "ppath-v0.0.3-linux-amd64.tar.gz",
		},
	}, {
		args: "knative-sandbox/kn-plugin-event!!kn-event",
		want: CLIInput{
			IsGithubInput: true,
			Owner:         "knative-sandbox",
			Repo:          "kn-plugin-event",
			BinaryName:    "kn-event",
		},
	}, {
		args: "knative-sandbox/kn-plugin-event@v1.9.1!!kn-event",
		want: CLIInput{
			IsGithubInput: true,
			Owner:         "knative-sandbox",
			Repo:          "kn-plugin-event",
			Tag:           "v1.9.1",
			BinaryName:    "kn-event",
		},
	}, {
		args: "knative-sandbox/kn-plugin-event@v1.9.1::kn-event-linux-amd64!!kn-event",
		want: CLIInput{
			IsGithubInput: true,
			Owner:         "knative-sandbox",
			Repo:          "kn-plugin-event",
			Tag:           "v1.9.1",
			Asset:         "kn-event-linux-amd64",
			BinaryName:    "kn-event",
		},
	}, {
		args:    "foo/bar!!foobar@1.2",
		wantErr: UnrecognizedInputError{"foo/bar!!foobar@1.2"},
	}}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.args, func(t *testing.T) {
			t.Parallel()
			got, err := ParseCLIInput(tt.args)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("error mismatch\n"+
					" got = %+v,\n"+
					"want = %+v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("results are not equal\n"+
					" got = %+v,\n"+
					"want = %+v", got, tt.want)
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
	}{{
		name: "test1",
		args: args{
			downloadedFilePath: filepath.Join(t.TempDir(), "ppath-v0.0.3-darwin-arm64.tar.gz"),
			tmpExtractionPath:  filepath.Join(t.TempDir(), "tmp"),
		},
		url:     "https://github.com/marwanhawari/ppath/releases/download/v0.0.3/ppath-v0.0.3-darwin-arm64.tar.gz",
		wantErr: false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := testsupport.NewDefaultRuntime(t)
			require.NoError(t, DownloadFile(rt, tt.args.downloadedFilePath, tt.url))

			install := Installation{
				DownloadedFilePath: tt.args.downloadedFilePath,
			}
			if err := extractBinary(rt, install, tt.args.tmpExtractionPath); (err != nil) != tt.wantErr {
				t.Errorf("extractBinary() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInstallBinary(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		pd      PackageData
		batch   bool
		wantErr error
	}{{
		name:  "test1",
		want:  "hello",
		batch: true,
		pd: PackageData{
			Source: "github",
			Owner:  "wavesoftware",
			Repo:   "asm-sandbox",
			Tag:    "v0.1.0",
			Asset:  "hello-0.1.0-linux-amd64.tar.xz",
			Binary: "hello",
			URL:    "https://github.com/wavesoftware/asm-sandbox/releases/download/v0.1.0/hello-0.1.0-linux-amd64.tar.xz",
		},
	}, {
		name: "test2",
		pd: PackageData{
			Source: "github",
			Owner:  "wavesoftware",
			Repo:   "asm-sandbox",
			Tag:    "v0.1.0",
			Asset:  "not-exists.txt",
			URL:    "https://github.com/marwanhawari/ppath/releases/download/v0.0.3/checksums.txt",
		},
		wantErr: prompt.ExitUserSelectionError{},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lockFile := LockFile{
				Os:       "darwin",
				Arch:     "arm64",
				Packages: []PackageData{tt.pd},
			}
			rt := testsupport.NewRuntime(t, lockFile.Os, lockFile.Arch)

			downloadedFilePath := filepath.Join(rt.PkgPath, tt.pd.Asset)
			err := DownloadFile(rt, downloadedFilePath, tt.pd.URL)

			if err != nil {
				t.Fatalf("Could not download file to %v", downloadedFilePath)
			}

			install := Installation{
				DownloadedFilePath: downloadedFilePath,
				Repo:               tt.pd.Repo,
				BatchMode:          tt.batch,
				BinaryName:         tt.pd.Binary,
			}
			got, err := InstallBinary(rt, install, &lockFile, false)
			if tt.wantErr != nil && reflect.TypeOf(tt.wantErr) != reflect.TypeOf(err) {
				t.Fatalf("error\n"+
					" got = %+v\n"+
					"want = %+v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
