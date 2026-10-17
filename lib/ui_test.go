package stew

import (
	"reflect"
	"testing"
)

// The tests set NonInteractive, thus no prompt waits for input, also when the tests run in a terminal.
func TestPrompts_NonInteractive(t *testing.T) {
	NonInteractive = true
	defer func() { NonInteractive = false }()

	if got, err := WarningPromptConfirm("Overwrite?"); err != nil || !got {
		t.Errorf("WarningPromptConfirm() = %v, %v, want true, nil", got, err)
	}
	if got, err := PromptInput("Set the stewPath.", "/default/path"); err != nil || got != "/default/path" {
		t.Errorf("PromptInput() = %v, %v, want /default/path, nil", got, err)
	}
	if got, err := PromptRenameBinary("fzf"); err != nil || got != "fzf" {
		t.Errorf("PromptRenameBinary() = %v, %v, want fzf, nil", got, err)
	}
	if got, err := PromptMultiSelect("Exclude?", []string{"a", "b"}, []string{"b"}); err != nil || !reflect.DeepEqual(got, []string{"b"}) {
		t.Errorf("PromptMultiSelect() = %v, %v, want [b], nil", got, err)
	}

	wantErr := NonInteractivePromptError{Message: "Select an asset:"}
	if got, err := WarningPromptSelect("Select an asset:", []string{"a", "b"}); err != wantErr || got != "" {
		t.Errorf("WarningPromptSelect() = %v, %v, want \"\", %v", got, err, wantErr)
	}
	if got, err := PromptSelect("Select an asset:", []string{"a", "b"}); err != wantErr || got != "" {
		t.Errorf("PromptSelect() = %v, %v, want \"\", %v", got, err, wantErr)
	}
}

// Without a terminal on stdin, a confirmation is not accepted automatically. Only the --yes flag accepts it.
func TestWarningPromptConfirm_NoTerminal(t *testing.T) {
	if canPrompt() {
		t.Skip("stdin is a terminal")
	}
	wantErr := NonInteractivePromptError{Message: "Overwrite?"}
	if got, err := WarningPromptConfirm("Overwrite?"); err != wantErr || got {
		t.Errorf("WarningPromptConfirm() = %v, %v, want false, %v", got, err, wantErr)
	}
}
