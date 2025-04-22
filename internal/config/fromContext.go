package config

import (
	"context"
)

const ConfigKey = "devgita-config-context"

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
