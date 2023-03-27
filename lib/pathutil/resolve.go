package pathutil

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

// Resolve will resolve the full path for an input path.
func Resolve(filePath string) (string, error) {
	var resolvedPath string
	var err error
	resolvedPath = filePath

	resolvedPath = strings.ReplaceAll(resolvedPath, "\"", "")

	if strings.HasPrefix(resolvedPath, "~") {
		var homeDir string
		homeDir, err = os.UserHomeDir()
		if err != nil {
			return "", errors.WithStack(err)
		}

		resolvedPath = filepath.Join(homeDir, strings.TrimLeft(resolvedPath, "~"))
	}

	resolvedPath = os.ExpandEnv(resolvedPath)

	if !filepath.IsAbs(resolvedPath) {
		resolvedPath, err = filepath.Abs(resolvedPath)
		if err != nil {
			return "", errors.WithStack(err)
		}
	}

	resolvedPath = strings.TrimRight(resolvedPath, "/")

	return resolvedPath, nil
}
