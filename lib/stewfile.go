package stew

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/marwanhawari/stew/constants"
)

// LockFile contains all the data for the lockfile
type LockFile struct {
	Os       string        `json:"os"`
	Arch     string        `json:"arch"`
	Packages []PackageData `json:"packages"`
}

// PackageData contains the information for an installed binary
type PackageData struct {
	Source string `json:"source"`
	Owner  string `json:"owner"`
	Repo   string `json:"repo"`
	Tag    string `json:"tag"`
	Asset  string `json:"asset"`
	Binary string `json:"binary"`
	URL    string `json:"url"`
}

func readLockFileJSON(lockFilePath string) (LockFile, error) {

	lockFileBytes, err := ioutil.ReadFile(lockFilePath)
	if err != nil {
		return LockFile{}, err
	}

	var lockFile LockFile
	err = json.Unmarshal(lockFileBytes, &lockFile)
	if err != nil {
		return LockFile{}, err
	}

	return lockFile, nil
}

// WriteLockFileJSON will write the lockfile JSON file
func WriteLockFileJSON(lockFileJSON LockFile, outputPath string) error {

	lockFileBytes, err := json.MarshalIndent(lockFileJSON, "", "\t")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(outputPath, lockFileBytes, 0644)
	if err != nil {
		return err
	}

	fmt.Printf("📄 Updated %v\n", constants.GreenColor(outputPath))

	return nil
}

// RemovePackage will remove a package from a LockFile.Packages slice
func RemovePackage(pkgs []PackageData, index int) ([]PackageData, error) {
	if len(pkgs) == 0 {
		return []PackageData{}, NoPackagesInLockfileError{}
	}

	if index < 0 || index >= len(pkgs) {
		return []PackageData{}, IndexOutOfBoundsInLockfileError{}
	}

	return append(pkgs[:index], pkgs[index+1:]...), nil
}

// ReadStewfileContents will read the contents of the Stewfile
func ReadStewfileContents(stewfilePath string) ([]string, error) {
	file, err := os.Open(stewfilePath)
	if err != nil {
		return []string{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var packages []string
	for scanner.Scan() {
		packages = append(packages, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return []string{}, err
	}

	return packages, nil
}

// NewLockFile creates a new instance of the LockFile struct
func NewLockFile(stewLockFilePath, userOS, userArch string) (LockFile, error) {
	var lockFile LockFile
	lockFileExists, err := PathExists(stewLockFilePath)
	if err != nil {
		return LockFile{}, err
	}
	if !lockFileExists {
		lockFile = LockFile{Os: userOS, Arch: userArch, Packages: []PackageData{}}
	} else {
		lockFile, err = readLockFileJSON(stewLockFilePath)
		if err != nil {
			return LockFile{}, err
		}
	}
	return lockFile, nil
}

// SystemInfo contains system specific info like OS, arch, and ~/.stew paths
type SystemInfo struct {
	Os               string
	Arch             string
	StewPath         string
	StewBinPath      string
	StewPkgPath      string
	StewLockFilePath string
	StewTmpPath      string
}

// NewSystemInfo creates a new instance of the SystemInfo struct
func NewSystemInfo() (SystemInfo, error) {
	stewPath, err := getStewPath()
	if err != nil {
		return SystemInfo{}, err
	}

	var systemInfo SystemInfo
	systemInfo.StewPath = stewPath
	systemInfo.StewBinPath = path.Join(stewPath, "bin")
	systemInfo.StewPkgPath = path.Join(stewPath, "pkg")
	systemInfo.StewLockFilePath = path.Join(stewPath, "Stewfile.lock.json")
	systemInfo.StewTmpPath = path.Join(stewPath, "tmp")
	systemInfo.Os = getOS()
	systemInfo.Arch = getArch()

	return systemInfo, nil
}

// DeleteAssetAndBinary will delete the asset from the ~/.stew/pkg path and delete the binary from the ~/.stew/bin path
func DeleteAssetAndBinary(stewPkgPath, stewBinPath, asset, binary string) error {
	assetPath := path.Join(stewPkgPath, asset)
	binPath := path.Join(stewBinPath, binary)
	err := os.RemoveAll(assetPath)
	if err != nil {
		return err
	}
	err = os.RemoveAll(binPath)
	if err != nil {
		return err
	}
	return nil
}
