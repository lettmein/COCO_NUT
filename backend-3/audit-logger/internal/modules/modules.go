package modules

import (
	"myapp/internal/config"
	"myapp/internal/handler"
	"myapp/internal/repository"
	"myapp/internal/router"
	"myapp/internal/server"
	"myapp/internal/service"
	"myapp/pkg/logger"

	"go.uber.org/fx"
)

var Module = fx.Module("myapp",
	fx.Provide(
		config.New,
		logger.New,
	),
)

var RepositoryModule = fx.Module("repositories",
	fx.Provide(
		repository.NewLoggerRepository,
	),
)

var ServiceModule = fx.Module("services",
	fx.Provide(
		// Исправлено: предоставляем конкретный тип *service.LoggerService
		service.NewLoggerService,
	),
)

var HandlerModule = fx.Module("handlers",
	fx.Provide(
		handler.NewLoggerHandler,
	),
)

var InfrastructureModule = fx.Module("infrastructure",
	fx.Provide(
		router.NewRouter,
		server.NewHTTPServer,
	),
)

// Комбинированный модуль для удобства
var AppModule = fx.Module("application",
	Module,
	RepositoryModule,
	ServiceModule,
	HandlerModule,
	InfrastructureModule,
)
