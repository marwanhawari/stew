package cmd

import (
	"os"
	"path/filepath"

	"github.com/marwanhawari/stew/constants"
	stew "github.com/marwanhawari/stew/lib"
	"github.com/marwanhawari/stew/lib/config"
	"github.com/marwanhawari/stew/lib/errs"
	"github.com/marwanhawari/stew/lib/pathutil"
	"github.com/marwanhawari/stew/lib/ui/progress"
	"github.com/marwanhawari/stew/lib/ui/prompt"
)

// Browse is executed when you run `stew browse`
func Browse(cliInput string) {
	rt := errs.Strip(config.Initialize())

	sp := progress.Spinner(rt)

	stewBinPath := rt.StewBinPath
	stewPkgPath := rt.PkgPath
	stewLockFilePath := rt.LockPath
	stewTmpPath := rt.TmpPath

	parsedInput := errs.Strip(stew.ParseCLIInput(cliInput))

	owner := parsedInput.Owner
	repo := parsedInput.Repo

	lockFile := errs.Strip(stew.NewLockFile(rt))

	errs.MaybeExit(os.RemoveAll(stewTmpPath))
	errs.MaybeExit(os.MkdirAll(stewTmpPath, 0o755))

	rt.Println(constants.GreenColor(owner + "/" + repo))
	sp.Start()
	githubProject := errs.Strip(stew.NewGithubProject(rt.Config, owner, repo))
	sp.Stop()

	releaseTags := errs.Strip(stew.GetGithubReleasesTags(githubProject))
	tag := errs.Strip(prompt.Select(rt, "Choose a release tag:", releaseTags))
	tagIndex, _ := stew.Contains(releaseTags, tag)

	releaseAssets := errs.Strip(stew.GetGithubReleasesAssets(githubProject, tag))
	asset := errs.Strip(prompt.Select(rt, "Download and install an asset", releaseAssets))
	assetIndex, _ := stew.Contains(releaseAssets, asset)

	downloadURL := githubProject.Releases[tagIndex].Assets[assetIndex].DownloadURL
	downloadPath := filepath.Join(stewPkgPath, asset)
	downloadPathExists := errs.Strip(pathutil.Exists(downloadPath))

	if downloadPathExists {
		errs.MaybeExit(stew.AssetAlreadyDownloadedError{Asset: asset})
	} else {
		errs.MaybeExit(stew.DownloadFile(rt, downloadPath, downloadURL))
		rt.Printf("✅ Downloaded %v to %v\n", constants.GreenColor(asset), constants.GreenColor(stewPkgPath))
	}

	binaryName := func() string {
		var (
			binaryName string
			err        error
		)
		if binaryName, err = stew.InstallBinary(rt, stew.Installation{
			DownloadedFilePath: downloadPath,
			Repo:               repo,
			BinaryName:         parsedInput.BinaryName,
		}, &lockFile, false); err != nil {
			errs.MaybeExit(os.RemoveAll(downloadPath))
			errs.MaybeExit(err)
		}
		return binaryName
	}()

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

	errs.MaybeExit(stew.WriteLockFileJSON(rt, lockFile, stewLockFilePath))

	rt.Printf("✨ Successfully installed the %v binary in %v\n",
		constants.GreenColor(binaryName), constants.GreenColor(stewBinPath))

}
