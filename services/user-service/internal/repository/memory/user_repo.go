package memory

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"gitlab.com/4uvirik/my-platform/services/user-service/internal/domain"
)

var ErrNotFound = errors.New("user not found")

type UserRepository struct {
	mu    sync.RWMutex
	store map[string]*domain.User
}

func New() *UserRepository {
	return &UserRepository{
		store: make(map[string]*domain.User),
	}
}

func (r *UserRepository) Create(email, name string) (*domain.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	u := &domain.User{
		ID:        uuid.NewString(),
		Email:     email,
		Name:      name,
		CreatedAt: time.Now(),
	}

	r.store[u.ID] = u
	return u, nil
}

func (r *UserRepository) Get(id string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	u, ok := r.store[id]
	if !ok {
		return nil, ErrNotFound
	}
	return u, nil
}
