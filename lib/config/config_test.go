package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/marwanhawari/stew/lib/config"
)

func TestGetDefaultStewPath(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		userOS string
		want   string
	}{{
		userOS: "darwin",
		want:   filepath.Join(homeDir, ".local", "share", "stew"),
	}, {
		userOS: "windows",
		want:   filepath.Join(homeDir, "AppData", "Local", "stew"),
	}}
	for _, tt := range tests {
		t.Run(tt.userOS, func(t *testing.T) {
			got, err := config.DefaultStewPath(tt.userOS)
			if err != nil {
				t.Fatalf("unexpected error: %+v", err)
			}
			if got != tt.want {
				t.Errorf("DefaultStewPath()\n"+
					" got = %v\n"+
					"want = %v", got, tt.want)
			}
		})
	}
}
