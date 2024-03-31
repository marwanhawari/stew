package stew

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/marwanhawari/stew/constants"
	"github.com/mholt/archiver"
	progressbar "github.com/schollz/progressbar/v3"
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
		return false, err
	}

	filePerm := fileInfo.Mode()
	isExecutable := filePerm&0111 != 0

	return isExecutable, nil
}

// CatchAndExit will catch errors and immediately exit
func CatchAndExit(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// PathExists checks if a given path exists
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)

	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// DownloadFile will download a file from url to a given path
func DownloadFile(downloadPath string, url string) error {
	sp := constants.LoadingSpinner
	sp.Start()
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	if strings.Contains(url, "api.github.com") {
		req.Header.Add("Accept", "application/octet-stream")
		githubToken := os.Getenv("GITHUB_TOKEN")
		if githubToken != "" {
			req.Header.Add("Authorization", fmt.Sprintf("token %v", githubToken))
		}
	}

	resp, err := client.Do(req)
	sp.Stop()

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return NonZeroStatusCodeDownloadError{StatusCode: resp.StatusCode}
	}

	outputFile, err := os.Create(downloadPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		"⬇️  Downloading asset:",
	)
	_, err = io.Copy(io.MultiWriter(outputFile, bar), resp.Body)
	if err != nil {
		return err
	}

	_, err = io.Copy(outputFile, resp.Body)
	if err != nil {
		return err
	}

	return nil

}

func copyFile(srcFile, destFile string) error {
	srcContents, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer srcContents.Close()

	destContents, err := os.Create(destFile)
	if err != nil {
		return err
	}
	defer destContents.Close()

	_, err = io.Copy(destContents, srcContents)
	if err != nil {
		return err
	}

	err = os.Chmod(destFile, 0755)
	if err != nil {
		return err
	}

	return nil
}

func walkDir(rootDir string) ([]string, error) {
	allFilePaths := []string{}
	err := filepath.Walk(rootDir, func(filePath string, fileInfo os.FileInfo, err error) error {
		if !fileInfo.IsDir() {
			allFilePaths = append(allFilePaths, filePath)
		}
		return nil
	})
	return allFilePaths, err
}

func getBinary(filePaths []string, repo string) (string, string, error) {
	binaryFile := ""
	binaryName := ""
	var err error
	executableFiles := []string{}
	for _, fullPath := range filePaths {
		fileNameBase := filepath.Base(fullPath)
		fileIsExecutable, err := isExecutableFile(fullPath)
		if err != nil {
			return "", "", err
		}
		if fileNameBase == repo && fileIsExecutable {
			binaryFile = fullPath
			binaryName = repo
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
			binaryFile, err = WarningPromptSelect("Could not automatically detect the binary. Please select it manually:", filePaths)
			if err != nil {
				return "", "", err
			}
			binaryName = filepath.Base(binaryFile)
			binaryName, err = PromptRenameBinary(binaryName)
			if err != nil {
				return "", "", nil
			}
		}
	}

	return binaryFile, binaryName, nil
}

// ValidateCLIInput makes sure the CLI input isn't empty
func ValidateCLIInput(cliInput string) error {
	if cliInput == "" {
		return EmptyCLIInputError{}
	}

	return nil
}

// CLIInput contains information about the parsed CLI input
type CLIInput struct {
	IsGithubInput bool
	Owner         string
	Repo          string
	Tag           string
	Asset         string
	DownloadURL   string
}

// ParseCLIInput creates a new instance of the CLIInput struct
func ParseCLIInput(cliInput string) (CLIInput, error) {
	err := ValidateCLIInput(cliInput)
	if err != nil {
		return CLIInput{}, err
	}

	reGithub, err := regexp.Compile(constants.RegexGithub)
	if err != nil {
		return CLIInput{}, err
	}
	reURL, err := regexp.Compile(constants.RegexURL)
	if err != nil {
		return CLIInput{}, err
	}
	var parsedInput CLIInput
	if reGithub.MatchString(cliInput) {
		parsedInput, err = parseGithubInput(cliInput)
	} else if reURL.MatchString(cliInput) {
		parsedInput, err = parseURLInput(cliInput)
	} else {
		return CLIInput{}, UnrecognizedInputError{}
	}
	if err != nil {
		return CLIInput{}, err
	}

	return parsedInput, nil

}

