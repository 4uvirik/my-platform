package app

import (
	"context"
	"gitlab.com/4uvirik/my-platform/services/api-gateway/config"
	"gitlab.com/4uvirik/my-platform/services/api-gateway/internal/server"
	"gitlab.com/4uvirik/my-platform/services/api-gateway/internal/server/httpserver"
	"log/slog"
	"net/http"
	"time"
)

type App struct {
	cfg        *config.Config
	httpServer *httpserver.Server
	logger     *slog.Logger
}

func New(cfg *config.Config, handler http.Handler, logger *slog.Logger) *App {
	server := httpserver.New(
		cfg.Server.Addr(),
		handler,
		5*time.Second,
		5*time.Second,
		logger,
	)

	return &App{
		cfg:        cfg,
		httpServer: server,
		logger:     logger,
	}
}

func (a *App) Run(ctx context.Context) error {
	ctx = server.SetupSignalContext(ctx)

	a.logger.Info("http server started", "addr", a.cfg.Server.Addr())

	go func() {
		if err := a.httpServer.Start(); err != nil && err != http.ErrServerClosed {
			a.logger.Error("http server error", "error", err)
		}
	}()

	<-ctx.Done()

	a.logger.Info("shutting down application")

	shutDownCtx, cancel := context.WithTimeout(context.Background(), a.cfg.Server.ShutdownTimeout)
	defer cancel()

	return a.Shutdown(shutDownCtx)
}

func (a *App) Shutdown(ctx context.Context) error {
	return a.httpServer.Shutdown(ctx)
}
