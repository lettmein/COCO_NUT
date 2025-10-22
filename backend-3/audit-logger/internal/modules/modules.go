package modules

import (
	"database/sql"
	"fmt"
	"myapp/internal/config"
	"myapp/internal/handler"
	"myapp/internal/repository"
	"myapp/internal/router"
	"myapp/internal/server"
	"myapp/internal/service"
	"myapp/pkg/logger"

	_ "github.com/lib/pq" // Импорт драйвера PostgreSQL
	"go.uber.org/fx"
)

// Re-export types for use in main
type Config = config.Config
type HTTPServer = server.HTTPServer
type Logger = logger.Logger

var Module = fx.Module("myapp",
	fx.Provide(
		config.New,
		logger.New,
		// Добавляем создание подключения к БД
		newDatabase,
	),
)

var RepositoryModule = fx.Module("repositories",
	fx.Provide(
		repository.NewLoggerRepository,
	),
)

var ServiceModule = fx.Module("services",
	fx.Provide(
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

// Функция для создания подключения к БД
func newDatabase(config *Config, log *Logger) (*sql.DB, error) {
	// Формируем строку подключения для PostgreSQL
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Database.Host,
		config.Database.Port,
		config.Database.User,
		config.Database.Password,
		config.Database.Name,
		config.Database.SSLMode,
	)

	log.Infof("Connecting to database: %s@%s:%s", config.Database.User, config.Database.Host, config.Database.Port)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Проверяем подключение
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Настраиваем пул соединений
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * 60) // 5 minutes

	log.Info("Successfully connected to database")
	return db, nil
}
