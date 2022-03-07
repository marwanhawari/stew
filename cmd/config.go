package cmd

import (
	"fmt"
	"os"

	"github.com/marwanhawari/stew/constants"
	stew "github.com/marwanhawari/stew/lib"
)

func Config() {

	userOS, _, stewConfig, _, err := stew.Initialize()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	stewConfigFilePath, err := stew.GetStewConfigFilePath(userOS)
	stew.CatchAndExit(err)

	inputStewPath, err := stew.PromptInput("Set the stewPath. This will contain all stew data other than the binaries.", stewConfig.StewPath)
	stew.CatchAndExit(err)
	inputStewBinPath, err := stew.PromptInput("Set the stewBinPath. This is where the binaries will be installed by stew.", stewConfig.StewBinPath)
	stew.CatchAndExit(err)

	fullStewPath, err := stew.ResolveTilde(inputStewPath)
	stew.CatchAndExit(err)
	fullStewBinPath, err := stew.ResolveTilde(inputStewBinPath)
	stew.CatchAndExit(err)

	newStewConfig := stew.StewConfig{StewPath: fullStewPath, StewBinPath: fullStewBinPath}
	err = stew.WriteStewConfigJSON(newStewConfig, stewConfigFilePath)
	stew.CatchAndExit(err)

	fmt.Printf("📄 Updated %v\n", constants.GreenColor(stewConfigFilePath))
}
