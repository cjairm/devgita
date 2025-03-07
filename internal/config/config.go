package config

import (
	"context"
	"os"
	"path/filepath"
)

const ConfigKey = "devgita-config"

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

func GetDevgitaConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	localConfig := filepath.Join(homeDir, ".config", "devgita")
	return localConfig, err
}
