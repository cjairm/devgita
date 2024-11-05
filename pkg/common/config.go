package common

import (
	"context"
)

type contextKey string

const (
	configKey contextKey = "config"
)

type Config struct {
	SelectedLanguages []string
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
