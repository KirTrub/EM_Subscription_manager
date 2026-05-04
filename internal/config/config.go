package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	App struct {
		Name string `yaml:"name"`
		Port int    `yaml:"port"`
	} `yaml:"app"`

	Postgres struct {
		DSN string `yaml:"dsn"`
	} `yaml:"postgres"`

	Logger struct {
		Level string `yaml:"level"`
	} `yaml:"logger"`

	Swagger struct {
		Enabled bool `yaml:"enabled"`
	} `yaml:"swagger"`
}

func Load(path string) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return nil, err
	}

	if cfg.App.Port == 0 {
		cfg.App.Port = 3000
	}
	if cfg.Logger.Level == "" {
		cfg.Logger.Level = "info"
	}
	return &cfg, nil
}
