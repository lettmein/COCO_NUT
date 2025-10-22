package config

import (
	"fmt"
	"time"
)

type Config struct {
	Environment string
	HTTP        HTTPConfig
	Database    DatabaseConfig
}

type HTTPConfig struct {
	Address string
	Port    int
	Timeout time.Duration
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func New() *Config {
	return &Config{
		Environment: "development",
		HTTP: HTTPConfig{
			Address: ":8080",
			Port:    8080,
			Timeout: 30 * time.Second,
		},
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "user",
			Password: "password",
			DBName:   "myapp",
			SSLMode:  "disable",
		},
	}
}

func (db *DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		db.Host, db.Port, db.User, db.Password, db.DBName, db.SSLMode)
}
