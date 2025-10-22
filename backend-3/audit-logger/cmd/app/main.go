package main

import (
	"context"
	"go.uber.org/fx"
	"log"
	"myapp/internal/config"
	"myapp/internal/handler"
	"myapp/internal/repository"
	"myapp/internal/service"
	"myapp/pkg/logger"
)

func main() {
	app := fx.New(
		// Модули приложения
		fx.Provide(
			config.New,
			logger.New,
			repository.NewRepository,
			service.NewService,
			handler.NewHandler,
		),

		// Запуск приложения
		fx.Invoke(registerHooks),
	)

	app.Run()
}

func registerHooks(
	lc fx.Lifecycle,
	handler *handler.Handler,
	config *config.Config,
	log *logger.Logger,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info("Starting application on " + config.HTTP.Address)
			go handler.Start(config.HTTP.Address)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("Stopping application")
			return nil
		},
	})
}
