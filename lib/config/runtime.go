package config

import (
	"runtime"

	"github.com/marwanhawari/stew/lib/ui/terminal"
)

// Runtime contains all the runtime information for Stew.
type Runtime struct {
	OS   string
	Arch string
	Config
	SystemInfo
	terminal.Terminal
}

// Initialize returns pertinent initialization information like OS, arch, configuration, and system info.
func Initialize() (Runtime, error) {
	tr := terminal.Standard()
	config, err := NewConfig(tr, runtime.GOOS)
	if err != nil {
		return Runtime{}, err
	}

	return NewRuntime(tr, config), nil
}

// NewRuntime creates a new instance of the Runtime struct.
func NewRuntime(tr terminal.Terminal, config Config) Runtime {
	systemInfo := NewSystemInfo(config)

	return Runtime{
		runtime.GOOS, runtime.GOARCH,
		config, systemInfo, tr,
	}
}
