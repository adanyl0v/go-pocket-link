package repository

import (
	"context"
	"github.com/google/uuid"
	"go-pocket-link/internal/domain"
)

type UsersRepository interface {
	Save(ctx context.Context, user *domain.User) error
	Get(ctx context.Context, id uuid.UUID) (domain.User, error)
	GetByCredentials(ctx context.Context, email, password string) (domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type Repositories struct {
	Users UsersRepository
}
