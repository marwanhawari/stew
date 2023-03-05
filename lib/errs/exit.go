package errs

import (
	"fmt"
	"os"

	"github.com/wavesoftware/go-retcode"
)

var exitFn = os.Exit

// MaybeExit will catch errors and immediately exit, if error is not nil.
func MaybeExit(err error) {
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%+v\n", err)
		exitFn(retcode.Calc(err))
	}
}

// WithExitFn will execute given function with given exit handler. This is
// useful for testing.
func WithExitFn(hndl func(int), fn func()) {
	original := exitFn
	exitFn = hndl
	defer func() {
		exitFn = original
	}()
	fn()
}
