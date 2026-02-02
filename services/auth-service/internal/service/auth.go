package service

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	"gitlab.com/4uvirik/my-platform/services/auth-service/internal/domain"
	"gitlab.com/4uvirik/my-platform/services/auth-service/internal/jwt"
	"gitlab.com/4uvirik/my-platform/services/auth-service/internal/repository"
)

type AuthService struct {
	repo   repository.UserRepository
	jwtMgr *jwt.Manager
}

func NewAuthService(repo repository.UserRepository, jwtMgr *jwt.Manager) *AuthService {
	return &AuthService{repo: repo, jwtMgr: jwtMgr}
}

func (s *AuthService) Register(email, password string) error {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return s.repo.Create(domain.User{
		Email:        email,
		PasswordHash: string(hash),
	})
}

func (s *AuthService) Login(email, password string) (string, string, error) {
	u, err := s.repo.GetByEmail(email)
	if err != nil {
		return "", "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return "", "", errors.New("invalid credentials")
	}

	access, refresh, err := s.jwtMgr.Generate(u.Email)
	return access, refresh, err
}
