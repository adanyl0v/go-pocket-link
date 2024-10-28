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
	// Update domain.User Name, Email and Password by ID
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type SessionsRepository interface {
	Save(ctx context.Context, session *domain.Session) error
	Get(ctx context.Context, id uuid.UUID) (domain.Session, error)
	GetByUserID(ctx context.Context, id uuid.UUID) (domain.Session, error)
	GetByRefreshToken(ctx context.Context, token string) (domain.Session, error)
	// Update domain.Session RefreshToken and IsInvoked by ID
	Update(ctx context.Context, session *domain.Session) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type Repositories struct {
	Users    UsersRepository
	Sessions SessionsRepository
}
