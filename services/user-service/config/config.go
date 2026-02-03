package config

import (
	"fmt"
	"os"
	"time"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Config struct {
	App    AppConfig    `yaml:"app"`
	Server ServerConfig `yaml:"server"`
	Logger LoggerConfig `yaml:"logger"`
}

type AppConfig struct {
	Name string `env:"APP_NAME" yaml:"name"`
	Env  string `env:"APP_ENV"  envDefault:"local" yaml:"env"`
}

type ServerConfig struct {
	Host            string        `env:"SERVER_HOST" yaml:"host"`
	Port            string        `env:"SERVER_PORT" yaml:"port"`
	ShutdownTimeout time.Duration `env:"SERVER_SHUTDOWN_TIMEOUT" envDefault:"5s" yaml:"shutdown_timeout"`
}

type LoggerConfig struct {
	Level string `env:"LOG_LEVEL" envDefault:"info" yaml:"level"`
}

func Load(yamlPath string) (*Config, error) {
	cfg := &Config{}

	// 1. .env
	if _, err := os.Stat(".env"); err == nil {
		_ = godotenv.Load(".env")
	}

	if err := env.Parse(cfg); err == nil {
		return cfg, nil
	}

	// 2. yaml
	if yamlPath == "" {
		return nil, fmt.Errorf("yaml path is empty and env config failed")
	}

	data, err := os.ReadFile(yamlPath)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) Addr() string {
	return c.Server.Host + ":" + c.Server.Port
}
