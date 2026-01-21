package logger

import (
	"log/slog"
	"os"

	"github.com/4uvirik/ProductService/pkg/logger/handler/slogpretty"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func ConfigLogger(logLevel string) (*slog.Logger, error) {
	var log *slog.Logger

	switch logLevel {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug, AddSource: true}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log, nil
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: true,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
