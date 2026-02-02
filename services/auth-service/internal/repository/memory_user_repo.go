package repository

import (
	"errors"
	"sync"

	"gitlab.com/4uvirik/my-platform/services/auth-service/internal/domain"
)

var ErrNotFound = errors.New("user not found")

type MemoryUserRepo struct {
	mu    sync.RWMutex
	users map[string]domain.User
}

func NewMemoryUserRepo() *MemoryUserRepo {
	return &MemoryUserRepo{
		users: make(map[string]domain.User),
	}
}

func (r *MemoryUserRepo) Create(u domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.users[u.Email]; ok {
		return errors.New("user already exists")
	}

	r.users[u.Email] = u
	return nil
}

func (r *MemoryUserRepo) GetByEmail(email string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	u, ok := r.users[email]
	if !ok {
		return nil, ErrNotFound
	}

	return &u, nil
}
