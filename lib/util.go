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
	_, err := archiver.ByExtension(filePath)
	return err == nil
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

	return os.Chmod(destFile, 0755)
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

type ExecutableFileInfo struct {
	fileName string
	filePath string
	fileHash string
}

func getBinary(filePaths []string, desiredBinaryRename, expectedBinaryHash string) (string, string, string, error) {
	executableFiles := []ExecutableFileInfo{}
	for _, fullPath := range filePaths {
		fileNameBase := filepath.Base(fullPath)
		fileIsExecutable, err := isExecutableFile(fullPath)
		if err != nil {
			return "", "", "", err
		}
		fileHash, err := calculateFileHash(fullPath)
		if err != nil {
			return "", "", "", err
		}
		if desiredBinaryRename != "" && expectedBinaryHash != "" && expectedBinaryHash == fileHash {
			return fullPath, desiredBinaryRename, expectedBinaryHash, nil
		}
		if !fileIsExecutable {
			continue
		}
		executableFiles = append(executableFiles, ExecutableFileInfo{fileName: fileNameBase, filePath: fullPath, fileHash: fileHash})
	}

	if len(executableFiles) != 1 {
		binaryFilePath, err := WarningPromptSelect("Could not automatically detect the binary. Please select it manually:", filePaths)
		if err != nil {
			return "", "", "", err
		}
		binaryName, err := PromptRenameBinary(filepath.Base(binaryFilePath))
		if err != nil {
			return "", "", "", err
		}
		binaryHash, err := calculateFileHash(binaryFilePath)
		if err != nil {
			return "", "", "", err
		}
		return binaryFilePath, binaryName, binaryHash, nil
	}

	if desiredBinaryRename != "" {
		return executableFiles[0].filePath, desiredBinaryRename, executableFiles[0].fileHash, nil
	}
	return executableFiles[0].filePath, executableFiles[0].fileName, executableFiles[0].fileHash, nil
}

// ValidateCLIInput makes sure the CLI input isn't empty
func ValidateCLIInput(cliInput string) error {
	if cliInput == "" {
		return EmptyCLIInputError{}
	}
	return nil
}

// ParseCLIInput creates a new instance of the PackageData struct
func ParseCLIInput(cliInput string) (PackageData, error) {
	err := ValidateCLIInput(cliInput)
	if err != nil {
		return PackageData{}, err
	}

	reGithub, err := regexp.Compile(constants.RegexGithub)
	if err != nil {
		return PackageData{}, err
	}
	reURL, err := regexp.Compile(constants.RegexURL)
	if err != nil {
		return PackageData{}, err
	}
	var parsedInput PackageData
	if reGithub.MatchString(cliInput) {
		parsedInput, err = parseGithubInput(cliInput)
	} else if reURL.MatchString(cliInput) {
		parsedInput, err = parseURLInput(cliInput)
	} else {
		return PackageData{}, UnrecognizedInputError{}
	}
	if err != nil {
		return PackageData{}, err
	}

	return parsedInput, nil

}

func parseGithubInput(cliInput string) (PackageData, error) {
	parsedInput := PackageData{}
	parsedInput.Source = "github"
	trimmedString := strings.Trim(strings.Trim(strings.TrimSpace(cliInput), "/"), "@")
	splitInput := strings.SplitN(trimmedString, "@", 2)

	ownerAndRepo := splitInput[0]
	splitOwnerAndRepo := strings.SplitN(ownerAndRepo, "/", 2)
	parsedInput.Owner = splitOwnerAndRepo[0]
	parsedInput.Repo = splitOwnerAndRepo[1]

	if len(splitInput) == 2 {
		parsedInput.Tag = splitInput[1]
	}

	return parsedInput, nil

}

func parseURLInput(cliInput string) (PackageData, error) {
	return PackageData{Source: "other", Asset: filepath.Base(cliInput), URL: cliInput}, nil
}

// Contains checks if a string slice contains a given target
func Contains[T comparable](slice []T, target T) (int, bool) {
	for index, element := range slice {
		if target == element {
			return index, true
		}
	}
	return -1, false
}

func FindBinaryInLockFile(lockFile LockFile, binaryName string) (int, bool) {
	for index, pkg := range lockFile.Packages {
		if pkg.Binary == binaryName {
			return index, true
		}
	}
	return -1, false
}

