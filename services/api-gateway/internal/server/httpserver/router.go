package httpserver

import (
	"gitlab.com/4uvirik/my-platform/services/api-gateway/internal/server/httpserver/handlers"
	"net/http"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	health := handlers.NewHealthHandler()
	mux.HandleFunc("/health", health.Health)

	return mux
}
