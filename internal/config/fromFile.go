package config

import (
	"os"
	"path/filepath"

	"github.com/cjairm/devgita/pkg/constants"
	"github.com/cjairm/devgita/pkg/files"
	"github.com/cjairm/devgita/pkg/paths"
	"gopkg.in/yaml.v3"
)

type IgnoredConfig struct {
	Packages      []string `yaml:"packages"`
	DesktopApps   []string `yaml:"desktop_apps"`
	Fonts         []string `yaml:"fonts"`
	Themes        []string `yaml:"themes"`
	TerminalTools []string `yaml:"terminal_tools"`
	DevLanguages  []string `yaml:"dev_languages"`
	Databases     []string `yaml:"databases"`
}

type InstalledConfig struct {
	Packages      []string `yaml:"packages"`
	DesktopApps   []string `yaml:"desktop_apps"`
	Fonts         []string `yaml:"fonts"`
	Themes        []string `yaml:"themes"`
	TerminalTools []string `yaml:"terminal_tools"`
	DevLanguages  []string `yaml:"dev_languages"`
	Databases     []string `yaml:"databases"`
}

type AlreadyInstalledConfig struct {
	Packages      []string `yaml:"packages"`
	DesktopApps   []string `yaml:"desktop_apps"`
	Fonts         []string `yaml:"fonts"`
	Themes        []string `yaml:"themes"`
	TerminalTools []string `yaml:"terminal_tools"`
	DevLanguages  []string `yaml:"dev_languages"`
	Databases     []string `yaml:"databases"`
}

type GlobalConfig struct {
	AppPath                string                 `yaml:"app_path"`
	ConfigPath             string                 `yaml:"config_path"`
	AlreadyInstalledConfig AlreadyInstalledConfig `yaml:"already_installed"`
	CurrentFont            string                 `yaml:"current_font"`
	CurrentTheme           string                 `yaml:"current_theme"`
	Ignored                IgnoredConfig          `yaml:"ignored"`
	Installed              InstalledConfig        `yaml:"installed"`
	Shortcuts              map[string]string      `yaml:"shortcuts"`
}

var globalConfigFilePath = filepath.Join(
	paths.ConfigDir,
	constants.AppName,
	constants.GlobalConfigFile,
)

func LoadGlobalConfig() (*GlobalConfig, error) {
	globalConfigFile, err := os.ReadFile(globalConfigFilePath)
	if err != nil {
		return nil, err
	}
	var globalConfig GlobalConfig
	err = yaml.Unmarshal(globalConfigFile, &globalConfig)
	if err != nil {
		return nil, err
	}
	return &globalConfig, nil
}

func SetGlobalConfig(globalConfig *GlobalConfig) error {
	file, err := yaml.Marshal(globalConfig)
	if err != nil {
		return err
	}
	return os.WriteFile(globalConfigFilePath, file, 0644)
}

func ResetGlobalConfig() error {
	data, err := yaml.Marshal(&GlobalConfig{})
	if err != nil {
		return err
	}
	return os.WriteFile(globalConfigFilePath, data, 0644)
}

func CreateGlobalConfig() error {
	if paths.FileAlreadyExist(globalConfigFilePath) {
		return nil
	}
	// Move file to keep the original clean
	if err := files.CopyFile(
		filepath.Join(paths.BashConfigAppDir, constants.GlobalConfigFile),
		globalConfigFilePath,
	); err != nil {
		return err
	}
	return ResetGlobalConfig()
}
