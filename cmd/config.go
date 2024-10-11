package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/marwanhawari/stew/constants"
	stew "github.com/marwanhawari/stew/lib"
)

func Config() {

	userOS := runtime.GOOS
	stewConfigFilePath, err := stew.GetStewConfigFilePath(userOS)
	stew.CatchAndExit(err)
	configExists, err := stew.PathExists(stewConfigFilePath)
	stew.CatchAndExit(err)

	if !configExists {
		_, err := stew.NewStewConfig(userOS)
		stew.CatchAndExit(err)
		return
	}
	_, _, _, systemInfo, _ := stew.Initialize()
	installedPackages, err := stew.ReadStewLockFileContents(systemInfo.StewLockFilePath)
	stew.CatchAndExit(err)

	config, err := stew.ReadStewConfigJSON(stewConfigFilePath)
	stew.CatchAndExit(err)
	newStewPath, newStewBinPath, newExcludedFromUpgradeAll, err := stew.PromptConfig(config.StewPath, config.StewBinPath, installedPackages, config.ExcludedFromUpgradeAll)
	stew.CatchAndExit(err)

	newStewConfig := stew.StewConfig{StewPath: newStewPath, StewBinPath: newStewBinPath, ExcludedFromUpgradeAll: newExcludedFromUpgradeAll}
	err = stew.WriteStewConfigJSON(newStewConfig, stewConfigFilePath)
	stew.CatchAndExit(err)

	fmt.Printf("📄 Updated %v\n", constants.GreenColor(stewConfigFilePath))

	pathVariable := os.Getenv("PATH")
	stew.ValidateStewBinPath(newStewBinPath, pathVariable)
}
