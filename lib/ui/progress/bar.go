package progress

import (
	"fmt"
	"time"

	"github.com/marwanhawari/stew/lib/ui/terminal"
	"github.com/schollz/progressbar/v3"
)

// Bar returns a new progress bar.
func Bar(io terminal.Terminal, size int64, description string) *progressbar.ProgressBar {
	writer := io.Output()
	b := progressbar.NewOptions64(
		size,
		progressbar.OptionSetDescription(description),
		progressbar.OptionSetWriter(writer),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(10),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionShowCount(),
		progressbar.OptionOnCompletion(func() {
			_, _ = fmt.Fprintln(writer)
		}),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionFullWidth(),
	)
	_ = b.RenderBlank()
	return b
}
