package stew

import (
	"github.com/AlecAivazis/survey/v2"
)

// PromptSelect launches the selection UI
func PromptSelect(message string, options []string) (string, error) {
	result := ""
	prompt := &survey.Select{
		Message: message,
		Options: options,
	}
	err := survey.AskOne(prompt, &result, survey.WithIcons(func(icons *survey.IconSet) {
		icons.Question.Text = "*"
	}))
	if err != nil {
		return "", ExitUserSelectionError{Err: err}
	}

	return result, nil
}

// WarningPromptSelect launches the selection UI with a warning styling
func WarningPromptSelect(message string, options []string) (string, error) {
	result := ""
	prompt := &survey.Select{
		Message: message,
		Options: options,
	}
	err := survey.AskOne(prompt, &result, survey.WithIcons(func(icons *survey.IconSet) {
		icons.Question.Text = "!"
		icons.Question.Format = "yellow+hb"
	}))
	if err != nil {
		return "", ExitUserSelectionError{Err: err}
	}

	return result, nil
}

// WarningPromptConfirm launches the confirm UI with a warning styling
func WarningPromptConfirm(message string) (bool, error) {
	result := false
	prompt := &survey.Confirm{
		Message: message,
	}
	err := survey.AskOne(prompt, &result, survey.WithIcons(func(icons *survey.IconSet) {
		icons.Question.Text = "!"
		icons.Question.Format = "yellow+hb"
	}))
	if err != nil {
		return false, ExitUserSelectionError{Err: err}
	}

	return result, nil
}

// warningPromptInput launches the input UI with a warning styling
func warningPromptInput(message string, defaultInput string) (string, error) {
	result := ""
	prompt := &survey.Input{
		Message: message,
		Default: defaultInput,
	}
	err := survey.AskOne(prompt, &result, survey.WithIcons(func(icons *survey.IconSet) {
		icons.Question.Text = "!"
		icons.Question.Format = "yellow+hb"
	}))
	if err != nil {
		return "", ExitUserSelectionError{Err: err}
	}

	return result, nil
}
