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
func Install(cliInputs []string) {
	var err error
	for _, cliInput := range cliInputs {
		if strings.Contains(cliInput, "Stewfile") {
			cliInputs, err = stew.ReadStewfileContents(cliInput)
			stew.CatchAndExit(err)
			break
		}
	}

	if len(cliInputs) == 0 {
		stew.CatchAndExit(stew.EmptyCLIInputError{})
	}

	for _, cliInput := range cliInputs {
		sp := constants.LoadingSpinner

		stewPath, err := stew.GetStewPath()
		stew.CatchAndExit(err)
		systemInfo := stew.NewSystemInfo(stewPath)

		userOS := systemInfo.Os
		userArch := systemInfo.Arch
		stewBinPath := systemInfo.StewBinPath
		stewPkgPath := systemInfo.StewPkgPath
		stewLockFilePath := systemInfo.StewLockFilePath
		stewTmpPath := systemInfo.StewTmpPath

		parsedInput, err := stew.ParseCLIInput(cliInput)
		stew.CatchAndExit(err)

		owner := parsedInput.Owner
		repo := parsedInput.Repo
		tag := parsedInput.Tag
		asset := parsedInput.Asset
		downloadURL := parsedInput.DownloadURL

		lockFile, err := stew.NewLockFile(stewLockFilePath, userOS, userArch)
		stew.CatchAndExit(err)

		err = os.RemoveAll(stewTmpPath)
		stew.CatchAndExit(err)
		err = os.MkdirAll(stewTmpPath, 0755)
		stew.CatchAndExit(err)

		var githubProject stew.GithubProject
		if parsedInput.IsGithubInput {
			fmt.Println(constants.GreenColor(owner + "/" + repo))
			sp.Start()
			githubProject, err = stew.NewGithubProject(owner, repo)
			sp.Stop()
			stew.CatchAndExit(err)

			// This will make sure that there are any tags at all
			releaseTags, err := stew.GetGithubReleasesTags(githubProject)
			stew.CatchAndExit(err)

			if tag == "" || tag == "latest" {
				tag = githubProject.Releases[0].TagName
			}

			// Need to make sure user input tag is in the tags
			tagIndex, tagFound := stew.Contains(releaseTags, tag)
			if !tagFound {
				tag, err = stew.WarningPromptSelect(fmt.Sprintf("Could not find a release with the tag %v - please select a release:", constants.YellowColor(tag)), releaseTags)
				stew.CatchAndExit(err)
				tagIndex, _ = stew.Contains(releaseTags, tag)
			}

			// Make sure there are any assets at all
			releaseAssets, err := stew.GetGithubReleasesAssets(githubProject, tag)
			stew.CatchAndExit(err)

			if asset == "" {
				asset, err = stew.DetectAsset(userOS, userArch, releaseAssets)
			}
			stew.CatchAndExit(err)

			assetIndex, assetFound := stew.Contains(releaseAssets, asset)
			if !assetFound {
				asset, err = stew.WarningPromptSelect(fmt.Sprintf("Could not find the asset %v - please select an asset:", constants.YellowColor(asset)), releaseAssets)
				stew.CatchAndExit(err)
				assetIndex, _ = stew.Contains(releaseAssets, asset)
			}

			downloadURL = githubProject.Releases[tagIndex].Assets[assetIndex].DownloadURL

		} else {
			fmt.Println(constants.GreenColor(asset))
		}

		downloadPath := filepath.Join(stewPkgPath, asset)
		downloadPathExists, err := stew.PathExists(downloadPath)
		stew.CatchAndExit(err)
		if downloadPathExists {
			fmt.Println(stew.AssetAlreadyDownloadedError{Asset: asset})
			continue
		} else {
			err = stew.DownloadFile(downloadPath, downloadURL)
			stew.CatchAndExit(err)
			fmt.Printf("✅ Downloaded %v to %v\n", constants.GreenColor(asset), constants.GreenColor(stewPkgPath))
		}

		binaryName, err := stew.InstallBinary(downloadPath, repo, systemInfo, &lockFile, false)
		if err != nil {
			os.RemoveAll(downloadPath)
			stew.CatchAndExit(err)
		}

		var packageData stew.PackageData
		if parsedInput.IsGithubInput {
			packageData = stew.PackageData{
				Source: "github",
				Owner:  githubProject.Owner,
				Repo:   githubProject.Repo,
				Tag:    tag,
				Asset:  asset,
				Binary: binaryName,
				URL:    downloadURL,
			}
		} else {
			packageData = stew.PackageData{
				Source: "other",
				Owner:  "",
				Repo:   "",
				Tag:    "",
				Asset:  asset,
				Binary: binaryName,
				URL:    downloadURL,
			}
		}

		lockFile.Packages = append(lockFile.Packages, packageData)

		err = stew.WriteLockFileJSON(lockFile, stewLockFilePath)
		stew.CatchAndExit(err)

		fmt.Printf("✨ Successfully installed the %v binary in %v\n", constants.GreenColor(binaryName), constants.GreenColor(stewBinPath))

	}
}
