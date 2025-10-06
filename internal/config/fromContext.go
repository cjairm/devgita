package config

import (
	"context"

	"github.com/cjairm/devgita/pkg/logger"
)

const ConfigKey = "devgita-config-context"

type ContextConfig struct {
	SelectedLanguages []string
	SelectedDbs       []string
}

// Function to store the config in context
func WithConfig(ctx context.Context, config ContextConfig) context.Context {
	logger.L().Info("Storing config in context")
	return context.WithValue(ctx, ConfigKey, config)
}

// Function to retrieve the config from context
func GetConfig(ctx context.Context) (ContextConfig, bool) {
	logger.L().Info("Retrieving config from context")
	config, ok := ctx.Value(ConfigKey).(ContextConfig)
	return config, ok
}
