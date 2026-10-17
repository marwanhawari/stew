package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/marwanhawari/stew/constants"
	stew "github.com/marwanhawari/stew/lib"
)

// Install is executed when you run `stew install`
func Install(cliInput string) {
	userOS, userArch, _, systemInfo, err := stew.Initialize()
	stew.CatchAndExit(err)

	if isLockFileInput(cliInput) {
		pkgs, err := stew.ReadStewLockFileContents(cliInput)
		stew.CatchAndExit(err)
		if len(pkgs) == 0 {
			stew.CatchAndExit(stew.EmptyCLIInputError{})
		}
		err = installFromLockFile(pkgs, userOS, userArch, systemInfo)
		stew.CatchAndExit(err)
	} else if filepath.Base(cliInput) == "Stewfile" {
		pkgs, err := stew.ReadStewfileContents(cliInput)
		stew.CatchAndExit(err)
		if len(pkgs) == 0 {
			stew.CatchAndExit(stew.EmptyCLIInputError{})
		}
		err = installFromStewfile(pkgs, userOS, userArch, systemInfo)
		stew.CatchAndExit(err)
	} else {
		pkg, err := stew.ParseCLIInput(cliInput)
		stew.CatchAndExit(err)
		err = installOne(pkg, userOS, userArch, systemInfo, false)
		stew.CatchAndExit(err)
	}
}

func installOne(pkg stew.PackageData, userOS, userArch string, systemInfo stew.SystemInfo, installingFromLockFile bool) error {
	sp := constants.LoadingSpinner

	stewBinPath := systemInfo.StewBinPath
	stewPkgPath := systemInfo.StewPkgPath
	stewLockFilePath := systemInfo.StewLockFilePath
	stewTmpPath := systemInfo.StewTmpPath

	source := pkg.Source
	owner := pkg.Owner
	repo := pkg.Repo
	tag := pkg.Tag
	asset := pkg.Asset
	desiredBinaryRename := pkg.Binary
	expectedBinaryHash := pkg.BinaryHash
	downloadURL := pkg.URL

	lockFile, err := stew.NewLockFile(stewLockFilePath, userOS, userArch)
	if err != nil {
		return err
	}

	if err = os.RemoveAll(stewTmpPath); err != nil {
		return err
	}
	if err = os.MkdirAll(stewTmpPath, 0755); err != nil {
		return err
	}

	var githubProject stew.GithubProject
	if source == "github" {
		fmt.Println(constants.GreenColor(owner + "/" + repo))
		sp.Start()
		githubProject, err = stew.NewGithubProject(owner, repo)
		sp.Stop()
		if err != nil {
			return err
		}

		releaseTags, err := stew.GetGithubReleasesTags(githubProject)
		if err != nil {
			return err
		}

		if tag == "" || tag == "latest" {
			// Find first non-prerelease tag
			for _, release := range githubProject.Releases {
				if !release.Prerelease {
					tag = release.TagName
					break
				}
			}
		}

		tagIndex, tagFound := stew.Contains(releaseTags, tag)
		if !tagFound {
			tag, err = stew.WarningPromptSelect(fmt.Sprintf("Could not find a release with the tag %v - please select a release:", constants.YellowColor(tag)), releaseTags)
			if err != nil {
				return err
			}
			tagIndex, _ = stew.Contains(releaseTags, tag)
		}

		releaseAssets, err := stew.GetGithubReleasesAssets(githubProject, tag)
		if err != nil {
			return err
		}

		if asset == "" {
			asset, err = stew.DetectAsset(userOS, userArch, releaseAssets)
		}
		if err != nil {
			return err
		}

		assetIndex, assetFound := stew.Contains(releaseAssets, asset)
		if !assetFound {
			asset, err = stew.WarningPromptSelect(fmt.Sprintf("Could not find the asset %v - please select an asset:", constants.YellowColor(asset)), releaseAssets)
			if err != nil {
				return err
			}
			assetIndex, _ = stew.Contains(releaseAssets, asset)
		}

		downloadURL = githubProject.Releases[tagIndex].Assets[assetIndex].DownloadURL
	} else {
		fmt.Println(constants.GreenColor(asset))
	}

	downloadPath := filepath.Join(stewPkgPath, asset)
	err = stew.DownloadFile(downloadPath, downloadURL)
	if err != nil {
		return err
	}
	fmt.Printf("✅ Downloaded %v to %v\n", constants.GreenColor(asset), constants.GreenColor(stewPkgPath))

	binaryName, binaryHash, err := stew.InstallBinary(downloadPath, repo, systemInfo, &lockFile, installingFromLockFile, desiredBinaryRename, expectedBinaryHash)
	if err != nil {
		if err := os.RemoveAll(downloadPath); err != nil {
			return err
		}
		return err
	}

	var packageData stew.PackageData
	if source == "github" {
		packageData = stew.PackageData{
			Source:     "github",
			Owner:      githubProject.Owner,
			Repo:       githubProject.Repo,
			Tag:        tag,
			Asset:      asset,
			Binary:     binaryName,
			URL:        downloadURL,
			BinaryHash: binaryHash,
		}
	} else {
		packageData = stew.PackageData{
			Source:     "other",
			Owner:      "",
			Repo:       "",
			Tag:        "",
			Asset:      asset,
			Binary:     binaryName,
			URL:        downloadURL,
			BinaryHash: binaryHash,
		}
	}

	indexInLockFile, binaryFoundInLockFile := stew.FindBinaryInLockFile(lockFile, binaryName)
	if installingFromLockFile && binaryFoundInLockFile {
		lockFile.Packages[indexInLockFile] = packageData
	} else {
		lockFile.Packages = append(lockFile.Packages, packageData)
	}

	err = stew.WriteLockFileJSON(lockFile, stewLockFilePath)
	if err != nil {
		return err
	}

	fmt.Printf("✨ Successfully installed the %v binary in %v\n", constants.GreenColor(binaryName), constants.GreenColor(stewBinPath))
	return nil
}

