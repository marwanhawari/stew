package prompt

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/marwanhawari/stew/lib/ui/terminal"
	"github.com/pkg/errors"
)

// ConfirmWarn launches the confirm UI with a warning styling.
func ConfirmWarn(tr terminal.Terminal, message string) (bool, error) {
	if err := ensureTty(tr); err != nil {
		return false, err
	}
	result := false
	prompt := &survey.Confirm{
		Message: message,
	}
	prompt.WithStdio(stdio(tr))
	err := survey.AskOne(prompt, &result, survey.WithIcons(func(icons *survey.IconSet) {
		icons.Question.Text = "!"
		icons.Question.Format = "yellow+hb"
	}))
	if err != nil {
		return false, ExitUserSelectionError{Err: errors.WithStack(err)}
	}

	return result, nil
}
