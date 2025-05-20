package config

import (
	"os"
	"path/filepath"

	"github.com/cjairm/devgita/pkg/constants"
	"github.com/cjairm/devgita/pkg/files"
	"github.com/cjairm/devgita/pkg/paths"
	"gopkg.in/yaml.v3"
)

type GlobalConfig struct {
	AppPath               string            `yaml:"app_path"`
	ConfigPath            string            `yaml:"config_path"`
	AvailableFonts        []string          `yaml:"available_fonts"`
	AvailableThemes       []string          `yaml:"available_themes"`
	CurrentFont           string            `yaml:"current_font"`
	CurrentTheme          string            `yaml:"current_theme"`
	InstalledPackages     []string          `yaml:"installed_packages"`
	InstalledDesktopApps  []string          `yaml:"installed_desktop_apps"`
	InstalledDevLanguages []string          `yaml:"installed_dev_languages"`
	InstalledDatabases    []string          `yaml:"installed_databases"`
	Shortcuts             map[string]string `yaml:"shortcuts"`
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
