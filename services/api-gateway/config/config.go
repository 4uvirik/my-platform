package config

import (
	"fmt"
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

type Config struct {
	App    AppConfig    `yaml:"app"`
	Server ServerConfig `yaml:"server"`
	Logger LoggerConfig `yaml:"logger"`
}

type AppConfig struct {
	Name string `env:"APP_NAME" yaml:"name"`
	Env  string `env:"APP_ENV" envDefault:"local" yaml:"env"`
}

type ServerConfig struct {
	Host            string        `env:"SERVER_HOST" yaml:"host"`
	Port            string        `env:"SERVER_PORT" yaml:"port"`
	ShutdownTimeout time.Duration `env:"SERVER_SHUTDOWN_TIMEOUT" envDefault:"5s" yaml:"shutdown_timeout"`
}

type LoggerConfig struct {
	Level string `env:"LOG_LEVEL" envDefault:"info" yaml:"level"`
}

// Load - загрузка конфига сначала из .env потом из yaml.
func Load(yamlPath string) (*Config, error) {
	cfg := &Config{}

	var err error

	// 1. Загрузка из .env
	err = loadFromEnv(cfg)
	if err == nil {
		return cfg, nil
	}

	// 2. Загрузка из yaml
	err = loadFromYaml(yamlPath, cfg)
	if err == nil {
		return cfg, nil
	}

	return nil, fmt.Errorf("failed to load config: %w", err)
}

func loadFromEnv(cfg *Config) error {
	if _, err := os.Stat(".env"); err == nil {
		_ = godotenv.Load(".env")
	}

	if err := env.Parse(cfg); err != nil {
		return fmt.Errorf("env parse error: %w", err)
	}

	return nil
}

func loadFromYaml(path string, cfg *Config) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read yaml error: %w", err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return fmt.Errorf("yaml unmarshal error: %w", err)
	}

	return nil
}
