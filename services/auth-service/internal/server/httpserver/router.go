package httpserver

import (
	"net/http"

	"gitlab.com/4uvirik/my-platform/services/auth-service/internal/transport/http/handlers"
)

func NewRouter(authSvc AuthService, logger Logger) http.Handler {
	mux := http.NewServeMux()

	health := handlers.NewHealthHandler()
	auth := handlers.NewAuthHandler(authSvc)

	mux.HandleFunc("/health", health.Health)
	mux.HandleFunc("/auth/register", auth.Register)
	mux.HandleFunc("/auth/login", auth.Login)

	return LoggingMiddleware(logger)(mux)
}
