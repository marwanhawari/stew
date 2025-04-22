package cmd

import (
	"fmt"

	"github.com/marwanhawari/stew/constants"
	stew "github.com/marwanhawari/stew/lib"
)

// List is executed when you run `stew list`
func List(cliTagsFlag bool) {

	userOS, userArch, _, systemInfo, err := stew.Initialize()
	stew.CatchAndExit(err)

	stewLockFilePath := systemInfo.StewLockFilePath

	lockFile, err := stew.NewLockFile(stewLockFilePath, userOS, userArch)
	stew.CatchAndExit(err)

	if len(lockFile.Packages) == 0 {
		return
	}

	for _, pkg := range lockFile.Packages {
		switch pkg.Source {
		case "other":
			fmt.Println(pkg.URL + constants.GreenColor(":"+pkg.Binary))
		case "github":
			defaultLine := pkg.Owner + "/" + pkg.Repo + constants.GreenColor(":"+pkg.Binary)
			if cliTagsFlag {
				fmt.Println(defaultLine + constants.CyanColor("@"+pkg.Tag))
			} else {
				fmt.Println(defaultLine)
			}
		}
	}
}
