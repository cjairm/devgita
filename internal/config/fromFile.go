package config

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"time"

	"github.com/cjairm/devgita/pkg/constants"
	"github.com/cjairm/devgita/pkg/files"
	"github.com/cjairm/devgita/pkg/logger"
	"github.com/cjairm/devgita/pkg/paths"
	"gopkg.in/yaml.v3"
)

// Used to store what this app installed
type InstalledConfig struct {
	Packages      []string `yaml:"packages"`
	DesktopApps   []string `yaml:"desktop_apps"`
	Fonts         []string `yaml:"fonts"`
	Themes        []string `yaml:"themes"`
	TerminalTools []string `yaml:"terminal_tools"`
	DevLanguages  []string `yaml:"dev_languages"`
	Databases     []string `yaml:"databases"`
}

// Used to store config that user already had installed before using this app
type AlreadyInstalledConfig struct {
	Packages      []string `yaml:"packages"`
	DesktopApps   []string `yaml:"desktop_apps"`
	Fonts         []string `yaml:"fonts"`
	Themes        []string `yaml:"themes"`
	TerminalTools []string `yaml:"terminal_tools"`
	DevLanguages  []string `yaml:"dev_languages"`
	Databases     []string `yaml:"databases"`
}

// ShellFeatures tracks which shell enhancements are enabled
type ShellFeatures struct {
	IsMac                 bool `yaml:"is_mac"`
	Mise                  bool `yaml:"mise"`
	Zoxide                bool `yaml:"zoxide"`
	ZshAutosuggestions    bool `yaml:"zsh_autosuggestions"`
	ZshSyntaxHighlighting bool `yaml:"zsh_syntax_highlighting"`
	Powerlevel10k         bool `yaml:"powerlevel10k"`
	ExtendedCapabilities  bool `yaml:"extended_capabilities"`
	LazyGit               bool `yaml:"lazy_git"`
	LazyDocker            bool `yaml:"lazy_docker"`
	Fzf                   bool `yaml:"fzf"`
	Neovim                bool `yaml:"neovim"`
	Tmux                  bool `yaml:"tmux"`
	Eza                   bool `yaml:"eza"`
	Bat                   bool `yaml:"bat"`
}

// FailedInstallation tracks packages that failed to install
type FailedInstallation struct {
	PackageName  string    `yaml:"package_name"`
	Category     string    `yaml:"category"` // "package" | "dev_language" | "database"
	ErrorMessage string    `yaml:"error_message"`
	FailedAt     time.Time `yaml:"failed_at"`
	AttemptCount int       `yaml:"attempt_count"`
}

type GlobalConfig struct {
	AppPath             string                 `yaml:"app_path"`
	ConfigPath          string                 `yaml:"config_path"`
	AlreadyInstalled    AlreadyInstalledConfig `yaml:"already_installed"`
	CurrentFont         string                 `yaml:"current_font"`
	CurrentTheme        string                 `yaml:"current_theme"`
	Installed           InstalledConfig        `yaml:"installed"`
	Shortcuts           map[string]string      `yaml:"shortcuts"`
	Shell               ShellFeatures          `yaml:"shell"`
	FailedInstallations []FailedInstallation   `yaml:"failed_installations,omitempty"`
}

func getGlobalConfigFilePath() string {
	return filepath.Join(
		paths.Paths.Config.Root,
		constants.App.Name,
		constants.App.File.GlobalConfig,
	)
}

func (gc *GlobalConfig) Load() error {
	globalConfigFile, err := os.ReadFile(getGlobalConfigFilePath())
	if err != nil {
		return err
	}
	return yaml.Unmarshal(globalConfigFile, gc)
}

func (gc *GlobalConfig) Save() error {
	file, err := yaml.Marshal(gc)
	if err != nil {
		return err
	}
	return os.WriteFile(getGlobalConfigFilePath(), file, 0644)
}

func (gc *GlobalConfig) Reset() error {
	logger.L().Debug("Resetting global config")
	*gc = GlobalConfig{}
	data, err := yaml.Marshal(gc)
	if err != nil {
		return err
	}
	return os.WriteFile(getGlobalConfigFilePath(), data, 0644)
}