func parseGithubInput(cliInput string) (CLIInput, error) {
	parsedInput := CLIInput{}
	parsedInput.IsGithubInput = true
	trimmedString := strings.Trim(strings.Trim(strings.Trim(strings.TrimSpace(cliInput), "/"), "@"), "::")
	splitInput := strings.SplitN(trimmedString, "@", 2)

	ownerAndRepo := splitInput[0]
	splitOwnerAndRepo := strings.SplitN(ownerAndRepo, "/", 2)
	parsedInput.Owner = splitOwnerAndRepo[0]
	parsedInput.Repo = splitOwnerAndRepo[1]

	if len(splitInput) == 2 {
		tagAndAsset := splitInput[1]
		splitTagAndAsset := strings.SplitN(tagAndAsset, "::", 2)
		parsedInput.Tag = splitTagAndAsset[0]
		if len(splitTagAndAsset) == 2 {
			parsedInput.Asset = splitTagAndAsset[1]
		}
	}

	return parsedInput, nil

}

func parseURLInput(cliInput string) (CLIInput, error) {
	return CLIInput{IsGithubInput: false, Asset: filepath.Base(cliInput), DownloadURL: cliInput}, nil
}

// Contains checks if a string slice contains a given target
func Contains(slice []string, target string) (int, bool) {
	for index, element := range slice {
		if target == element {
			return index, true
		}
	}
	return -1, false
}

func extractBinary(downloadedFilePath, tmpExtractionPath string) error {
	isArchive := isArchiveFile(downloadedFilePath)
	if isArchive {
		err := archiver.Unarchive(downloadedFilePath, tmpExtractionPath)
		if err != nil {
			return err
		}

	} else {
		originalBinaryName := filepath.Base(downloadedFilePath)

		renamedBinaryName, err := PromptRenameBinary(originalBinaryName)
		if err != nil {
			return err
		}
		err = copyFile(downloadedFilePath, filepath.Join(tmpExtractionPath, renamedBinaryName))
		if err != nil {
			return err
		}
	}
	return nil
}

// InstallBinary will extract the binary and copy it to the ~/.stew/bin path
func InstallBinary(downloadedFilePath string, repo string, systemInfo SystemInfo, lockFile *LockFile, overwriteFromUpgrade bool) (string, error) {

	tmpExtractionPath := systemInfo.StewTmpPath
	assetDownloadPath := systemInfo.StewPkgPath
	binaryInstallPath := systemInfo.StewBinPath

	err := extractBinary(downloadedFilePath, tmpExtractionPath)
	if err != nil {
		return "", err
	}

	allFilePaths, err := walkDir(tmpExtractionPath)
	if err != nil {
		return "", err
	}

	binaryFile, binaryName, err := getBinary(allFilePaths, repo)
	if err != nil {
		return "", err
	}

	// Check if the binary already exists
	for index, pkg := range lockFile.Packages {
		previousAssetPath := filepath.Join(assetDownloadPath, pkg.Asset)
		newAssetPath := downloadedFilePath
		var overwrite bool
		if pkg.Binary == binaryName {
			if !overwriteFromUpgrade {
				overwrite, err = WarningPromptConfirm(fmt.Sprintf("The binary %v version: %v is already installed, would you like to overwrite it?", constants.YellowColor(binaryName), constants.YellowColor(pkg.Tag)))
				if err != nil {
					os.RemoveAll(newAssetPath)
					return "", err
				}
			} else {
				overwrite = true
			}

			if overwrite {
				err := os.RemoveAll(previousAssetPath)
				if err != nil {
					return "", err
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
					return "", err
				}

				err = os.RemoveAll(tmpExtractionPath)
				if err != nil {
					return "", err
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
		return "", err
	}

	return binaryName, nil
}

// PromptRenameBinary takes in the original name of the binary and will return the new name of the binary.
func PromptRenameBinary(originalBinaryName string) (string, error) {
	renamedBinaryName, err := warningPromptInput("Rename the binary?", originalBinaryName)
	if err != nil {
		return "", err
	}
	return renamedBinaryName, nil
}

// ResolvePath will resolve the full path for an input path
func ResolvePath(filePath string) (string, error) {
	var resolvedPath string
	var err error
	resolvedPath = filePath

	resolvedPath = strings.ReplaceAll(resolvedPath, "\"", "")

	if strings.HasPrefix(resolvedPath, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}

		resolvedPath = filepath.Join(homeDir, strings.TrimLeft(resolvedPath, "~"))
	}

	resolvedPath = os.ExpandEnv(resolvedPath)

	if !filepath.IsAbs(resolvedPath) {
		resolvedPath, err = filepath.Abs(resolvedPath)
		if err != nil {
			return "", err
		}
	}

	resolvedPath = strings.TrimRight(resolvedPath, "/")

	return resolvedPath, nil
}
