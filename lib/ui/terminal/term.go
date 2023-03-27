package terminal

import "io"

// Terminal can be used to print and read from the terminal.
type Terminal interface {
	Print(a ...any)
	Printf(format string, a ...any)
	Println(a ...any)
	Output() FileWriter
	Input() FileReader
}

// FileWriter provides a minimal interface for Stdin.
type FileWriter interface {
	io.Writer
	Fd() uintptr
}

// FileReader provides a minimal interface for Stdout.
type FileReader interface {
	io.Reader
	Fd() uintptr
}
