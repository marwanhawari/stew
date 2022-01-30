package cmd

import (
	"fmt"

	"github.com/marwanhawari/stew/constants"
	stew "github.com/marwanhawari/stew/lib"
)

func Uninstall(cliFlag bool, binaryName string) {
	if cliFlag && binaryName != "" {
		stew.CatchAndExit(stew.UninstallCLIInputError{})
	} else if !cliFlag {
		err := stew.ValidateCLIInput(binaryName)
		stew.CatchAndExit(err)
	}

	systemInfo, err := stew.NewSystemInfo()
	stew.CatchAndExit(err)

	userOS := systemInfo.Os
	userArch := systemInfo.Arch
	stewBinPath := systemInfo.StewBinPath
	stewPkgPath := systemInfo.StewPkgPath
	stewLockFilePath := systemInfo.StewLockFilePath

	lockFile, err := stew.NewLockFile(stewLockFilePath, userOS, userArch)
	stew.CatchAndExit(err)

	if len(lockFile.Packages) == 0 {
		stew.CatchAndExit(stew.NoBinariesInstalledError{})
	}

	if cliFlag {
		for _, pkg := range lockFile.Packages {
			err = stew.DeleteAssetAndBinary(stewPkgPath, stewBinPath, pkg.Asset, pkg.Binary)
			stew.CatchAndExit(err)
		}
		lockFile.Packages = []stew.PackageData{}
	} else {
		var binaryFound bool
		for index, pkg := range lockFile.Packages {
			if pkg.Binary == binaryName {
				err = stew.DeleteAssetAndBinary(stewPkgPath, stewBinPath, pkg.Asset, pkg.Binary)
				stew.CatchAndExit(err)
				lockFile.Packages, err = stew.RemovePackage(lockFile.Packages, index)
				stew.CatchAndExit(err)
				binaryFound = true
				break
			}
		}
		if !binaryFound {
			stew.CatchAndExit(stew.BinaryNotInstalledError{Binary: binaryName})
		}
	}

	err = stew.WriteLockFileJSON(lockFile, stewLockFilePath)
	stew.CatchAndExit(err)
	if cliFlag {
		fmt.Printf("✨ Successfully uninstalled all binaries from %v\n", constants.GreenColor(stewBinPath))
	} else {
		fmt.Printf("✨ Successfully uninstalled the %v binary from %v\n", constants.GreenColor(binaryName), constants.GreenColor(stewBinPath))
	}
}
