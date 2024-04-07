package stew

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/marwanhawari/stew/constants"
)

// GetDefaultStewPath will return the default path to the top-level stew directory
func GetDefaultStewPath(userOS string) (string, error) {

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	var stewPath string
	switch userOS {
	case "windows":
		stewPath = filepath.Join(homeDir, "AppData", "Local", "stew")
	default:
		xdgDataHomePath := os.Getenv("XDG_DATA_HOME")
		if xdgDataHomePath == "" {
			stewPath = filepath.Join(homeDir, ".local", "share", "stew")
		} else {
			stewPath = filepath.Join(xdgDataHomePath, "stew")

		}
	}

	return stewPath, nil
}

// GetDefaultStewBinPath will return the default path where binaries are installed by stew
func GetDefaultStewBinPath(userOS string) (string, error) {

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	var stewBinPath string
	switch userOS {
	case "windows":
		stewBinPath = filepath.Join(homeDir, "AppData", "Local", "stew", "bin")
	default:
		stewBinPath = filepath.Join(homeDir, ".local", "bin")

	}

	return stewBinPath, nil
}

// GetStewConfigFilePath will return the stew config file path
func GetStewConfigFilePath(userOS string) (string, error) {

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	var stewConfigFilePath string
	switch userOS {
	case "windows":
		stewConfigFilePath = filepath.Join(homeDir, "AppData", "Local", "stew", "Config", "stew.config.json")
	default:
		xdgConfigHomePath := os.Getenv("XDG_CONFIG_HOME")
		if xdgConfigHomePath == "" {
			stewConfigFilePath = filepath.Join(homeDir, ".config", "stew", "stew.config.json")
		} else {
			stewConfigFilePath = filepath.Join(xdgConfigHomePath, "stew", "stew.config.json")
		}
	}

	return stewConfigFilePath, nil
}

// StewConfig contains all the stew configuration data
type StewConfig struct {
	StewPath    string `json:"stewPath"`
	StewBinPath string `json:"stewBinPath"`
}

func ReadStewConfigJSON(stewConfigFilePath string) (StewConfig, error) {

	stewConfigFileBytes, err := os.ReadFile(stewConfigFilePath)
	if err != nil {
		return StewConfig{}, err
	}

	var stewConfig StewConfig
	err = json.Unmarshal(stewConfigFileBytes, &stewConfig)
	if err != nil {
		return StewConfig{}, err
	}

	return stewConfig, nil
}

// WriteStewConfigJSON will write the config JSON file
func WriteStewConfigJSON(stewConfigFileJSON StewConfig, outputPath string) error {

	stewConfigFileBytes, err := json.MarshalIndent(stewConfigFileJSON, "", "\t")
	if err != nil {
		return err
	}

	err = os.WriteFile(outputPath, stewConfigFileBytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

// NewStewConfig creates a new instance of the StewConfig struct
func NewStewConfig(userOS string) (StewConfig, error) {
	var stewConfig StewConfig

	stewConfigFilePath, err := GetStewConfigFilePath(userOS)
	if err != nil {
		return StewConfig{}, err
	}
	defaultStewPath, err := GetDefaultStewPath(userOS)
	if err != nil {
		return StewConfig{}, err
	}
	defaultStewBinPath, err := GetDefaultStewBinPath(userOS)
	if err != nil {
		return StewConfig{}, err
	}

	configExists, err := PathExists(stewConfigFilePath)
	if err != nil {
		return StewConfig{}, err
	}
	if configExists {
		stewConfig, err = ReadStewConfigJSON(stewConfigFilePath)
		if err != nil {
			return StewConfig{}, err
		}

		if stewConfig.StewPath == "" {
			stewConfig.StewPath = defaultStewPath
		}

		if stewConfig.StewBinPath == "" {
			stewConfig.StewBinPath = defaultStewBinPath
		}
	} else {
		selectedStewPath, selectedStewBinPath, err := PromptConfig(defaultStewPath, defaultStewBinPath)
		if err != nil {
			return StewConfig{}, err
		}
		stewConfig.StewPath = selectedStewPath
		stewConfig.StewBinPath = selectedStewBinPath
		fmt.Printf("📄 Updated %v\n", constants.GreenColor(stewConfigFilePath))
	}

	pathVariable := os.Getenv("PATH")
	ValidateStewBinPath(stewConfig.StewBinPath, pathVariable)

	err = createStewDirsAndFiles(stewConfig, stewConfigFilePath)
	if err != nil {
		return StewConfig{}, err
	}

	return stewConfig, nil
}

func createStewDirsAndFiles(stewConfig StewConfig, stewConfigFilePath string) error {
	var err error

	err = os.MkdirAll(stewConfig.StewPath, 0755)
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Join(stewConfig.StewPath, "pkg"), 0755)
	if err != nil {
		return err
	}

	err = os.MkdirAll(stewConfig.StewBinPath, 0755)
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(stewConfigFilePath), 0755)
	if err != nil {
		return err
	}
	err = WriteStewConfigJSON(stewConfig, stewConfigFilePath)
	if err != nil {
		return err
	}

	return nil
}

