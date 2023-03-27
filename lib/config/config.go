package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/marwanhawari/stew/constants"
	"github.com/marwanhawari/stew/lib/pathutil"
	"github.com/marwanhawari/stew/lib/ui/prompt"
	"github.com/marwanhawari/stew/lib/ui/terminal"
	"github.com/pkg/errors"
)

const (
	// PathEnvVar is the environment variable that can be used to specify the
	// path to the stew config file.
	PathEnvVar = "STEW_CONFIG_PATH"
	// DefaultGithubAPI is the default GitHub API endpoint.
	DefaultGithubAPI = "api.github.com"

	osWindows = "windows"
)

// DefaultStewPath will return the default path to the top-level stew directory.
func DefaultStewPath(userOS string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", errors.WithStack(err)
	}

	var stewPath string
	switch userOS {
	case osWindows:
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

// DefaultBinPath will return the default path where binaries are installed by stew.
func DefaultBinPath(userOS string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", errors.WithStack(err)
	}

	var stewBinPath string
	switch userOS {
	case osWindows:
		stewBinPath = filepath.Join(homeDir, "AppData", "Local", "stew", "bin")
	default:
		stewBinPath = filepath.Join(homeDir, ".local", "bin")
	}

	return stewBinPath, nil
}

// FilePath will return the stew config file path.
func FilePath(userOS string) (string, error) {
	if cp := os.Getenv(PathEnvVar); cp != "" {
		return cp, nil
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", errors.WithStack(err)
	}

	var stewConfigFilePath string
	switch userOS {
	case osWindows:
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

// Config contains all the stew configuration data.
type Config struct {
	StewPath    string `json:"stewPath,omitempty"`
	StewBinPath string `json:"stewBinPath,omitempty"`
	GithubAPI   string `json:"githubAPI,omitempty"`
	GithubToken string `json:"githubToken,omitempty"`
}

func readStewConfigJSON(stewConfigFilePath string, config *Config) error {
	stewConfigFileBytes, err := os.ReadFile(stewConfigFilePath)
	if err != nil {
		return errors.WithStack(err)
	}

	err = json.Unmarshal(stewConfigFileBytes, config)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// WriteStewConfigJSON will write the config JSON file.
func (c Config) WriteStewConfigJSON(outputPath string) error {
	stewConfigFileBytes, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return errors.WithStack(err)
	}

	err = os.WriteFile(outputPath, stewConfigFileBytes, 0o644)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// NewConfig creates a new instance of the Config struct.
func NewConfig(prt terminal.Terminal, userOS string) (Config, error) {
	stewConfig := Config{}

	stewConfigFilePath, err := FilePath(userOS)
	if err != nil {
		return Config{}, err
	}
	stewPath, err := DefaultStewPath(userOS)
	if err != nil {
		return Config{}, err
	}
	binPath, err := DefaultBinPath(userOS)
	if err != nil {
		return Config{}, err
	}

	configExists, err := pathutil.Exists(stewConfigFilePath)
	if err != nil {
		return Config{}, err
	}
	if configExists {
		err = readStewConfigJSON(stewConfigFilePath, &stewConfig)
		if err != nil {
			return Config{}, err
		}

		if stewConfig.StewPath == "" {
			stewConfig.StewPath = stewPath
		}

		if stewConfig.StewBinPath == "" {
			stewConfig.StewBinPath = binPath
		}
	} else {
		selectedStewPath, selectedBinPath, err := PromptConfig(prt, stewPath, binPath)
		if err != nil {
			return Config{}, err
		}
		stewConfig.StewPath = selectedStewPath
		stewConfig.StewBinPath = selectedBinPath
		prt.Printf("📄 Updated %v\n", constants.GreenColor(stewConfigFilePath))
	}

	pathVariable := os.Getenv("PATH")
	ValidateBinPath(prt, stewConfig.StewBinPath, pathVariable)

	err = stewConfig.CreateFiles(userOS)
	if err != nil {
		return Config{}, err
	}

	if stewConfig.GithubAPI == "" {
		stewConfig.GithubAPI = DefaultGithubAPI
	}

	return stewConfig, nil
}

// CreateFiles will create the stew config files and directories.
func (c Config) CreateFiles(userOS string) error {
	stewConfigFilePath, err := FilePath(userOS)
	if err != nil {
		return err
	}

	err = os.MkdirAll(c.StewPath, 0o755)
	if err != nil {
		return errors.WithStack(err)
	}
	err = os.MkdirAll(filepath.Join(c.StewPath, "pkg"), 0o755)
	if err != nil {
		return errors.WithStack(err)
	}

	err = os.MkdirAll(c.StewBinPath, 0o755)
	if err != nil {
		return errors.WithStack(err)
	}

	err = os.MkdirAll(filepath.Dir(stewConfigFilePath), 0o755)
	if err != nil {
		return errors.WithStack(err)
	}
	err = c.WriteStewConfigJSON(stewConfigFilePath)
	if err != nil {
		return err
	}

	return nil
}

// PromptConfig launches an interactive UI for setting the stew config values.
// It returns the resolved stewPath and stewBinPath.
func PromptConfig(io terminal.Terminal, suggestedStewPath, suggestedStewBinPath string) (string, string, error) {
	inputStewPath, err := prompt.Input(io,
		"Set the stewPath. This will contain all stew data other than the binaries.",
		suggestedStewPath)
	if err != nil {
		return "", "", err
	}
	inputStewBinPath, err := prompt.Input(io,
		"Set the stewBinPath. This is where the binaries will be installed by stew.",
		suggestedStewBinPath)
	if err != nil {
		return "", "", err
	}

	fullStewPath, err := pathutil.Resolve(inputStewPath)
	if err != nil {
		return "", "", err
	}
	fullStewBinPath, err := pathutil.Resolve(inputStewBinPath)
	if err != nil {
		return "", "", err
	}

	return fullStewPath, fullStewBinPath, nil
}

func ValidateBinPath(prt terminal.Terminal, stewBinPath, pathVariable string) bool {
	if !strings.Contains(pathVariable, stewBinPath) {
		prt.Printf("%v The stewBinPath %v is not in your PATH variable.\nYou need to add %v to PATH.\n", constants.YellowColor("WARNING:"), constants.YellowColor(stewBinPath), constants.YellowColor(stewBinPath))
		prt.Printf("Add the following line to your ~/.zshrc or ~/.bashrc file then start a new terminal session:\n\nexport PATH=\"%v:$PATH\"\n\n", stewBinPath)
		return false
	}

	return true
}
