package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	HTTP HTTPConfig
}

type HTTPConfig struct {
	Address string `envconfig:"HTTP_ADDRESS" default:":8080"`
}

func New() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
