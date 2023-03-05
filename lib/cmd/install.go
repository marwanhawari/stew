package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/marwanhawari/stew/constants"
	stew "github.com/marwanhawari/stew/lib"
	"github.com/marwanhawari/stew/lib/config"
	"github.com/marwanhawari/stew/lib/pathutil"
	"github.com/marwanhawari/stew/lib/ui/progress"
	"github.com/marwanhawari/stew/lib/ui/prompt"
)

// Install is executed when you run `stew install`.
func Install(rt config.Runtime, cliInputs []string, batchMode bool) error {
	var err error
	for _, cliInput := range cliInputs {
		if strings.Contains(cliInput, "Stewfile.lock.json") {
			if cliInputs, err = stew.ReadStewLockFileContents(cliInput); err != nil {
				return err
			}
			break
		}

		if strings.Contains(cliInput, "Stewfile") {
			if cliInputs, err = stew.ReadStewfileContents(cliInput); err != nil {
				return err
			}
			break
		}
	}

	if len(cliInputs) == 0 {
		return stew.EmptyCLIInputError{}
	}

	for _, cliInput := range cliInputs {
		sp := progress.Spinner(rt)

		var input stew.CLIInput
		if input, err = stew.ParseCLIInput(cliInput); err != nil {
			return err
		}

		var (
			owner       = input.Owner
			repo        = input.Repo
			tag         = input.Tag
			asset       = input.Asset
			downloadURL = input.DownloadURL
		)

		var lockFile stew.LockFile
		if lockFile, err = stew.NewLockFile(rt); err != nil {
			return err
		}

		var githubProject stew.GithubProject
		if input.IsGithubInput {
			rt.Println(constants.GreenColor(owner + "/" + repo))
			if err = func() error {
				sp.Start()
				defer sp.Stop()
				githubProject, err = stew.NewGithubProject(rt.Config, owner, repo)
				return err
			}(); err != nil {
				return err
			}

			// This will make sure that there are any tags at all
			var releaseTags []string
			if releaseTags, err = stew.GetGithubReleasesTags(githubProject); err != nil {
				return err
			}

			if tag == "" || tag == "latest" {
				tag = githubProject.Releases[0].TagName
			}

			// Need to make sure user input tag is in the tags
			tagIndex, tagFound := stew.Contains(releaseTags, tag)
			if !tagFound {
				if tag, err = prompt.SelectWarn(rt, fmt.Sprintf(
					"Could not find a release with the tag %v - please select a release:",
					constants.YellowColor(tag)), releaseTags); err != nil {
					return err
				}
				tagIndex, _ = stew.Contains(releaseTags, tag)
			}

			// Make sure there are any assets at all
			var releaseAssets []string
			if releaseAssets, err = stew.GetGithubReleasesAssets(githubProject, tag); err != nil {
				return err
			}

			if asset == "" {
				if asset, err = stew.DetectAsset(rt, releaseAssets); err != nil {
					return err
				}
			}

			assetIndex, assetFound := stew.Contains(releaseAssets, asset)
			if !assetFound {
				if asset, err = prompt.SelectWarn(rt, fmt.Sprintf(
					"Could not find the asset %v - please select an asset:",
					constants.YellowColor(asset)), releaseAssets); err != nil {
					return err
				}
				assetIndex, _ = stew.Contains(releaseAssets, asset)
			}

			downloadURL = githubProject.Releases[tagIndex].Assets[assetIndex].DownloadURL
		} else {
			rt.Println(constants.GreenColor(asset))
		}

		downloadPath := filepath.Join(rt.PkgPath, asset)
		var downloadPathExists bool
		if downloadPathExists, err = pathutil.Exists(downloadPath); err != nil {
			return err
		}
		if downloadPathExists {
			rt.Println(stew.AssetAlreadyDownloadedError{Asset: asset})
			continue
		} else {
			if err = stew.DownloadFile(rt, downloadPath, downloadURL); err != nil {
				return err
			}
			rt.Printf("✅ Downloaded %v to %v\n",
				constants.GreenColor(asset), constants.GreenColor(rt.PkgPath))
		}

		var binaryName string
		if binaryName, err = stew.InstallBinary(rt, stew.Installation{
			DownloadedFilePath: downloadPath,
			Repo:               repo,
			BinaryName:         input.BinaryName,
			BatchMode:          batchMode,
		}, &lockFile, false); err != nil {
			_ = os.RemoveAll(downloadPath)
			return err
		}

		var packageData stew.PackageData
		if input.IsGithubInput {
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

		if err = stew.WriteLockFileJSON(rt, lockFile, rt.LockPath); err != nil {
			return err
		}

		rt.Printf("✨ Successfully installed the %v binary in %v\n",
			constants.GreenColor(binaryName), constants.GreenColor(rt.StewBinPath))
	}
	return err
}
