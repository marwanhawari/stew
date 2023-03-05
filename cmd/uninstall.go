package cmd

import (
	"github.com/marwanhawari/stew/constants"
	stew "github.com/marwanhawari/stew/lib"
	"github.com/marwanhawari/stew/lib/config"
	"github.com/marwanhawari/stew/lib/errs"
)

// Uninstall is executed when you run `stew uninstall`
func Uninstall(cliFlag bool, binaryName string) {
	rt := errs.Strip(config.Initialize())

	if cliFlag && binaryName != "" {
		errs.MaybeExit(stew.CLIFlagAndInputError{})
	} else if !cliFlag {
		errs.MaybeExit(stew.ValidateCLIInput(binaryName))
	}

	lockFile := errs.Strip(stew.NewLockFile(rt))

	if len(lockFile.Packages) == 0 {
		errs.MaybeExit(stew.NoBinariesInstalledError{})
	}

	if cliFlag {
		for _, pkg := range lockFile.Packages {
			errs.MaybeExit(stew.DeleteAssetAndBinary(rt.PkgPath, rt.StewBinPath, pkg.Asset, pkg.Binary))
		}
		lockFile.Packages = []stew.PackageData{}
	} else {
		var binaryFound bool
		for index, pkg := range lockFile.Packages {
			if pkg.Binary == binaryName {
				errs.MaybeExit(stew.DeleteAssetAndBinary(
					rt.PkgPath, rt.StewBinPath, pkg.Asset, pkg.Binary))
				lockFile.Packages = errs.Strip(stew.RemovePackage(lockFile.Packages, index))
				binaryFound = true
				break
			}
		}
		if !binaryFound {
			errs.MaybeExit(stew.BinaryNotInstalledError{Binary: binaryName})
		}
	}

	errs.MaybeExit(stew.WriteLockFileJSON(rt, lockFile, rt.LockPath))
	if cliFlag {
		rt.Printf("✨ Successfully uninstalled all binaries from %v\n",
			constants.GreenColor(rt.StewBinPath))
	} else {
		rt.Printf("✨ Successfully uninstalled the %v binary from %v\n",
			constants.GreenColor(binaryName), constants.GreenColor(rt.StewBinPath))
	}
}
