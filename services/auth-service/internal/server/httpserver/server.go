package httpserver

import (
	"context"
	"net/http"
)

type Server struct {
	srv *http.Server
}

func New(addr string, authSvc AuthService, logger Logger) *Server {
	router := NewRouter(authSvc, logger)

	return &Server{
		srv: &http.Server{
			Addr:    addr,
			Handler: router,
		},
	}
}

func (s *Server) Start() error {
	return s.srv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func (s *Server) Addr() string {
	return s.srv.Addr
}
