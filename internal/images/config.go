package images

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	BucketName  string `yaml:"bucket"`
	Region      string `yaml:"region"`
	URL         string `yaml:"url"`
	PartitionId string `yaml:"partitionId"`
	Port        int    `yaml:"port"`
	Address     string `yaml:"address"`
}



func NewConfig(path string) (Config, error) {
	cfg := Config{}

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return cfg, fmt.Errorf("failed to read config file: %w", err)
	}

	return cfg, nil
}
