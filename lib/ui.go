package stew

import (
	"github.com/charmbracelet/huh"
)

// PromptSelect launches the selection UI
func PromptSelect(message string, options []string) (string, error) {
	for i, option := range options {
		max := 128
		if len(option) > max {
			options[i] = option[:max] + "..."
		}
	}
	padding := 2
	height := 10
	if len(options) < height {
		height = len(options)
	}
	result := ""
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().Title(message).Options(
				huh.NewOptions(options...)...,
			).Height(height + padding).Value(&result),
		),
	).WithTheme(huh.ThemeCatppuccin()).Run()
	if err != nil {
		return "", ExitUserSelectionError{Err: err}
	}
	return result, nil
}

// PromptInput launches the input UI
func PromptInput(message string, defaultInput string) (string, error) {
	result := ""
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title(message).
				Prompt("> ").
				Value(&result).
				Placeholder(defaultInput),
		),
	).WithTheme(huh.ThemeCatppuccin()).Run()
	if err != nil {
		return "", ExitUserSelectionError{Err: err}
	}
	if result == "" {
		result = defaultInput
	}

	return result, nil
}

// WarningPromptSelect launches the selection UI with a warning styling
func WarningPromptSelect(message string, options []string) (string, error) {
	for i, option := range options {
		max := 128
		if len(option) > max {
			options[i] = option[:max] + "..."
		}
	}
	padding := 2
	height := 10
	if len(options) < height {
		height = len(options)
	}
	result := ""
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().Title("! " + message).Options(
				huh.NewOptions(options...)...,
			).Height(height + padding).Value(&result),
		),
	).WithTheme(huh.ThemeCatppuccin()).Run()
	if err != nil {
		return "", ExitUserSelectionError{Err: err}
	}
	return result, nil
}

// WarningPromptConfirm launches the confirm UI with a warning styling
func WarningPromptConfirm(message string) (bool, error) {
	var result bool
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("! " + message).
				Affirmative("Yes").
				Value(&result).
				Negative("No"),
		),
	).WithTheme(huh.ThemeCatppuccin()).Run()
	if err != nil {
		return false, ExitUserSelectionError{Err: err}
	}

	return result, nil
}

// PromptInput launches the input UI with a warning styling
func warningPromptInput(message string, defaultInput string) (string, error) {
	result := ""
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("! " + message).
				Prompt("> ").
				Value(&result).
				Placeholder(defaultInput),
		),
	).WithTheme(huh.ThemeCatppuccin()).Run()
	if err != nil {
		return "", ExitUserSelectionError{Err: err}
	}
	if result == "" {
		result = defaultInput
	}

	return result, nil
}
