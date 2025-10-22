package main

import (
	"context"
	"myapp/internal/modules"
	"myapp/internal/server"
	"myapp/pkg/logger"

	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		// Регистрация всех модулей
		modules.Module,               // Базовые зависимости (config, logger)
		modules.RepositoryModule,     // Репозитории
		modules.ServiceModule,        // Сервисы
		modules.HandlerModule,        // Обработчики
		modules.InfrastructureModule, // Инфраструктура (router, server)

		// Запуск приложения
		fx.Invoke(registerHooks),
	)

	app.Run()
}

func registerHooks(
	lc fx.Lifecycle,
	server *server.HTTPServer,
	log *logger.Logger,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info("Starting application")
			go server.Start()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("Stopping application")
			return server.Shutdown(ctx)
		},
	})
}
