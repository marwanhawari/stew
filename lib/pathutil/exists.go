package pathutil

import (
	"os"

	"github.com/pkg/errors"
)

// Exists checks if a given path exists.
func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, errors.WithStack(err)
	}

	return true, nil
}