func (gc *GlobalConfig) Create() error {
	globalConfigFilePath := getGlobalConfigFilePath()
	if paths.FileAlreadyExist(globalConfigFilePath) {
		return nil
	}
	appFolder := filepath.Join(
		paths.Paths.Config.Root,
		constants.App.Name,
	)
	if !files.DirAlreadyExist(appFolder) {
		if err := os.MkdirAll(appFolder, files.DirPermission); err != nil {
			return err
		}
	}
	// Initialize with empty config structure instead of copying template
	// This avoids dependency on extracted embedded files
	return gc.Reset()
}

func (gc *GlobalConfig) getSliceByType(configType, itemType string) *[]string {
	switch configType {
	case "installed":
		return gc.getInstalledSlice(itemType)
	case "already_installed":
		return gc.getAlreadyInstalledSlice(itemType)
	}
	return nil
}

func (gc *GlobalConfig) getInstalledSlice(itemType string) *[]string {
	switch itemType {
	case "package":
		return &gc.Installed.Packages
	case "desktop_app":
		return &gc.Installed.DesktopApps
	case "font":
		return &gc.Installed.Fonts
	case "theme":
		return &gc.Installed.Themes
	case "terminal_tool":
		return &gc.Installed.TerminalTools
	case "dev_language":
		return &gc.Installed.DevLanguages
	case "database":
		return &gc.Installed.Databases
	}
	return nil
}

func (gc *GlobalConfig) getAlreadyInstalledSlice(itemType string) *[]string {
	switch itemType {
	case "package":
		return &gc.AlreadyInstalled.Packages
	case "desktop_app":
		return &gc.AlreadyInstalled.DesktopApps
	case "font":
		return &gc.AlreadyInstalled.Fonts
	case "theme":
		return &gc.AlreadyInstalled.Themes
	case "terminal_tool":
		return &gc.AlreadyInstalled.TerminalTools
	case "dev_language":
		return &gc.AlreadyInstalled.DevLanguages
	case "database":
		return &gc.AlreadyInstalled.Databases
	}
	return nil
}

func (gc *GlobalConfig) IsTracked(itemName, itemType, configType string) bool {
	slice := gc.getSliceByType(configType, itemType)
	if slice == nil {
		return false
	}
	return slices.Contains(*slice, itemName)
}

func (gc *GlobalConfig) AddToConfig(itemName, itemType, configType string) {
	slice := gc.getSliceByType(configType, itemType)
	if slice == nil {
		return
	}
	if !slices.Contains(*slice, itemName) {
		*slice = append(*slice, itemName)
	}
}

// AddToInstalled adds an item to the installed config
func (gc *GlobalConfig) AddToInstalled(itemName, itemType string) {
	gc.AddToConfig(itemName, itemType, "installed")
}

func (gc *GlobalConfig) AddToAlreadyInstalled(itemName, itemType string) {
	gc.AddToConfig(itemName, itemType, "already_installed")
}

// AddToFailed adds a package to the failed installations list
// It stores the package name, category, error message, timestamp, and attempt count
func (gc *GlobalConfig) AddToFailed(packageName, category, errorMessage string, attemptCount int) {
	// Check if package already in failed list, update if exists
	for i := range gc.FailedInstallations {
		if gc.FailedInstallations[i].PackageName == packageName {
			gc.FailedInstallations[i].ErrorMessage = errorMessage
			gc.FailedInstallations[i].FailedAt = time.Now()
			gc.FailedInstallations[i].AttemptCount = attemptCount
			logger.L().Warnw("Updated failed installation",
				"package", packageName,
				"category", category,
				"error", errorMessage,
			)
			return
		}
	}

	// Add new failed installation
	gc.FailedInstallations = append(gc.FailedInstallations, FailedInstallation{
		PackageName:  packageName,
		Category:     category,
		ErrorMessage: errorMessage,
		FailedAt:     time.Now(),
		AttemptCount: attemptCount,
	})
	logger.L().Warnw("Added to failed installations",
		"package", packageName,
		"category", category,
		"error", errorMessage,
	)
}

