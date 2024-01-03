package stew

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/huh"
)

func NewKeyMap() *huh.KeyMap {
	keymap := *huh.NewDefaultKeyMap()
	keymap.Input.Next = key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "submit"))
	keymap.Input.Prev = key.NewBinding(key.WithDisabled())
	keymap.Select.Next = key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select"))
	keymap.Select.Prev = key.NewBinding(key.WithDisabled())
	return &keymap
}

// PromptSelect launches the selection UI
func PromptSelect(message string, options []string) (string, error) {
	var result string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title(message).
				Height(10).
				Options(huh.NewOptions(options...)...).
				Value(&result),
		),
	).WithKeyMap(NewKeyMap())

	err := form.Run()
	if err != nil {
		return "", ExitUserSelectionError{Err: err}
	}
	return result, nil
}

// PromptConfirm launches the confirm UI with a warning styling
func PromptConfirm(message string) (bool, error) {
	var result bool

	form := huh.NewConfirm().
		Title(message).
		Value(&result)

	err := form.Run()
	if err != nil {
		return false, ExitUserSelectionError{Err: err}
	}

	return result, nil
}

// PromptInput launches the input UI
func PromptInput(message string, defaultInput string) (string, error) {
	form := huh.NewInput().
		Title(message).
		Value(&defaultInput)
	err := form.Run()
	if err != nil {
		return "", ExitUserSelectionError{Err: err}
	}
	return defaultInput, nil
}
