package httpserver

import (
	"log/slog"
	"net/http"
	"time"
)

type Logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
}

type AuthService interface {
	Register(email, password string) error
	Login(email, password string) (access, refresh string, err error)
}

func LoggingMiddleware(logger Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			logger.Info("http request",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Duration("duration", time.Since(start)),
			)
		})
	}
}
