package handlers

import (
	"encoding/json"
	"net/http"

	"gitlab.com/4uvirik/my-platform/services/user-service/internal/repository/memory"
)

type UserHandler struct {
	repo *memory.UserRepository
}

func NewUserHandler(repo *memory.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	u, err := h.repo.Create(req.Email, req.Name)
	if err != nil {
		http.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(u)
}
