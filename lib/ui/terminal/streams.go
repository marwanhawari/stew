package terminal

import (
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/pkg/errors"
)

// InOut is a Terminal that writes to an io.Writer.
type InOut struct {
	Out FileWriter
	In  FileReader
}

// Standard returns a Terminal that writes to color.Output.
func Standard() InOut {
	out, ok := color.Output.(FileWriter)
	if !ok {
		// unreachable
		log.Fatal("can't use color.Output as FileWriter")
	}
	return InOut{
		Out: out,
		In:  os.Stdin,
	}
}

func (w InOut) Print(a ...any) {
	if _, err := fmt.Fprint(w.Out, a...); err != nil {
		w.fatal(err)
	}
}

func (w InOut) Printf(format string, a ...any) {
	if _, err := fmt.Fprintf(w.Out, format, a...); err != nil {
		w.fatal(err)
	}
}

func (w InOut) Println(a ...any) {
	if _, err := fmt.Fprintln(w.Out, a...); err != nil {
		w.fatal(err)
	}
}

func (w InOut) Output() FileWriter {
	return w.Out
}

func (w InOut) Input() FileReader {
	return w.In
}

var _ Terminal = InOut{}

func (w InOut) fatal(err error) {
	log.Fatalf("%+v\n", errors.WithStack(err))
}