func extractBinary(downloadedFilePath, tmpExtractionPath, desiredBinaryRename string) error {
	isArchive := isArchiveFile(downloadedFilePath)
	if isArchive {
		err := archiver.Unarchive(downloadedFilePath, tmpExtractionPath)
		if err != nil {
			return err
		}
		return nil
	}
	originalBinaryName := filepath.Base(downloadedFilePath)
	if desiredBinaryRename != "" {
		return copyFile(downloadedFilePath, filepath.Join(tmpExtractionPath, desiredBinaryRename))
	}
	renamedBinaryName, err := PromptRenameBinary(originalBinaryName)
	if err != nil {
		return err
	}
	return copyFile(downloadedFilePath, filepath.Join(tmpExtractionPath, renamedBinaryName))
}

// InstallBinary will extract the binary and copy it to the ~/.stew/bin path
func InstallBinary(downloadedFilePath string, repo string, systemInfo SystemInfo, lockFile *LockFile, overwriteFromUpgrade bool, desiredBinaryRename, expectedBinaryHash string) (string, string, error) {
	tmpExtractionPath, stewPkgPath, binaryInstallPath := systemInfo.StewTmpPath, systemInfo.StewPkgPath, systemInfo.StewBinPath
	if err := extractBinary(downloadedFilePath, tmpExtractionPath, desiredBinaryRename); err != nil {
		return "", "", err
	}

	allFilePaths, err := walkDir(tmpExtractionPath)
	if err != nil {
		return "", "", err
	}

	binaryFileInTmpExtractionPath, binaryName, binaryHash, err := getBinary(allFilePaths, desiredBinaryRename, expectedBinaryHash)
	if err != nil {
		return "", "", err
	}

	if err = handleExistingBinary(lockFile, binaryName, downloadedFilePath, stewPkgPath, overwriteFromUpgrade); err != nil {
		return "", "", err
	}

	err = copyFile(binaryFileInTmpExtractionPath, filepath.Join(binaryInstallPath, binaryName))
	if err != nil {
		return "", "", err
	}

	err = os.RemoveAll(tmpExtractionPath)
	if err != nil {
		return "", "", err
	}

	return binaryName, binaryHash, nil
}

func handleExistingBinary(lockFile *LockFile, binaryName, newlyDownloadedAssetPath, stewPkgPath string, overwriteFromUpgrade bool) error {
	indexInLockFile, binaryFoundInLockFile := FindBinaryInLockFile(*lockFile, binaryName)
	if !binaryFoundInLockFile {
		return nil
	}
	pkg := lockFile.Packages[indexInLockFile]
	if !overwriteFromUpgrade {
		userChoosingToOverwrite, err := WarningPromptConfirm(fmt.Sprintf("The binary %v version: %v is already installed, would you like to overwrite it?", constants.YellowColor(binaryName), constants.YellowColor(pkg.Tag)))
		if err != nil {
			if err := os.RemoveAll(newlyDownloadedAssetPath); err != nil {
				return err
			}
			return err
		}
		if !userChoosingToOverwrite {
			if err := os.RemoveAll(newlyDownloadedAssetPath); err != nil {
				return err
			}
			return AbortBinaryOverwriteError{Binary: binaryName}
		}
	}
	return overwriteBinary(lockFile, indexInLockFile, newlyDownloadedAssetPath, stewPkgPath, overwriteFromUpgrade)
}

func overwriteBinary(lockFile *LockFile, indexInLockFile int, newlyDownloadedAssetPath, stewPkgPath string, overwriteFromUpgrade bool) error {
	pkg := lockFile.Packages[indexInLockFile]
	previousAssetPath := filepath.Join(stewPkgPath, pkg.Asset)
	if previousAssetPath != newlyDownloadedAssetPath {
		if err := os.RemoveAll(previousAssetPath); err != nil {
			return err
		}
	}
	// If not overwriting as part of an upgrade, remove the package entry from the lock file
	// This is because the upgrade command will update the package entry in place
	// but the install command will add a new package entry
	if !overwriteFromUpgrade {
		var err error
		lockFile.Packages, err = RemovePackage(lockFile.Packages, indexInLockFile)
		if err != nil {
			return err
		}
	}
	return nil
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
