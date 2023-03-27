package cmd

import (
	stewcmd "github.com/marwanhawari/stew/lib/cmd"
	"github.com/marwanhawari/stew/lib/config"
	"github.com/marwanhawari/stew/lib/errs"
)

// Install is executed when you run `stew install`
func Install(cliInputs []string, batchMode bool) {
	rt := errs.Strip(config.Initialize())
	err := stewcmd.Install(rt, cliInputs, batchMode)
	errs.MaybeExit(err)
}
