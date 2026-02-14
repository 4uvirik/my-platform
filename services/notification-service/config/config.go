package config

import (
	"os"
	"time"

	"github.com/caarlos0/env/v11"
	"gopkg.in/yaml.v3"
)

type Config struct {
	App    AppConfig    `yaml:"app"`
	Logger LoggerConfig `yaml:"logger"`
	Kafka  KafkaConfig  `yaml:"kafka"`
	Redis  RedisConfig  `yaml:"redis"`
}

type AppConfig struct {
	Name string `env:"APP_NAME" yaml:"name" envDefault:"notification-service"`
	Env  string `env:"APP_ENV"  yaml:"env"  envDefault:"local.yaml"`
}

type LoggerConfig struct {
	Level string `env:"LOG_LEVEL" yaml:"level" envDefault:"info"`
	JSON  bool   `env:"LOG_JSON"  yaml:"json"  envDefault:"true"`
}

type KafkaConfig struct {
	Brokers           []string      `env:"KAFKA_BROKERS" yaml:"brokers" envSeparator:","`
	TopicOrderCreated string        `env:"KAFKA_TOPIC_ORDER_CREATED" yaml:"topic_order_created"`
	GroupID           string        `env:"KAFKA_GROUP_ID" yaml:"group_id" envDefault:"notification-service"`
	ConsumeTimeout    time.Duration `env:"KAFKA_CONSUME_TIMEOUT" yaml:"consume_timeout" envDefault:"10s"`
}

type RedisConfig struct {
	Addr     string `env:"REDIS_ADDR" yaml:"addr" envDefault:"localhost:6379"`
	Password string `env:"REDIS_PASSWORD" yaml:"password"`
	DB       int    `env:"REDIS_DB" yaml:"db" envDefault:"0"`
}

func Load(path string) (*Config, error) {
	cfg := &Config{}

	if path != "" {
		b, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		if err := yaml.Unmarshal(b, cfg); err != nil {
			return nil, err
		}
	}

	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