func (gc *GlobalConfig) IsInstalledByDevgita(itemName, itemType string) bool {
	return gc.IsTracked(itemName, itemType, "installed")
}

func (gc *GlobalConfig) IsAlreadyInstalled(itemName, itemType string) bool {
	return gc.IsTracked(itemName, itemType, "already_installed")
}

// EnableShellFeature enables a shell feature by name
func (gc *GlobalConfig) EnableShellFeature(featureName string) {
	switch featureName {
	case constants.Mise:
		gc.Shell.Mise = true
	case constants.Zoxide:
		gc.Shell.Zoxide = true
	case constants.ZshAutosuggestions:
		gc.Shell.ZshAutosuggestions = true
	case constants.Syntaxhighlighting:
		gc.Shell.ZshSyntaxHighlighting = true
	case constants.Powerlevel10k:
		gc.Shell.Powerlevel10k = true
	case "extended_capabilities":
		gc.Shell.ExtendedCapabilities = true
	case constants.LazyGit:
		gc.Shell.LazyGit = true
	case constants.LazyDocker:
		gc.Shell.LazyDocker = true
	case constants.Fzf:
		gc.Shell.Fzf = true
	case constants.Neovim:
		gc.Shell.Neovim = true
	case constants.Tmux:
		gc.Shell.Tmux = true
	case constants.Eza:
		gc.Shell.Eza = true
	case constants.Bat:
		gc.Shell.Bat = true
	}
}

// DisableShellFeature disables a shell feature by name
func (gc *GlobalConfig) DisableShellFeature(featureName string) {
	switch featureName {
	case constants.Mise:
		gc.Shell.Mise = false
	case constants.Zoxide:
		gc.Shell.Zoxide = false
	case constants.ZshAutosuggestions:
		gc.Shell.ZshAutosuggestions = false
	case constants.Syntaxhighlighting:
		gc.Shell.ZshSyntaxHighlighting = false
	case constants.Powerlevel10k:
		gc.Shell.Powerlevel10k = false
	case "extended_capabilities":
		gc.Shell.ExtendedCapabilities = false
	case constants.LazyGit:
		gc.Shell.LazyGit = false
	case constants.LazyDocker:
		gc.Shell.LazyDocker = false
	case constants.Fzf:
		gc.Shell.Fzf = false
	case constants.Neovim:
		gc.Shell.Neovim = false
	case constants.Tmux:
		gc.Shell.Tmux = false
	case constants.Eza:
		gc.Shell.Eza = false
	case constants.Bat:
		gc.Shell.Bat = false
	}
}

// IsShellFeatureEnabled checks if a shell feature is enabled
func (gc *GlobalConfig) IsShellFeatureEnabled(featureName string) bool {
	switch featureName {
	case constants.Mise:
		return gc.Shell.Mise
	case constants.Zoxide:
		return gc.Shell.Zoxide
	case constants.ZshAutosuggestions:
		return gc.Shell.ZshAutosuggestions
	case constants.Syntaxhighlighting:
		return gc.Shell.ZshSyntaxHighlighting
	case constants.Powerlevel10k:
		return gc.Shell.Powerlevel10k
	case "extended_capabilities":
		return gc.Shell.ExtendedCapabilities
	case constants.LazyGit:
		return gc.Shell.LazyGit
	case constants.LazyDocker:
		return gc.Shell.LazyDocker
	case constants.Fzf:
		return gc.Shell.Fzf
	case constants.Neovim:
		return gc.Shell.Neovim
	case constants.Tmux:
		return gc.Shell.Tmux
	case constants.Eza:
		return gc.Shell.Eza
	case constants.Bat:
		return gc.Shell.Bat
	}
	return false
}

func (gc *GlobalConfig) RegenerateShellConfig() error {
	templatePath := filepath.Join(
		paths.Paths.App.Configs.Templates,
		constants.App.Template.ShellConfig,
	)
	outputPath := filepath.Join(paths.Paths.App.Root, fmt.Sprintf("%s.zsh", constants.App.Name))
	return files.GenerateFromTemplate(templatePath, outputPath, gc.Shell)
}
