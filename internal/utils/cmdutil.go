package utils

import (
	"os"

	"go.uber.org/zap"
)

func NewLogger(service string) *zap.Logger {
	env := os.Getenv("ENV")

	logger, _ := zap.NewProduction(zap.Fields(
		zap.String("env", env),
		zap.String("service", service),
	))

	if env == "" {
		logger, _ = zap.NewDevelopment()
	}

	return logger
}
