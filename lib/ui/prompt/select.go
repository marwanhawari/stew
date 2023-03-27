package prompt

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/marwanhawari/stew/lib/ui/terminal"
	"github.com/pkg/errors"
)

// Select launches the selection UI.
func Select(tr terminal.Terminal, message string, options []string) (string, error) {
	if err := ensureTty(tr); err != nil {
		return "", err
	}
	result := ""
	prompt := &survey.Select{
		Message: message,
		Options: options,
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

// SelectWarn launches the selection UI with a warning styling.
func SelectWarn(tr terminal.Terminal, message string, options []string) (string, error) {
	if err := ensureTty(tr); err != nil {
		return "", err
	}
	result := ""
	prompt := &survey.Select{
		Message: message,
		Options: options,
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
