package config

import (
	"os"
)

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type Config struct {
	Port           string
	ServiceName    string
	DatabaseConfig DatabaseConfig
	AuditServiceURL string
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func NewConfig() *Config {
	return &Config{
		Port:        getEnv("PORT", "8081"),
		ServiceName: getEnv("SERVICE_NAME", "request-service"),
		DatabaseConfig: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "request_db"),
		},
		AuditServiceURL: getEnv("AUDIT_SERVICE_URL", "http://localhost:8083"),
	}
}

func (c *Config) GetDSN() string {
	return "host=" + c.DatabaseConfig.Host +
		" port=" + c.DatabaseConfig.Port +
		" user=" + c.DatabaseConfig.User +
		" password=" + c.DatabaseConfig.Password +
		" dbname=" + c.DatabaseConfig.DBName +
		" sslmode=disable"
}
