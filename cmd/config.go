package cmd

import (
	"os"
	"runtime"

	"github.com/marwanhawari/stew/constants"
	"github.com/marwanhawari/stew/lib/config"
	"github.com/marwanhawari/stew/lib/errs"
	"github.com/marwanhawari/stew/lib/pathutil"
	"github.com/marwanhawari/stew/lib/ui/terminal"
)

// Config is executed when you run `stew config`.
func Config() {
	userOS := runtime.GOOS
	prt := terminal.Standard()
	configPath, err := config.FilePath(userOS)
	errs.MaybeExit(err)
	configExists := errs.Strip(pathutil.Exists(configPath))

	if !configExists {
		_, err = config.NewConfig(prt, userOS)
		errs.MaybeExit(err)
		return
	}

	defaultStewPath, err := config.DefaultStewPath(userOS)
	errs.MaybeExit(err)
	defaultBinPath, err := config.DefaultBinPath(userOS)
	errs.MaybeExit(err)

	newStewPath, newStewBinPath, err := config.PromptConfig(prt, defaultStewPath, defaultBinPath)
	errs.MaybeExit(err)

	cfg := config.Config{StewPath: newStewPath, StewBinPath: newStewBinPath}
	err = cfg.WriteStewConfigJSON(configPath)
	errs.MaybeExit(err)

	prt.Printf("📄 Updated %v\n", constants.GreenColor(configPath))

	pathVariable := os.Getenv("PATH")
	config.ValidateBinPath(prt, newStewBinPath, pathVariable)
}
