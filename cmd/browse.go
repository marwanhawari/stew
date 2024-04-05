package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/marwanhawari/stew/constants"
	stew "github.com/marwanhawari/stew/lib"
)

// Browse is executed when you run `stew browse`
func Browse(cliInput string) {

	userOS, userArch, _, systemInfo, err := stew.Initialize()
	stew.CatchAndExit(err)

	sp := constants.LoadingSpinner

	stewBinPath := systemInfo.StewBinPath
	stewPkgPath := systemInfo.StewPkgPath
	stewLockFilePath := systemInfo.StewLockFilePath
	stewTmpPath := systemInfo.StewTmpPath

	parsedInput, err := stew.ParseCLIInput(cliInput)
	stew.CatchAndExit(err)

	owner := parsedInput.Owner
	repo := parsedInput.Repo

	lockFile, err := stew.NewLockFile(stewLockFilePath, userOS, userArch)
	stew.CatchAndExit(err)

	err = os.RemoveAll(stewTmpPath)
	stew.CatchAndExit(err)
	err = os.MkdirAll(stewTmpPath, 0755)
	stew.CatchAndExit(err)

	fmt.Println(constants.GreenColor(owner + "/" + repo))
	sp.Start()
	githubProject, err := stew.NewGithubProject(owner, repo)
	sp.Stop()
	stew.CatchAndExit(err)

	releaseTags, err := stew.GetGithubReleasesTags(githubProject)
	stew.CatchAndExit(err)
	tag, err := stew.PromptSelect("Choose a release tag:", releaseTags)
	stew.CatchAndExit(err)
	tagIndex, _ := stew.Contains(releaseTags, tag)

	releaseAssets, err := stew.GetGithubReleasesAssets(githubProject, tag)
	stew.CatchAndExit(err)
	asset, err := stew.PromptSelect("Download and install an asset", releaseAssets)
	stew.CatchAndExit(err)
	assetIndex, _ := stew.Contains(releaseAssets, asset)

	downloadURL := githubProject.Releases[tagIndex].Assets[assetIndex].DownloadURL
	downloadPath := filepath.Join(stewPkgPath, asset)
	err = stew.DownloadFile(downloadPath, downloadURL)
	stew.CatchAndExit(err)
	fmt.Printf("✅ Downloaded %v to %v\n", constants.GreenColor(asset), constants.GreenColor(stewPkgPath))

	binaryName, err := stew.InstallBinary(downloadPath, repo, systemInfo, &lockFile, false, "")
	if err != nil {
		os.RemoveAll(downloadPath)
		stew.CatchAndExit(err)
	}

	packageData := stew.PackageData{
		Source: "github",
		Owner:  githubProject.Owner,
		Repo:   githubProject.Repo,
		Tag:    tag,
		Asset:  asset,
		Binary: binaryName,
		URL:    downloadURL,
	}

	lockFile.Packages = append(lockFile.Packages, packageData)

	err = stew.WriteLockFileJSON(lockFile, stewLockFilePath)
	stew.CatchAndExit(err)

	fmt.Printf("✨ Successfully installed the %v binary in %v\n", constants.GreenColor(binaryName), constants.GreenColor(stewBinPath))

}
