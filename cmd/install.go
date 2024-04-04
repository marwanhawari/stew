package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/marwanhawari/stew/constants"
	stew "github.com/marwanhawari/stew/lib"
)

// Install is executed when you run `stew install`
func Install(cliInput string) {
	var err error

	userOS, userArch, _, systemInfo, err := stew.Initialize()
	stew.CatchAndExit(err)

	var pkgs []stew.PackageData
	if filepath.Base(cliInput) == "Stewfile.lock.json" {
		pkgs, err = stew.ReadStewLockFileContents(cliInput)
		stew.CatchAndExit(err)
	} else if filepath.Base(cliInput) == "Stewfile" {
		pkgs, err = stew.ReadStewfileContents(cliInput)
		stew.CatchAndExit(err)
	} else {
		pkg, err := stew.ParseCLIInput(cliInput)
		stew.CatchAndExit(err)
		pkgs = append(pkgs, pkg)
	}

	if len(pkgs) == 0 {
		stew.CatchAndExit(stew.EmptyCLIInputError{})
	}

	for _, pkg := range pkgs {
		installOne(pkg, userOS, userArch, systemInfo)
	}
}

func installOne(pkg stew.PackageData, userOS, userArch string, systemInfo stew.SystemInfo) {
	sp := constants.LoadingSpinner

	stewBinPath := systemInfo.StewBinPath
	stewPkgPath := systemInfo.StewPkgPath
	stewLockFilePath := systemInfo.StewLockFilePath
	stewTmpPath := systemInfo.StewTmpPath

	owner := pkg.Owner
	repo := pkg.Repo
	tag := pkg.Tag
	asset := pkg.Asset
	downloadURL := pkg.URL

	lockFile, err := stew.NewLockFile(stewLockFilePath, userOS, userArch)
	stew.CatchAndExit(err)

	err = os.RemoveAll(stewTmpPath)
	stew.CatchAndExit(err)
	err = os.MkdirAll(stewTmpPath, 0755)
	stew.CatchAndExit(err)

	var githubProject stew.GithubProject
	if pkg.Source == "github" {
		fmt.Println(constants.GreenColor(owner + "/" + repo))
		sp.Start()
		githubProject, err = stew.NewGithubProject(owner, repo)
		sp.Stop()
		stew.CatchAndExit(err)

		releaseTags, err := stew.GetGithubReleasesTags(githubProject)
		stew.CatchAndExit(err)

		if tag == "" || tag == "latest" {
			tag = githubProject.Releases[0].TagName
		}

		tagIndex, tagFound := stew.Contains(releaseTags, tag)
		if !tagFound {
			tag, err = stew.WarningPromptSelect(fmt.Sprintf("Could not find a release with the tag %v - please select a release:", constants.YellowColor(tag)), releaseTags)
			stew.CatchAndExit(err)
			tagIndex, _ = stew.Contains(releaseTags, tag)
		}

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
	err = stew.DownloadFile(downloadPath, downloadURL)
	stew.CatchAndExit(err)
	fmt.Printf("✅ Downloaded %v to %v\n", constants.GreenColor(asset), constants.GreenColor(stewPkgPath))

	binaryName, err := stew.InstallBinary(downloadPath, repo, systemInfo, &lockFile, false)
	if err != nil {
		os.RemoveAll(downloadPath)
		stew.CatchAndExit(err)
	}

	var packageData stew.PackageData
	if pkg.Source == "github" {
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
