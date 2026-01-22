package config

import (
	"fmt"
	"strings"
)

func (cfg *Config) Validate() error {
	var errorsMsg []string

	if cfg.App.Name == "" {
		errorsMsg = append(errorsMsg, "invalid app name")
	}

	if cfg.Server.Host == "" {
		errorsMsg = append(errorsMsg, "invalid app host")
	}

	if cfg.Server.Port == "" {
		errorsMsg = append(errorsMsg, "invalid app port")
	}

	if len(errorsMsg) > 0 {
		return fmt.Errorf("%s", strings.Join(errorsMsg, ": "))
	}

	return nil
}
