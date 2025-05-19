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

func LoadGlobalConfig() (*GlobalConfig, error) {
	file, err := os.ReadFile(
		filepath.Join(paths.ConfigDir, constants.AppName, "global_config.yaml"),
	)
	if err != nil {
		return nil, err
	}
	var config GlobalConfig
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func SetGlobalConfig(config *GlobalConfig) error {
	filePath := filepath.Join(paths.ConfigDir, constants.AppName, "global_config.yaml")
	file, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, file, 0644)
}

func ResetGlobalConfig() error {
	filePath := filepath.Join(paths.ConfigDir, constants.AppName, "global_config.yaml")
	data, err := yaml.Marshal(&GlobalConfig{})
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0644)
}

func CreateGlobalConfig() error {
	newConfigFile := filepath.Join(paths.ConfigDir, constants.AppName, "global_config.yaml")
	if paths.FileAlreadyExist(newConfigFile) {
		return nil
	}
	if err := files.CopyFile(
		filepath.Join(paths.BashConfigAppDir, "global_config.yaml"),
		newConfigFile,
	); err != nil {
		return err
	}
	return ResetGlobalConfig()
}

//Example of how to use the config package
// configFile := "./configs/bash/devgita_config.json"
//
// // Load the configuration
// c, err := config.LoadConfig(configFile)
// if err != nil {
// 	fmt.Println("Error loading config:", err)
// 	return
// }
//
// // Print the loaded configuration
// fmt.Printf("Loaded Config: %+v\n", c)
//
// // Modify the configuration
// c.SelectedTheme = "light"
// c.InstalledPackages = append(c.InstalledPackages, "new-package")
//
// // Save the updated configuration
// err = config.SaveConfig(configFile, c)
// if err != nil {
// 	fmt.Println("Error saving config:", err)
// 	return
// }
//
// fmt.Println("Configuration saved successfully.")
