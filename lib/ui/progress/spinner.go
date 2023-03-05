package progress

import (
	"time"

	"github.com/briandowns/spinner"
	"github.com/marwanhawari/stew/lib/ui/terminal"
)

// Spinner returns a new spinner.
func Spinner(io terminal.Terminal) *spinner.Spinner {
	return spinner.New(
		spinner.CharSets[9],
		100*time.Millisecond,
		spinner.WithColor("cyan"),
		spinner.WithHiddenCursor(true),
		spinner.WithWriter(io.Output()),
	)
}
