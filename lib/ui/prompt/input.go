package prompt

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/marwanhawari/stew/lib/ui/terminal"
	"github.com/pkg/errors"
)

// Input launches the input UI.
func Input(tr terminal.Terminal, message string, defaultInput string) (string, error) {
	if err := ensureTty(tr); err != nil {
		return "", err
	}
	result := ""
	prompt := &survey.Input{
		Message: message,
		Default: defaultInput,
	}
	prompt.WithStdio(stdio(tr))
	err := survey.AskOne(prompt, &result, survey.WithIcons(func(icons *survey.IconSet) {
		icons.Question.Text = "*"
	}))
	if err != nil {
		return "", ExitUserSelectionError{Err: err}
	}

	return result, nil
}

// InputWarn launches the input UI with a warning styling.
func InputWarn(tr terminal.Terminal, message string, defaultInput string) (string, error) {
	if err := ensureTty(tr); err != nil {
		return "", err
	}
	result := ""
	prompt := &survey.Input{
		Message: message,
		Default: defaultInput,
	}
	prompt.WithStdio(stdio(tr))
	err := survey.AskOne(prompt, &result, survey.WithIcons(func(icons *survey.IconSet) {
		icons.Question.Text = "!"
		icons.Question.Format = "yellow+hb"
	}))
	if err != nil {
		return "", ExitUserSelectionError{Err: errors.WithStack(err)}
	}

	return result, nil
}
