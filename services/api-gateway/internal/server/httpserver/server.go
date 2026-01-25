package httpserver

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
	logger     *slog.Logger
}

func New(addr string, handler http.Handler, readTimeout time.Duration, writeTimeuot time.Duration, errorLogger *slog.Logger) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:         addr,
			Handler:      handler,
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeuot,
		},
		logger: errorLogger,
	}
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
