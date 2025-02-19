package config

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

const ConfigKey = "devgita-config"

type GlobalConfig struct {
	LocalConfigPath      string            `json:"localConfigPath"`
	RemoteConfigPath     string            `json:"remoteConfigPath"`
	SelectedTheme        string            `json:"selectedTheme"`
	Font                 string            `json:"font"`
	InstalledPackages    []string          `json:"installedPackages"`
	InstalledDesktopApps []string          `json:"installedDesktopApps"`
	Shortcuts            map[string]string `json:"shortcuts"`
}

type ContextConfig struct {
	SelectedLanguages []string
	SelectedDbs       []string
}

// Function to store the config in context
func WithConfig(ctx context.Context, config ContextConfig) context.Context {
	return context.WithValue(ctx, ConfigKey, config)
}

// Function to retrieve the config from context
func GetConfig(ctx context.Context) (ContextConfig, bool) {
	config, ok := ctx.Value(ConfigKey).(ContextConfig)
	return config, ok
}

func GetDevgitaPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	// DO NOT CHANGE THIS PATH
	return filepath.Join(homeDir, ".local", "share", "devgita"), nil
}

func GetLocalConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".config"), nil
}

func GetDevgitaConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	localConfig := filepath.Join(homeDir, ".config", "devgita")
	return localConfig, err
}

func LoadConfig(filename string) (*GlobalConfig, error) {
	file, err := ioutil.ReadFile(filename)
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

func SaveConfig(filename string, config *GlobalConfig) error {
	file, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, file, 0644)
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
