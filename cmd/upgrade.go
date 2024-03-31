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

	sp := constants.LoadingSpinner

	stewPkgPath := systemInfo.StewPkgPath
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

	var binaryFound bool
	for index, pkg := range lockFile.Packages {
		shouldUpgradePackage := (upgradeAllCliFlag || pkg.Binary == binaryName)
		if !shouldUpgradePackage {
			continue
		}
		fmt.Println(constants.GreenColor(pkg.Binary))
		binaryFound = true

		if pkg.Source == "other" {
			fmt.Fprintln(os.Stderr, stew.InstalledFromURLError{Binary: pkg.Binary})
			continue
		}
		owner := pkg.Owner
		repo := pkg.Repo

		sp.Start()
		githubProject, err := stew.NewGithubProject(owner, repo)
		sp.Stop()
		stew.CatchAndExit(err)

		// This will make sure that there are any tags at all
		_, err = stew.GetGithubReleasesTags(githubProject)
		stew.CatchAndExit(err)

		// Get the latest tag
		tagIndex := 0
		tag := githubProject.Releases[tagIndex].TagName

		if pkg.Tag == tag {
			fmt.Fprintln(os.Stderr, stew.AlreadyInstalledLatestTagError{Tag: tag})
			continue
		}

		// Make sure there are any assets at all
		releaseAssets, err := stew.GetGithubReleasesAssets(githubProject, tag)
		stew.CatchAndExit(err)

		asset, err := stew.DetectAsset(userOS, userArch, releaseAssets)
		stew.CatchAndExit(err)
		assetIndex, _ := stew.Contains(releaseAssets, asset)

		downloadURL := githubProject.Releases[tagIndex].Assets[assetIndex].DownloadURL

		downloadPath := filepath.Join(stewPkgPath, asset)
		downloadPathExists, err := stew.PathExists(downloadPath)
		stew.CatchAndExit(err)
		if downloadPathExists {
			stew.CatchAndExit(stew.AssetAlreadyDownloadedError{Asset: asset})
		} else {
			err = stew.DownloadFile(downloadPath, downloadURL)
			stew.CatchAndExit(err)
			fmt.Printf("✅ Downloaded %v to %v\n", constants.GreenColor(asset), constants.GreenColor(stewPkgPath))
		}

		_, err = stew.InstallBinary(downloadPath, repo, systemInfo, &lockFile, true)
		if err != nil {
			os.RemoveAll(downloadPath)
			stew.CatchAndExit(err)
		}

		lockFile.Packages[index].Tag = tag
		lockFile.Packages[index].Asset = asset
		lockFile.Packages[index].URL = downloadURL

		err = stew.WriteLockFileJSON(lockFile, stewLockFilePath)
		stew.CatchAndExit(err)

		fmt.Printf("✨ Successfully upgraded the %v binary from %v to %v\n", constants.GreenColor(pkg.Binary), constants.GreenColor(pkg.Tag), constants.GreenColor(tag))
	}
	if !upgradeAllCliFlag && !binaryFound {
		stew.CatchAndExit(stew.BinaryNotInstalledError{Binary: binaryName})
	}
}
