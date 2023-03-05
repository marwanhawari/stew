package testsupport

import (
	"path"

	"github.com/marwanhawari/stew/lib/config"
	"github.com/marwanhawari/stew/lib/ui/terminal"
)

// NewDefaultRuntime creates a default testing stew.Runtime.
func NewDefaultRuntime(t TestingT) config.Runtime {
	return NewRuntime(t, "linux", "amd64")
}

// NewRuntime creates a testing stew.Runtime.
func NewRuntime(t TestingT, os, arch string) config.Runtime {
	tmp := t.TempDir()
	t.Setenv(config.PathEnvVar, path.Join(tmp, "config.json"))
	run := config.Runtime{
		OS:   os,
		Arch: arch,
		Config: config.Config{
			StewPath:    path.Join(tmp, "data"),
			StewBinPath: path.Join(tmp, "bin"),
			GithubAPI:   config.DefaultGithubAPI,
		},
		Terminal: terminal.TestTerminal{TestingT: t},
	}
	if err := run.CreateFiles("linux"); err != nil {
		t.Errorf("unexpected error: %+v", err)
		t.FailNow()
	}
	run.SystemInfo = config.NewSystemInfo(run.Config)

	return run
}

type TestingT interface {
	terminal.TestingT
	Setenv(key, value string)
	TempDir() string
	Errorf(format string, args ...any)
	FailNow()
}
