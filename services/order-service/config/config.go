package config

import (
	"time"

	"github.com/caarlos0/env/v11"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	App    AppConfig    `yaml:"app"`
	Server ServerConfig `yaml:"server"`
	Logger LoggerConfig `yaml:"logger"`
	DB     DBConfig     `yaml:"db"`
	Kafka  KafkaConfig  `yaml:"kafka"`
	Redis  RedisConfig  `yaml:"redis"`
}

type AppConfig struct {
	Name string `env:"APP_NAME" yaml:"name" envDefault:"order-service"`
	Env  string `env:"APP_ENV"  yaml:"env"  envDefault:"local.yaml"`
}

type ServerConfig struct {
	GRPCAddr        string        `env:"GRPC_ADDR" yaml:"grpc_addr" envDefault:":50051"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT" yaml:"shutdown_timeout" envDefault:"10s"`
}

type LoggerConfig struct {
	Level string `env:"LOG_LEVEL" yaml:"level" envDefault:"info"`
	JSON  bool   `env:"LOG_JSON"  yaml:"json"  envDefault:"true"`
}

type DBConfig struct {
	DSN string `env:"DB_DSN" yaml:"dsn"`
}

type KafkaConfig struct {
	Brokers           []string `env:"KAFKA_BROKERS" yaml:"brokers" envSeparator:","`
	TopicOrderCreated string   `env:"KAFKA_TOPIC_ORDER_CREATED" yaml:"topic_order_created"`
	GroupID           string   `env:"KAFKA_GROUP_ID" yaml:"group_id"`
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
