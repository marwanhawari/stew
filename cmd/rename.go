package cmd

import (
	"os"
	"path/filepath"

	"github.com/marwanhawari/stew/constants"
	stew "github.com/marwanhawari/stew/lib"
	"github.com/marwanhawari/stew/lib/config"
	"github.com/marwanhawari/stew/lib/errs"
)

// Rename is executed when you run `stew rename`
func Rename(cliInput string) {
	rt := errs.Strip(config.Initialize())

	errs.MaybeExit(stew.ValidateCLIInput(cliInput))

	lockFile := errs.Strip(stew.NewLockFile(rt))

	if len(lockFile.Packages) == 0 {
		errs.MaybeExit(stew.NoBinariesInstalledError{})
	}

	var binaryFound bool
	var renamedBinaryName string
	for index, pkg := range lockFile.Packages {
		if pkg.Binary == cliInput {
			renamedBinaryName = errs.Strip(stew.PromptRenameBinary(rt, stew.RenameBinaryArgs{
				Default: cliInput,
			}))
			errs.MaybeExit(os.Rename(filepath.Join(rt.StewBinPath, cliInput),
				filepath.Join(rt.StewBinPath, renamedBinaryName)))

			lockFile.Packages[index].Binary = renamedBinaryName
			binaryFound = true
			break
		}
	}
	if !binaryFound {
		errs.MaybeExit(stew.BinaryNotInstalledError{Binary: cliInput})
	}

	errs.MaybeExit(stew.WriteLockFileJSON(rt, lockFile, rt.LockPath))

	rt.Printf("✨ Successfully renamed the %v binary to %v\n",
		constants.GreenColor(cliInput), constants.GreenColor(renamedBinaryName))
}
