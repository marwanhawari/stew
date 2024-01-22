package stew

import (
	"os"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/assert"
)

func TestConfigGetDefaultStewPath(t *testing.T) {
	assert := assert.New(t)
	t.Run("Test on non-Windows OS with XDG_DATA_HOME unset", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		patches.ApplyFunc(os.UserHomeDir, func() (string, error) {
			return "/home/testuser", nil
		})
		patches.ApplyFunc(os.Getenv, func(key string) string {
			if key == "XDG_DATA_HOME" {
				return ""
			}
			return os.Getenv(key)
		})

		path, err := GetDefaultStewPath("linux")
		assert.NoError(err)
		expectedPath := "/home/testuser/.local/share/stew"
		assert.Equal(expectedPath, path)
	})

	t.Run("Test on non-Windows OS with XDG_DATA_HOME set", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()

		patches.ApplyFunc(os.UserHomeDir, func() (string, error) {
			return "/home/testuser", nil
		})
		patches.ApplyFunc(os.Getenv, func(key string) string {
			if key == "XDG_DATA_HOME" {
				return "/custom/path"
			}
			return os.Getenv(key)
		})

		path, err := GetDefaultStewPath("linux")
		assert.NoError(err)
		expectedPath := "/custom/path/stew"
		assert.Equal(expectedPath, path)
	})

	// Additional tests for other scenarios...
}
