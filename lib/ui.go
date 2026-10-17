package stew

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"golang.org/x/term"
)

// NonInteractive is set by the --yes flag. Then stew does not ask the user: a confirmation is accepted,
// an input keeps its default value, and a selection returns a NonInteractivePromptError.
var NonInteractive bool

// canPrompt checks if stew can ask the user. The survey prompts need a terminal on stdin.
// Without a terminal, survey reads piped input as key presses and can panic.
func canPrompt() bool {
	return !NonInteractive && term.IsTerminal(int(os.Stdin.Fd()))
}

// askOne launches a survey prompt. It returns an error, not a panic, when the prompt fails.
func askOne(message string, prompt survey.Prompt, result interface{}, opts ...survey.AskOpt) (err error) {
	if !canPrompt() {
		return NonInteractivePromptError{Message: message}
	}
	defer func() {
		if r := recover(); r != nil {
			err = ExitUserSelectionError{Err: fmt.Errorf("%v", r)}
		}
	}()
	if err := survey.AskOne(prompt, result, opts...); err != nil {
		return ExitUserSelectionError{Err: err}
	}
	return nil
}

// PromptSelect launches the selection UI
func PromptSelect(message string, options []string) (string, error) {
	result := ""
	prompt := &survey.Select{
		Message: message,
		Options: options,
	}
	err := askOne(message, prompt, &result, survey.WithIcons(func(icons *survey.IconSet) {
		icons.Question.Text = "*"
	}))
	if err != nil {
		return "", err
	}

	return result, nil
}

// PromptMultiSelect launches the multiple selection UI. Without a terminal, it returns the default selections.
func PromptMultiSelect(message string, options []string, defaultSelections []string) ([]string, error) {
	if !canPrompt() {
		return defaultSelections, nil
	}
	result := []string{}
	prompt := &survey.MultiSelect{
		Message: message,
		Options: options,
		Default: defaultSelections,
	}
	err := askOne(message, prompt, &result, survey.WithIcons(func(icons *survey.IconSet) {
		icons.Question.Text = "*"
	}))
	if err != nil {
		return []string{}, err
	}

	return result, nil
}

// PromptInput launches the input UI. Without a terminal, it returns the default input.
func PromptInput(message string, defaultInput string) (string, error) {
	if !canPrompt() {
		return defaultInput, nil
	}
	result := ""
	prompt := &survey.Input{
		Message: message,
		Default: defaultInput,
	}
	err := askOne(message, prompt, &result, survey.WithIcons(func(icons *survey.IconSet) {
		icons.Question.Text = "*"
	}))
	if err != nil {
		return "", err
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
	err := askOne(message, prompt, &result, survey.WithIcons(func(icons *survey.IconSet) {
		icons.Question.Text = "!"
		icons.Question.Format = "yellow+hb"
	}))
	if err != nil {
		return "", err
	}

	return result, nil
}

// WarningPromptConfirm launches the confirm UI with a warning styling. With the --yes flag, it returns true.
func WarningPromptConfirm(message string) (bool, error) {
	if NonInteractive {
		return true, nil
	}
	result := false
	prompt := &survey.Confirm{
		Message: message,
	}
	err := askOne(message, prompt, &result, survey.WithIcons(func(icons *survey.IconSet) {
		icons.Question.Text = "!"
		icons.Question.Format = "yellow+hb"
	}))
	if err != nil {
		return false, err
	}

	return result, nil
}

// warningPromptInput launches the input UI with a warning styling. Without a terminal, it returns the default input.
func warningPromptInput(message string, defaultInput string) (string, error) {
	if !canPrompt() {
		return defaultInput, nil
	}
	result := ""
	prompt := &survey.Input{
		Message: message,
		Default: defaultInput,
	}
	err := askOne(message, prompt, &result, survey.WithIcons(func(icons *survey.IconSet) {
		icons.Question.Text = "!"
		icons.Question.Format = "yellow+hb"
	}))
	if err != nil {
		return "", err
	}

	return result, nil
}
