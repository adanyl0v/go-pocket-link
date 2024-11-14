package repository

import (
	"context"
	"github.com/adanyl0v/go-pocket-link/internal/domain"
	"github.com/google/uuid"
	"time"
)

type UsersRepository interface {
	Save(ctx context.Context, user *domain.User) error
	Get(ctx context.Context, id uuid.UUID) (domain.User, error)
	GetByCredentials(ctx context.Context, email, password string) (domain.User, error)
	// Update domain.User Name, Email and Password by ID
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type TokensRepository interface {
	Get(ctx context.Context, userID, tokenID uuid.UUID) (domain.Token, error)
	GetByKey(ctx context.Context, key string) (domain.Token, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Token, error)
	GetByTokenID(ctx context.Context, tokenID uuid.UUID) (domain.Token, error)
	Set(ctx context.Context, token *domain.Token, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
	DeleteByTokenID(ctx context.Context, tokenID uuid.UUID) error
}

type Repositories struct {
	Users  UsersRepository
	Tokens TokensRepository
}
