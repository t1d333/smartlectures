package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Address           string `yaml:"address"            env:"ADDRESS"            env-default:":8000"`
	RecognizerAddress string `yaml:"recognizer_address" env:"RECOGNIZER_ADDRESS"`
	StorageAddress    string `yaml:"storage_address"    env:"STORAGE_ADDRESS"`
	PostgresURL       string `yaml:"postgres_url"       env:"DB_URL"`
}

func NewConfig(path string) (Config, error) {
	cfg := Config{}

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return cfg, fmt.Errorf("failed to read config file: %w", err)
	}

	return cfg, nil
}
