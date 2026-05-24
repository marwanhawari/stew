package cmd

import (
	"io"
	"os"
	"strings"
	"testing"
)

func TestGetZshCompletion(t *testing.T) {
	completion := GetZshCompletion()

	for _, snippet := range []string{
		"compdef _stew stew",
		"'completion:Generate shell completion scripts'",
	} {
		if !strings.Contains(completion, snippet) {
			t.Errorf("GetZshCompletion() missing %q", snippet)
		}
	}
}

func TestRunCompletion(t *testing.T) {
	t.Run("zsh", func(t *testing.T) {
		oldStdout := os.Stdout
		r, w, err := os.Pipe()
		if err != nil {
			t.Fatalf("os.Pipe() error = %v", err)
		}
		os.Stdout = w
		t.Cleanup(func() {
			os.Stdout = oldStdout
		})

		err = RunCompletion("zsh")
		w.Close()
		if err != nil {
			t.Fatalf("RunCompletion() error = %v", err)
		}

		output, readErr := io.ReadAll(r)
		if readErr != nil {
			t.Fatalf("io.ReadAll() error = %v", readErr)
		}
		if !strings.Contains(string(output), "compdef _stew stew") {
			t.Errorf("RunCompletion() output missing compdef registration")
		}
	})

	t.Run("unsupported shell", func(t *testing.T) {
		err := RunCompletion("bash")
		if err == nil {
			t.Fatal("RunCompletion() error = nil, want unsupported shell error")
		}
		if !strings.Contains(err.Error(), "unsupported shell: bash") {
			t.Errorf("RunCompletion() error = %v, want unsupported shell message", err)
		}
	})
}
