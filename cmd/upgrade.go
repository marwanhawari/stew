package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/marwanhawari/stew/constants"
	stew "github.com/marwanhawari/stew/lib"
)

// Upgrade is executed when you run `stew upgrade`
func Upgrade(upgradeAllCliFlag bool, binaryName string) {

	userOS, userArch, _, systemInfo, err := stew.Initialize()
	stew.CatchAndExit(err)

	if upgradeAllCliFlag && binaryName != "" {
		stew.CatchAndExit(stew.CLIFlagAndInputError{})
	} else if !upgradeAllCliFlag {
		err := stew.ValidateCLIInput(binaryName)
		stew.CatchAndExit(err)
	}

	stewTmpPath := systemInfo.StewTmpPath
	stewLockFilePath := systemInfo.StewLockFilePath

	lockFile, err := stew.NewLockFile(stewLockFilePath, userOS, userArch)
	stew.CatchAndExit(err)

	err = os.RemoveAll(stewTmpPath)
	stew.CatchAndExit(err)
	err = os.MkdirAll(stewTmpPath, 0755)
	stew.CatchAndExit(err)

	if len(lockFile.Packages) == 0 {
		stew.CatchAndExit(stew.NoBinariesInstalledError{})
	}

	if upgradeAllCliFlag {
		upgradeAll(userOS, userArch, lockFile, systemInfo)
	} else {
		err := upgradeOne(binaryName, userOS, userArch, lockFile, systemInfo)
		stew.CatchAndExit(err)
	}
}

func upgradeOne(binaryName, userOS, userArch string, lockFile stew.LockFile, systemInfo stew.SystemInfo) error {
	sp := constants.LoadingSpinner
	stewPkgPath := systemInfo.StewPkgPath
	stewLockFilePath := systemInfo.StewLockFilePath

	indexInLockFile, binaryFoundInLockFile := stew.FindBinaryInLockFile(lockFile, binaryName)
	if !binaryFoundInLockFile {
		return stew.BinaryNotInstalledError{Binary: binaryName}
	}

	pkg := lockFile.Packages[indexInLockFile]
	fmt.Println(constants.GreenColor(pkg.Binary))
	if pkg.Source == "other" {
		return stew.InstalledFromURLError{Binary: pkg.Binary}
	}
	owner := pkg.Owner
	repo := pkg.Repo

	sp.Start()
	githubProject, err := stew.NewGithubProject(owner, repo)
	sp.Stop()
	if err != nil {
		return err
	}

	// This will make sure that there are any tags at all
	_, err = stew.GetGithubReleasesTags(githubProject)
	if err != nil {
		return err
	}

	// Get the latest tag
	tagIndex := 0
	tag := githubProject.Releases[tagIndex].TagName

	if pkg.Tag == tag {
		return stew.AlreadyInstalledLatestTagError{Tag: tag}
	}

	// Make sure there are any assets at all
	releaseAssets, err := stew.GetGithubReleasesAssets(githubProject, tag)
	if err != nil {
		return err
	}

	asset, err := stew.DetectAsset(userOS, userArch, releaseAssets)
	if err != nil {
		return err
	}
	assetIndex, _ := stew.Contains(releaseAssets, asset)
	downloadURL := githubProject.Releases[tagIndex].Assets[assetIndex].DownloadURL
	downloadPath := filepath.Join(stewPkgPath, asset)
	err = stew.DownloadFile(downloadPath, downloadURL)
	if err != nil {
		return err
	}
	fmt.Printf("✅ Downloaded %v to %v\n", constants.GreenColor(asset), constants.GreenColor(stewPkgPath))

	_, binaryHash, err := stew.InstallBinary(downloadPath, repo, systemInfo, &lockFile, true, pkg.Binary, "")
	if err != nil {
		if err := os.RemoveAll(downloadPath); err != nil {
			return err
		}
		return err
	}

	lockFile.Packages[indexInLockFile].Tag = tag
	lockFile.Packages[indexInLockFile].Asset = asset
	lockFile.Packages[indexInLockFile].URL = downloadURL
	lockFile.Packages[indexInLockFile].BinaryHash = binaryHash
	if err := stew.WriteLockFileJSON(lockFile, stewLockFilePath); err != nil {
		return err
	}

	fmt.Printf("✨ Successfully upgraded the %v binary from %v to %v\n", constants.GreenColor(pkg.Binary), constants.GreenColor(pkg.Tag), constants.GreenColor(tag))
	return nil
}

func upgradeAll(userOS, userArch string, lockFile stew.LockFile, systemInfo stew.SystemInfo) {
	stewConfigFilePath, err := stew.GetStewConfigFilePath(userOS)
	stew.CatchAndExit(err)
	stewConfig, err := stew.ReadStewConfigJSON(stewConfigFilePath)
	stew.CatchAndExit(err)
	excludedPackages := stewConfig.ExcludedFromUpgradeAll

	for _, pkg := range lockFile.Packages {
		if pkg.Source != "other" {
			// Skip upgrading explicitly excluded packages
			pkgName := stew.GetPackageDisplayName(pkg, false)
			if _, packageIsExcluded := stew.Contains(excludedPackages, pkgName); packageIsExcluded {
				fmt.Printf("%v (excluded)\n", constants.YellowColor(pkgName))
				continue
			}
		}
		if err := upgradeOne(pkg.Binary, userOS, userArch, lockFile, systemInfo); err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
	}
}
