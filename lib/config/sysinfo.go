package config

import "path/filepath"

// SystemInfo contains system specific info like OS, arch, and ~/.stew paths.
type SystemInfo struct {
	PkgPath  string
	LockPath string
	TmpPath  string
}

// NewSystemInfo creates a new instance of the SystemInfo struct.
func NewSystemInfo(stewConfig Config) SystemInfo {
	var systemInfo SystemInfo
	systemInfo.PkgPath = filepath.Join(stewConfig.StewPath, "pkg")
	systemInfo.LockPath = filepath.Join(stewConfig.StewPath, "Stewfile.lock.json")
	systemInfo.TmpPath = filepath.Join(stewConfig.StewPath, "tmp")
	return systemInfo
}
