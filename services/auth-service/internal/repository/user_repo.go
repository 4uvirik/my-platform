package repository

import "gitlab.com/4uvirik/my-platform/services/auth-service/internal/domain"

type UserRepository interface {
	Create(u domain.User) error
	GetByEmail(email string) (*domain.User, error)
}
