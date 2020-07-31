package app

import (
	"context"

	"github.com/alonelegion/service_template/internal/app/config"
	grpc "github.com/alonelegion/service_template/internal/app/grpc/server"
	http "github.com/alonelegion/service_template/internal/app/http/server"
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

func (app *Application) Run(ctx context.Context) {
	grpcServerErrCh := grpc.NewServer(ctx, app.logger, app.config)
	httpServerErrCh := http.NewServer(ctx, app.logger, app.config)

	select {
	case <-grpcServerErrCh:
		<-httpServerErrCh
	case <-httpServerErrCh:
		<-grpcServerErrCh
	}
}

func (app *Application) Shutdown() {
	_ = app.logger.Sync()
}
