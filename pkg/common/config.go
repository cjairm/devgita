package common

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

const (
	configKey string = "config"
)

type Config struct {
	SelectedLanguages []string
	SelectedDbs       []string
}

// Function to store the config in context
func WithConfig(ctx context.Context, config Config) context.Context {
	return context.WithValue(ctx, configKey, config)
}

// Function to retrieve the config from context
func GetConfig(ctx context.Context) (Config, bool) {
	config, ok := ctx.Value(configKey).(Config)
	return config, ok
}

func GetDevgitaPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error getting home directory: %w", err)
	}
	devgitaPath := filepath.Join(homeDir, ".local", "share", "devgita")
	return devgitaPath, nil
}
