package app

import (
	"context"
	"gitlab.com/4uvirik/my-platform/services/api-gateway/internal/server/httpserver"
	"log/slog"
	"net/http"
	"time"
)

type App struct {
	httpServer *httpserver.Server
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

func (a *App) Run() error {
	return a.httpServer.Start()
}

func (a *App) Shutdown(ctx context.Context) error {
	return a.httpServer.Shutdown(ctx)
}
