package cmd

import (
	"fmt"

	stew "github.com/marwanhawari/stew/lib"
)

// List is executed when you run `stew list`
func List(cliTagsFlag bool, cliAssetsFlag bool) {

	stewPath, err := stew.GetStewPath()
	stew.CatchAndExit(err)
	systemInfo := stew.NewSystemInfo(stewPath)

	userOS := systemInfo.Os
	userArch := systemInfo.Arch
	stewLockFilePath := systemInfo.StewLockFilePath

	lockFile, err := stew.NewLockFile(stewLockFilePath, userOS, userArch)
	stew.CatchAndExit(err)

	if len(lockFile.Packages) == 0 {
		return
	}

	for _, pkg := range lockFile.Packages {
		switch pkg.Source {
		case "other":
			fmt.Println(pkg.URL)
		case "github":
			if !cliTagsFlag && !cliAssetsFlag {
				fmt.Println(pkg.Owner + "/" + pkg.Repo)
			} else if cliTagsFlag && !cliAssetsFlag {
				fmt.Println(pkg.Owner + "/" + pkg.Repo + "@" + pkg.Tag)
			} else {
				fmt.Println(pkg.Owner + "/" + pkg.Repo + "@" + pkg.Tag + "::" + pkg.Asset)
			}
		}
	}
}
