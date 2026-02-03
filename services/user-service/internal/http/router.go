package http

import (
	"net/http"

	"gitlab.com/4uvirik/my-platform/services/user-service/internal/http/handlers"
	"gitlab.com/4uvirik/my-platform/services/user-service/internal/repository/memory"
)

func NewRouter(repo *memory.UserRepository) http.Handler {
	mux := http.NewServeMux()

	userHandler := handlers.NewUserHandler(repo)

	mux.HandleFunc("/health", handlers.Health)
	mux.HandleFunc("/users", userHandler.Create)

	return mux
}
