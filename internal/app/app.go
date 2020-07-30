package app

import (
	"github.com/alonelegion/service_template/internal/app/config"
	"go.uber.org/zap"
)

type Application struct {
	config *config.AppConfig
	logger *zap.Logger

	Name        string
	Version     string
	Environment string
}

func New(
	name, version, environment string,
	config *config.AppConfig,
	logger *zap.Logger,
) *Application {
	return &Application{
		config: config,
		logger: logger,

		Name:        name,
		Version:     version,
		Environment: environment,
	}
}
