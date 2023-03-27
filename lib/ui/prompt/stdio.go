package prompt

import (
	"github.com/AlecAivazis/survey/v2/terminal"
	term "github.com/marwanhawari/stew/lib/ui/terminal"
)

func stdio(tr term.Terminal) terminal.Stdio {
	return terminal.Stdio{
		In:  tr.Input(),
		Out: tr.Output(),
		Err: tr.Output(),
	}
}
