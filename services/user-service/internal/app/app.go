package app

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"gitlab.com/4uvirik/my-platform/services/user-service/internal/server"
	"gitlab.com/4uvirik/my-platform/services/user-service/internal/server/httpserver"
)

type App struct {
	server *httpserver.Server
	logger *slog.Logger
}

func New(addr string, handler http.Handler, logger *slog.Logger) *App {
	srv := httpserver.New(addr, handler, 5*time.Second, 5*time.Second, logger)

	return &App{
		server: srv,
		logger: logger,
	}
}

func (a *App) Run(ctx context.Context) error {
	ctx = server.SetupSignalContext(ctx)

	a.logger.Info("user-service started")

	go func() {
		if err := a.server.Start(); err != nil && err != http.ErrServerClosed {
			a.logger.Error("http server error", "error", err)
		}
	}()

	<-ctx.Done()
	a.logger.Info("shutting down user-service")

	return a.server.Shutdown(context.Background())
}
