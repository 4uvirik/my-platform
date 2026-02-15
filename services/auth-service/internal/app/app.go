package app

import (
	"context"
	"gitlab.com/4uvirik/my-platform/services/auth-service/config"
	"log/slog"
	"net/http"
	"time"

	"gitlab.com/4uvirik/my-platform/services/auth-service/internal/jwt"
	"gitlab.com/4uvirik/my-platform/services/auth-service/internal/repository"
	"gitlab.com/4uvirik/my-platform/services/auth-service/internal/server/httpserver"
	"gitlab.com/4uvirik/my-platform/services/auth-service/internal/service"
)

type App struct {
	httpServer *httpserver.Server
	logger     *slog.Logger
}

func New(cfg *config.Config, logger *slog.Logger) (*App, error) {
	userRepo := repository.NewMemoryUserRepo()
	jwtMgr := jwt.NewManager(cfg.JWTSecret, time.Minute*time.Duration(cfg.AccessTTLMin), time.Minute*time.Duration(cfg.RefreshTTLMin))
	authSvc := service.NewAuthService(userRepo, jwtMgr)

	srv := httpserver.New(cfg.HTTPAddr, authSvc, logger)

	return &App{
		httpServer: srv,
		logger:     logger,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	a.logger.Info("auth-service started", "addr", a.httpServer.Addr())

	go func() {
		if err := a.httpServer.Start(); err != nil && err != http.ErrServerClosed {
			a.logger.Error("http server error", "error", err)
		}
	}()

	<-ctx.Done()

	a.logger.Info("shutting down auth-service")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return a.httpServer.Shutdown(shutdownCtx)
}