// SystemInfo contains system specific info like OS, arch, and ~/.stew paths
type SystemInfo struct {
	StewPath         string
	StewBinPath      string
	StewPkgPath      string
	StewLockFilePath string
	StewTmpPath      string
}

// NewSystemInfo creates a new instance of the SystemInfo struct
func NewSystemInfo(stewConfig StewConfig) SystemInfo {
	var systemInfo SystemInfo
	systemInfo.StewPath = stewConfig.StewPath
	systemInfo.StewBinPath = stewConfig.StewBinPath
	systemInfo.StewPkgPath = filepath.Join(stewConfig.StewPath, "pkg")
	systemInfo.StewLockFilePath = filepath.Join(stewConfig.StewPath, "Stewfile.lock.json")
	systemInfo.StewTmpPath = filepath.Join(stewConfig.StewPath, "tmp")
	return systemInfo
}

// Initialize returns pertinent initialization information like OS, arch, configuration, and system info
func Initialize() (string, string, StewConfig, SystemInfo, error) {
	userOS := runtime.GOOS
	userArch := runtime.GOARCH
	stewConfig, err := NewStewConfig(userOS)
	if err != nil {
		return "", "", StewConfig{}, SystemInfo{}, err
	}
	systemInfo := NewSystemInfo(stewConfig)

	return userOS, userArch, stewConfig, systemInfo, nil
}

// PromptConfig launches an interactive UI for setting the stew config values. It returns the resolved stewPath and stewBinPath.
func PromptConfig(suggestedStewPath, suggestedStewBinPath string) (string, string, error) {
	inputStewPath, err := PromptInput("Set the stewPath. This will contain all stew data other than the binaries.", suggestedStewPath)
	if err != nil {
		return "", "", err
	}
	inputStewBinPath, err := PromptInput("Set the stewBinPath. This is where the binaries will be installed by stew.", suggestedStewBinPath)
	if err != nil {
		return "", "", err
	}

	fullStewPath, err := ResolvePath(inputStewPath)
	if err != nil {
		return "", "", err
	}
	fullStewBinPath, err := ResolvePath(inputStewBinPath)
	if err != nil {
		return "", "", err
	}

	return fullStewPath, fullStewBinPath, nil
}

func ValidateStewBinPath(stewBinPath, pathVariable string) bool {
	if !strings.Contains(pathVariable, stewBinPath) {
		fmt.Printf("%v The stewBinPath %v is not in your PATH variable.\nYou need to add %v to PATH.\n", constants.YellowColor("WARNING:"), constants.YellowColor(stewBinPath), constants.YellowColor(stewBinPath))
		fmt.Printf("Add the following line to your ~/.zshrc or ~/.bashrc file then start a new terminal session:\n\nexport PATH=\"%v:$PATH\"\n\n", stewBinPath)
		return false
	}

	return true
}
