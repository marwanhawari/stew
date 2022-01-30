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

type LockFile struct {
	Os       string        `json:"os"`
	Arch     string        `json:"arch"`
	Packages []PackageData `json:"packages"`
}
type PackageData struct {
	Source string `json:"source"`
	Owner  string `json:"owner"`
	Repo   string `json:"repo"`
	Tag    string `json:"tag"`
	Asset  string `json:"asset"`
	Binary string `json:"binary"`
	URL    string `json:"url"`
}

func ReadLockFileJSON(lockFilePath string) (LockFile, error) {

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

func RemovePackage(pkgs []PackageData, index int) ([]PackageData, error) {
	if len(pkgs) == 0 {
		return []PackageData{}, NoPackagesInLockfileError{}
	}

	if index < 0 || index >= len(pkgs) {
		return []PackageData{}, IndexOutOfBoundsInLockfileError{}
	}

	return append(pkgs[:index], pkgs[index+1:]...), nil
}

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

func NewLockFile(stewLockFilePath, userOS, userArch string) (LockFile, error) {
	var lockFile LockFile
	lockFileExists, err := PathExists(stewLockFilePath)
	if err != nil {
		return LockFile{}, err
	}
	if !lockFileExists {
		lockFile = LockFile{Os: userOS, Arch: userArch, Packages: []PackageData{}}
	} else {
		lockFile, err = ReadLockFileJSON(stewLockFilePath)
		if err != nil {
			return LockFile{}, err
		}
	}
	return lockFile, nil
}

type SystemInfo struct {
	Os               string
	Arch             string
	StewPath         string
	StewBinPath      string
	StewPkgPath      string
	StewLockFilePath string
	StewTmpPath      string
}

func NewSystemInfo() (SystemInfo, error) {
	stewPath, err := GetStewPath()
	if err != nil {
		return SystemInfo{}, err
	}

	var systemInfo SystemInfo
	systemInfo.StewPath = stewPath
	systemInfo.StewBinPath = path.Join(stewPath, "bin")
	systemInfo.StewPkgPath = path.Join(stewPath, "pkg")
	systemInfo.StewLockFilePath = path.Join(stewPath, "Stewfile.lock.json")
	systemInfo.StewTmpPath = path.Join(stewPath, "tmp")
	systemInfo.Os = GetOS()
	systemInfo.Arch = GetArch()

	return systemInfo, nil
}

func DeleteAssetAndBinary(stewPkgPath, stewBinPath, asset, binary string) error {
	assetPath := path.Join(stewPkgPath, asset)
	pkgPath := path.Join(stewBinPath, binary)
	err := os.RemoveAll(assetPath)
	if err != nil {
		return err
	}
	err = os.RemoveAll(pkgPath)
	if err != nil {
		return err
	}
	return nil
}
