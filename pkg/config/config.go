package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Port string `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		URL  string `yaml:"url"`
		Name string `yaml:"name"`
	} `yaml:"database"`
	Migrations struct {
		Path string `yaml:"path"`
	} `yaml:"migrations"`
}

func NewConfig() (*Config, error) {
	path := os.Getenv("CONFIG_PATH")
	data, err := os.ReadFile(path)

	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", path, err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	return &cfg, nil
}
