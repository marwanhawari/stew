package prompt

import (
	"errors"

	"github.com/marwanhawari/stew/lib/ui/terminal"
	"github.com/mattn/go-isatty"
)

// ErrStdinNotTerminal is returned when stdin is not a terminal.
var ErrStdinNotTerminal = errors.New("stdin is not interactive terminal")

func ensureTty(io terminal.Terminal) error {
	i := io.Input()
	err := ExitUserSelectionError{
		Err: ErrStdinNotTerminal,
	}
	if i == nil {
		return err
	}
	if fi, ok := i.(filelike); !ok || !isatty.IsTerminal(fi.Fd()) {
		return err
	}
	return nil
}

type filelike interface {
	Fd() uintptr
}
