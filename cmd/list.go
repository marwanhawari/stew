package cmd

import (
	"fmt"
	"os"

	stew "github.com/marwanhawari/stew/lib"
)

// List is executed when you run `stew list`
func List(cliTagsFlag bool, cliAssetsFlag bool) {

	userOS, userArch, _, systemInfo, err := stew.Initialize()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

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
