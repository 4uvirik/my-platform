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

func New(addr string, handler any, errorLogger *slog.Logger) *App {
	server := httpserver.New(
		addr,
		handler.(interface {
			ServeHTTP(http.ResponseWriter, *http.Request)
		}),
		5*time.Second,
		5*time.Second,
		errorLogger,
	)

	return &App{
		httpServer: server,
	}
}

func (a *App) Run(ctx context.Context) error {
	ctx = server.SetupSignalContext(ctx)

	go func() {
		if err := a.httpServer.Start(); err != nil && err != http.ErrServerClosed {
			a.logger.Error("http server error", "error", err)
		}
	}()

	<-ctx.Done()

	shutDownCtx, cancel := context.WithTimeout(context.Background(), a.cfg.Server.ShutdownTimeout)
	defer cancel()

	return a.httpServer.Shutdown(shutDownCtx)
}

func (a *App) Shutdown(ctx context.Context) error {
	return a.httpServer.Shutdown(ctx)
}
