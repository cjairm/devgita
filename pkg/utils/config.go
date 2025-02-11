package utils

import (
	"context"
	"os"
	"path/filepath"
)

type Config struct {
	SelectedLanguages []string
	SelectedDbs       []string
}

// Function to store the config in context
func WithConfig(ctx context.Context, config Config) context.Context {
	return context.WithValue(ctx, ConfigKey, config)
}

// Function to retrieve the config from context
func GetConfig(ctx context.Context) (Config, bool) {
	config, ok := ctx.Value(ConfigKey).(Config)
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
