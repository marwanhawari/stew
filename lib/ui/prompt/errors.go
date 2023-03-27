package prompt

import (
	"fmt"

	"github.com/marwanhawari/stew/constants"
)

// ExitUserSelectionError occurs when exiting from the terminal UI.
type ExitUserSelectionError struct {
	Err error
}

func (e ExitUserSelectionError) Error() string {
	return fmt.Sprintf("%v Exited from user selection: %v",
		constants.RedColor("Error:"), constants.RedColor(fmt.Sprintf("%+v", e.Err)))
}
