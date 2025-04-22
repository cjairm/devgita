package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/cjairm/devgita/pkg/paths"
)

type GlobalConfig struct {
	AppPath               string            `json:"appPath"`
	AvailableFonts        []string          `json:"availableFonts"`
	AvailableThemes       []string          `json:"availableThemes"`
	CurrentFont           string            `json:"currentFont"`
	CurrentTheme          string            `json:"currentTheme"`
	InstalledPackages     []string          `json:"installedPackages"`
	InstalledDesktopApps  []string          `json:"installedDesktopApps"`
	InstalledDevLanguages []string          `json:"installedDevLanguages"`
	InstalledDatabases    []string          `json:"installedDatabases"`
	Shortcuts             map[string]string `json:"shortcuts"`
}

func LoadGlobalConfig() (*GlobalConfig, error) {
	file, err := os.ReadFile(filepath.Join(paths.BashConfigAppDir, "global_config.json"))
	if err != nil {
		return nil, err
	}
	var config GlobalConfig
	err = json.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func SetGlobalConfig(config *GlobalConfig) error {
	filePath := filepath.Join(paths.BashConfigAppDir, "global_config.json")
	file, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, file, 0644)
}

func ResetGlobalConfig() error {
	filePath := filepath.Join(paths.BashConfigAppDir, "global_config.json")
	return os.WriteFile(filePath, []byte("{}"), 0644)
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
