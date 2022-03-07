package stew

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetDefaultStewPath(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Errorf("Could not get os.UserHomeDir()")
	}
	type args struct {
		userOS string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				userOS: "darwin",
			},
			want:    filepath.Join(homeDir, ".local", "share", "stew"),
			wantErr: false,
		},
		{
			name: "test2",
			args: args{
				userOS: "windows",
			},
			want:    filepath.Join(homeDir, "AppData", "Local", "stew"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetDefaultStewPath(tt.args.userOS)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDefaultStewPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetDefaultStewPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