// isLockFileInput checks if the input is a lockfile: Stewfile.lock.json, or an existing file whose name ends with .lock.json
func isLockFileInput(cliInput string) bool {
	if filepath.Base(cliInput) == "Stewfile.lock.json" {
		return true
	}
	if !strings.HasSuffix(cliInput, ".lock.json") {
		return false
	}
	fileInfo, err := os.Stat(cliInput)
	return err == nil && fileInfo.Mode().IsRegular()
}

func installFromLockFile(pkgs []stew.PackageData, userOS, userArch string, systemInfo stew.SystemInfo) error {
	return installAll(pkgs, userOS, userArch, systemInfo, true)
}

func installFromStewfile(pkgs []stew.PackageData, userOS, userArch string, systemInfo stew.SystemInfo) error {
	return installAll(pkgs, userOS, userArch, systemInfo, false)
}

// installAll installs each package. After an error, it continues with the next package.
// It returns an InstallFailedError that names the packages that could not be installed.
func installAll(pkgs []stew.PackageData, userOS, userArch string, systemInfo stew.SystemInfo, installingFromLockFile bool) error {
	failedPackages := []string{}
	for _, pkg := range pkgs {
		err := installOne(pkg, userOS, userArch, systemInfo, installingFromLockFile)
		if err == nil {
			continue
		}
		fmt.Fprintln(os.Stderr, err)
		// The user chose not to overwrite the binary, thus this is not a failure
		if _, aborted := err.(stew.AbortBinaryOverwriteError); aborted {
			continue
		}
		failedPackages = append(failedPackages, packageName(pkg))
	}
	if len(failedPackages) > 0 {
		return stew.InstallFailedError{Packages: failedPackages, Total: len(pkgs)}
	}
	return nil
}

func packageName(pkg stew.PackageData) string {
	if pkg.Source == "github" {
		return pkg.Owner + "/" + pkg.Repo
	}
	if pkg.Binary != "" {
		return pkg.Binary
	}
	return pkg.Asset
}
