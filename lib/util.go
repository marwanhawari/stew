package stew

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/marwanhawari/stew/constants"
	"github.com/marwanhawari/stew/lib/config"
	"github.com/marwanhawari/stew/lib/ui/progress"
	"github.com/marwanhawari/stew/lib/ui/prompt"
	"github.com/marwanhawari/stew/lib/ui/terminal"
	"github.com/mholt/archiver"
	"github.com/pkg/errors"
)

func isArchiveFile(filePath string) bool {
	fileExtension := filepath.Ext(filePath)
	if fileExtension == ".br" || fileExtension == ".bz2" || fileExtension == ".zip" || fileExtension == ".gz" || fileExtension == ".lz4" || fileExtension == ".sz" || fileExtension == ".xz" || fileExtension == ".zst" || fileExtension == ".tar" || fileExtension == ".rar" {
		return true
	}
	return false
}

func isExecutableFile(filePath string) (bool, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return false, errors.WithStack(err)
	}

	filePerm := fileInfo.Mode()
	isExecutable := filePerm&0o111 != 0

	return isExecutable, nil
}

// DownloadFile will download a file from url to a given path.
func DownloadFile(rt config.Runtime, downloadPath string, url string) error {
	sp := progress.Spinner(rt)
	sp.Start()
	resp, err := http.Get(url)
	sp.Stop()

	if err != nil {
		return errors.WithStack(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return NonZeroStatusCodeDownloadError{StatusCode: resp.StatusCode}
	}

	outputFile, err := os.Create(downloadPath)
	if err != nil {
		return errors.WithStack(err)
	}
	defer outputFile.Close()

	bar := progress.Bar(
		rt, resp.ContentLength, "⬇️  Downloading asset:",
	)
	_, err = io.Copy(io.MultiWriter(outputFile, bar), resp.Body)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func copyFile(srcFile, destFile string) error {
	srcContents, err := os.Open(srcFile)
	if err != nil {
		return errors.WithStack(err)
	}
	defer srcContents.Close()

	destContents, err := os.Create(destFile)
	if err != nil {
		return errors.WithStack(err)
	}
	defer destContents.Close()

	_, err = io.Copy(destContents, srcContents)
	if err != nil {
		return errors.WithStack(err)
	}

	err = os.Chmod(destFile, 0o755)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func walkDir(rootDir string) ([]string, error) {
	var allFilePaths []string
	err := filepath.Walk(rootDir, func(filePath string, fileInfo os.FileInfo, err error) error {
		if !fileInfo.IsDir() {
			allFilePaths = append(allFilePaths, filePath)
		}
		return nil
	})
	return allFilePaths, err
}

func getBinary(io terminal.Terminal, filePaths []string, installation Installation) (string, string, error) {
	binaryFile := ""
	binaryName := ""
	var executableFiles []string
	for _, fullPath := range filePaths {
		fileNameBase := filepath.Base(fullPath)
		fileIsExecutable, err := isExecutableFile(fullPath)
		if err != nil {
			return "", "", err
		}
		if fileNameBase == installation.Repo && fileIsExecutable {
			binaryFile = fullPath
			binaryName = installation.Repo
			executableFiles = append(executableFiles, fullPath)
		} else if filepath.Ext(fullPath) == ".exe" {
			binaryFile = fullPath
			binaryName = filepath.Base(fullPath)
			executableFiles = append(executableFiles, fullPath)
		} else if fileIsExecutable {
			executableFiles = append(executableFiles, fullPath)
		}
	}

	if binaryFile == "" {
		if len(executableFiles) == 1 {
			binaryFile = executableFiles[0]
			binaryName = filepath.Base(binaryFile)
		} else if len(executableFiles) != 1 {
			if installation.BatchMode {
				return "", "", CouldntDetectBinaryError{filePaths}
			}
			var err error
			binaryFile, err = prompt.SelectWarn(
				io, "Could not automatically detect the binary. "+
					"Please select it manually:", filePaths)
			if err != nil {
				return "", "", err
			}
			binaryName = filepath.Base(binaryFile)
			binaryName, err = PromptRenameBinary(io, RenameBinaryArgs{
				Default: binaryName,
			})
			if err != nil {
				return "", "", err
			}
		}
	}

	if installation.BinaryName != "" {
		binaryName = installation.BinaryName
	}

	return binaryFile, binaryName, nil
}

// ValidateCLIInput makes sure the CLI input isn't empty.
func ValidateCLIInput(cliInput string) error {
	if cliInput == "" {
		return EmptyCLIInputError{}
	}

	return nil
}

// CLIInput contains information about the parsed CLI input.
type CLIInput struct {
	IsGithubInput bool
	Owner         string
	Repo          string
	Tag           string
	Asset         string
	DownloadURL   string
	BinaryName    string
}

// ParseCLIInput creates a new instance of the CLIInput struct.
func ParseCLIInput(cliInput string) (CLIInput, error) {
	err := ValidateCLIInput(cliInput)
	if err != nil {
		return CLIInput{}, err
	}

	var parsedInput CLIInput
	switch {
	case constants.RegexGithub.MatchString(cliInput):
		parsedInput, err = parseGithubInput(cliInput)
	case constants.RegexURL.MatchString(cliInput):
		parsedInput, err = parseURLInput(cliInput)
	default:
		return CLIInput{}, UnrecognizedInputError{
			Input: cliInput,
		}
	}
	if err != nil {
		return CLIInput{}, err
	}

	return parsedInput, nil
}

func parseGithubInput(cliInput string) (CLIInput, error) {
	parsedInput := CLIInput{
		IsGithubInput: true,
	}
	m := constants.RegexGithub.FindStringSubmatch(cliInput)
	if len(m) != 6 {
		return CLIInput{}, UnrecognizedInputError{
			Input: cliInput,
		}
	}
	parsedInput.Owner = m[1]
	parsedInput.Repo = m[2]
	parsedInput.Tag = m[3]
	parsedInput.Asset = m[4]
	parsedInput.BinaryName = m[5]

	return parsedInput, nil
}

func parseURLInput(cliInput string) (CLIInput, error) {
	if !constants.RegexURL.MatchString(cliInput) {
		return CLIInput{}, UnrecognizedInputError{
			Input: cliInput,
		}
	}
	return CLIInput{
		Asset:       filepath.Base(cliInput),
		DownloadURL: cliInput,
	}, nil
}

// Contains checks if a string slice contains a given target.
func Contains(slice []string, target string) (int, bool) {
	for index, element := range slice {
		if target == element {
			return index, true
		}
	}
	return -1, false
}

// Installation contains information about the installation process.
type Installation struct {
	DownloadedFilePath string
	Repo               string
	BinaryName         string
	BatchMode          bool
}

func extractBinary(io terminal.Terminal, installation Installation, tmpExtractionPath string) error {
	if isArchiveFile(installation.DownloadedFilePath) {
		err := archiver.Unarchive(installation.DownloadedFilePath, tmpExtractionPath)
		if err != nil {
			return errors.WithStack(err)
		}
		return nil
	}
	renameArgs := RenameBinaryArgs{
		Default:   installation.BinaryName,
		BatchMode: installation.BatchMode,
	}
	if renameArgs.Default == "" {
		renameArgs.Default = filepath.Base(installation.DownloadedFilePath)
	} else {
		// If the user specified a binary name, we assume they want to use it
		renameArgs.BatchMode = true
	}

	renamedBinaryName, err := PromptRenameBinary(io, renameArgs)
	if err != nil {
		return err
	}
	err = copyFile(installation.DownloadedFilePath, filepath.Join(tmpExtractionPath, renamedBinaryName))
	if err != nil {
		return err
	}
	return nil
}

// InstallBinary will extract the binary and copy it to the ~/.stew/bin path.
func InstallBinary(rt config.Runtime, installation Installation, lockFile *LockFile, overwriteFromUpgrade bool) (string, error) {
	tmpExtractionPath := rt.TmpPath
	assetDownloadPath := rt.PkgPath
	binaryInstallPath := rt.StewBinPath

	if err := os.RemoveAll(tmpExtractionPath); err != nil {
		return "", errors.WithStack(err)
	}
	if err := os.MkdirAll(tmpExtractionPath, 0o755); err != nil {
		return "", errors.WithStack(err)
	}

	err := extractBinary(rt, installation, tmpExtractionPath)
	if err != nil {
		return "", err
	}

	allFilePaths, err := walkDir(tmpExtractionPath)
	if err != nil {
		return "", err
	}

	binaryFile, binaryName, err := getBinary(rt, allFilePaths, installation)
	if err != nil {
		return "", err
	}

	// Check if the binary already exists
	for index, pkg := range lockFile.Packages {
		previousAssetPath := filepath.Join(assetDownloadPath, pkg.Asset)
		newAssetPath := installation.DownloadedFilePath
		var overwrite bool
		if pkg.Binary == binaryName {
			if !overwriteFromUpgrade && !installation.BatchMode {
				overwrite, err = prompt.ConfirmWarn(rt,
					fmt.Sprintf("The binary %v version: %v is already installed, "+
						"would you like to overwrite it?",
						constants.YellowColor(binaryName),
						constants.YellowColor(pkg.Tag)))
				if err != nil {
					_ = os.RemoveAll(newAssetPath)
					return "", err
				}
			} else {
				overwrite = true
			}

			if overwrite {
				err = os.RemoveAll(previousAssetPath)
				if err != nil {
					return "", errors.WithStack(err)
				}

				if !overwriteFromUpgrade {
					lockFile.Packages, err = RemovePackage(lockFile.Packages, index)
					if err != nil {
						return "", err
					}
				}
			} else {
				err = os.RemoveAll(newAssetPath)
				if err != nil {
					return "", errors.WithStack(err)
				}

				err = os.RemoveAll(tmpExtractionPath)
				if err != nil {
					return "", errors.WithStack(err)
				}

				return "", AbortBinaryOverwriteError{Binary: pkg.Binary}
			}
		}
	}

	err = copyFile(binaryFile, filepath.Join(binaryInstallPath, binaryName))
	if err != nil {
		return "", err
	}

	err = os.RemoveAll(tmpExtractionPath)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return binaryName, nil
}

// RenameBinaryArgs contains information about the renaming process.
type RenameBinaryArgs struct {
	Default   string
	BatchMode bool
}

// PromptRenameBinary takes in the original name of the binary and will return the new name of the binary.
func PromptRenameBinary(io terminal.Terminal, rename RenameBinaryArgs) (string, error) {
	if rename.BatchMode {
		return rename.Default, nil
	}
	renamedBinaryName, err := prompt.InputWarn(io, "Rename the binary?", rename.Default)
	if err != nil {
		return "", err
	}
	return renamedBinaryName, nil
}
