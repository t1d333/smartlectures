package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	DragonflyURL      string `yaml:"dragnofly_url"      env:"DRAGONFLY_URL"`
	DragonflyPassword string `yaml:"dragonfly_password" env:"DRAGONFLY_PASSWORD"`
	PostgresURL       string `yaml:"postgres_url"       env:"DB_URL"`
}

func NewConfig(path string) (Config, error) {
	cfg := Config{}

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return cfg, fmt.Errorf("failed to read config file: %w", err)
	}

	return cfg, nil
}
