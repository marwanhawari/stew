package cmd

import (
	stew "github.com/marwanhawari/stew/lib"
	"github.com/marwanhawari/stew/lib/config"
	"github.com/marwanhawari/stew/lib/errs"
)

// List is executed when you run `stew list`.
func List(cliTagsFlag bool, cliAssetsFlag bool) {
	rt := errs.Strip(config.Initialize())

	lockFile := errs.Strip(stew.NewLockFile(rt))

	if len(lockFile.Packages) == 0 {
		return
	}

	for _, pkg := range lockFile.Packages {
		switch pkg.Source {
		case "other":
			rt.Println(pkg.URL)
		case "github":
			if !cliTagsFlag && !cliAssetsFlag {
				rt.Println(pkg.Owner + "/" + pkg.Repo)
			} else if cliTagsFlag && !cliAssetsFlag {
				rt.Println(pkg.Owner + "/" + pkg.Repo + "@" + pkg.Tag)
			} else {
				rt.Println(pkg.Owner + "/" + pkg.Repo + "@" + pkg.Tag + "::" + pkg.Asset + "!!" + pkg.Binary)
			}
		}
	}
}
