package terminal

import (
	"fmt"
	"path"
	"runtime"
	"strings"
)

const (
	testLogerSkip = 2
	maxUint       = ^uint(0)
)

// TestingT is part of the testing interface.
type TestingT interface {
	Log(args ...any)
	Logf(format string, args ...any)
	Fatal(args ...any)
}

// TestTerminal is a Terminal that writes to a TestingT.
type TestTerminal struct {
	TestingT
}

func (p TestTerminal) Print(a ...any) {
	p.Log(append([]any{caller(testLogerSkip)}, a...)...)
}

func (p TestTerminal) Printf(format string, a ...any) {
	p.Logf(caller(testLogerSkip)+" "+format, a...)
}

func (p TestTerminal) Println(a ...any) {
	p.Log(append([]any{caller(testLogerSkip)}, a...)...)
}

func (p TestTerminal) Output() FileWriter {
	return testingTWriter{TestingT: p.TestingT}
}

func (p TestTerminal) Input() FileReader {
	return nil
}

var _ Terminal = TestTerminal{}

func caller(skip int) string {
	_, f, l, _ := runtime.Caller(skip)
	_, here, _, _ := runtime.Caller(0)
	base := path.Dir(path.Dir(path.Dir(here)))
	f = strings.Replace(f, base, ".", 1)
	return fmt.Sprintf("%s:%d", f, l)
}

type testingTWriter struct {
	TestingT
}

func (w testingTWriter) Write(bytes []byte) (int, error) {
	w.Log(string(bytes))
	return len(bytes), nil
}

func (w testingTWriter) Fd() uintptr {
	return uintptr(maxUint)
}
