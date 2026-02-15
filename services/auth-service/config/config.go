package config

import (
	"os"
)

type Config struct {
	HTTPAddr      string
	JWTSecret     string
	AccessTTLMin  int
	RefreshTTLMin int
}

func Load() (*Config, error) {
	cfg := &Config{
		HTTPAddr:      getEnv("AUTH_HTTP_ADDR", ":8081"),
		JWTSecret:     getEnv("AUTH_JWT_SECRET", "dev-secret"),
		AccessTTLMin:  15,
		RefreshTTLMin: 60 * 24,
	}

	return cfg, nil
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
