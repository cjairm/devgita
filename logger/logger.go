package logger

import (
	"go.uber.org/zap"
)

var log *zap.SugaredLogger

// Init initializes the global logger
func Init(verbose bool) {
	var zapLogger *zap.Logger
	var err error

	if verbose {
		zapLogger, err = zap.NewDevelopment()
	} else {
		zapLogger, err = zap.NewProduction()
	}

	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}

	log = zapLogger.Sugar()
}

// L returns the global logger instance
func L() *zap.SugaredLogger {
	if log == nil {
		panic("logger not initialized. Call logger.Init() first.")
	}
	return log
}
