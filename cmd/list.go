package cmd

import (
	"fmt"
	"os"

	"github.com/gookit/color"
	"github.com/marwanhawari/stew/constants"
	stew "github.com/marwanhawari/stew/lib"
	"golang.org/x/term"
)

// List is executed when you run `stew list`
func List(cliTagsFlag bool) {
	if !term.IsTerminal(int(os.Stdout.Fd())) {
		color.Disable()
	}

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
			fmt.Println(constants.GreenColor(pkg.Binary+":") + pkg.URL)
		case "github":
			defaultLine := constants.GreenColor(pkg.Binary+":") + pkg.Owner + "/" + pkg.Repo
			if cliTagsFlag {
				fmt.Println(defaultLine + "@" + pkg.Tag)
			} else {
				fmt.Println(defaultLine)
			}
		}
	}
}
