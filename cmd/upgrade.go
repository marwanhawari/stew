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
)

// Upgrade is executed when you run `stew upgrade`
func Upgrade(cliFlag bool, binaryName string) {
	rt := errs.Strip(config.Initialize())

	if cliFlag && binaryName != "" {
		errs.MaybeExit(stew.CLIFlagAndInputError{})
	} else if !cliFlag {
		errs.MaybeExit(stew.ValidateCLIInput(binaryName))
	}

	sp := progress.Spinner(rt)

	lockFile := errs.Strip(stew.NewLockFile(rt))

	errs.MaybeExit(os.RemoveAll(rt.TmpPath))
	errs.MaybeExit(os.MkdirAll(rt.TmpPath, 0o755))

	if len(lockFile.Packages) == 0 {
		errs.MaybeExit(stew.NoBinariesInstalledError{})
	}

	var binaryFound bool
	for index, pkg := range lockFile.Packages {

		var upgradeCondition bool
		if cliFlag || pkg.Binary == binaryName {
			upgradeCondition = true
		}

		if upgradeCondition {
			rt.Println(constants.GreenColor(pkg.Binary))
			binaryFound = true

			if pkg.Source == "other" {
				rt.Println(stew.InstalledFromURLError{Binary: pkg.Binary})
				continue
			}
			owner := pkg.Owner
			repo := pkg.Repo

			sp.Start()
			githubProject, err := stew.NewGithubProject(rt.Config, owner, repo)
			sp.Stop()
			errs.MaybeExit(err)

			// This will make sure that there are any tags at all
			_, err = stew.GetGithubReleasesTags(githubProject)
			errs.MaybeExit(err)

			// Get the latest tag
			tagIndex := 0
			tag := githubProject.Releases[tagIndex].TagName

			if pkg.Tag == tag {
				rt.Println(stew.AlreadyInstalledLatestTagError{Tag: tag})
				continue
			}

			// Make sure there are any assets at all
			releaseAssets, err := stew.GetGithubReleasesAssets(githubProject, tag)
			errs.MaybeExit(err)

			asset, err := stew.DetectAsset(rt, releaseAssets)
			errs.MaybeExit(err)
			assetIndex, _ := stew.Contains(releaseAssets, asset)

			downloadURL := githubProject.Releases[tagIndex].Assets[assetIndex].DownloadURL

			downloadPath := filepath.Join(rt.PkgPath, asset)
			downloadPathExists, err := pathutil.Exists(downloadPath)
			errs.MaybeExit(err)
			if downloadPathExists {
				errs.MaybeExit(stew.AssetAlreadyDownloadedError{Asset: asset})
			} else {
				err = stew.DownloadFile(rt, downloadPath, downloadURL)
				errs.MaybeExit(err)
				rt.Printf("✅ Downloaded %v to %v\n",
					constants.GreenColor(asset), constants.GreenColor(rt.PkgPath))
			}

			_, err = stew.InstallBinary(rt, stew.Installation{
				DownloadedFilePath: downloadPath,
				Repo:               repo,
				BinaryName:         pkg.Binary,
			}, &lockFile, true)
			if err != nil {
				errs.MaybeExit(os.RemoveAll(downloadPath))
				errs.MaybeExit(err)
			}

			lockFile.Packages[index].Tag = tag
			lockFile.Packages[index].Asset = asset
			lockFile.Packages[index].URL = downloadURL

			err = stew.WriteLockFileJSON(rt, lockFile, rt.LockPath)
			errs.MaybeExit(err)

			rt.Printf("✨ Successfully upgraded the %v binary from %v to %v\n",
				constants.GreenColor(pkg.Binary), constants.GreenColor(pkg.Tag),
				constants.GreenColor(tag))

		}
	}
	if !cliFlag && !binaryFound {
		errs.MaybeExit(stew.BinaryNotInstalledError{Binary: binaryName})
	}
}
