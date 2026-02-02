package handlers

import (
	"encoding/json"
	"net/http"
)

type AuthService interface {
	Register(email, password string) error
	Login(email, password string) (access, refresh string, err error)
}

type AuthHandler struct {
	svc AuthService
}

func NewAuthHandler(svc AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

type authReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authResp struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req authReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if err := h.svc.Register(req.Email, req.Password); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req authReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	access, refresh, err := h.svc.Login(req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	_ = json.NewEncoder(w).Encode(authResp{Access: access, Refresh: refresh})
}
