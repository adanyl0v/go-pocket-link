package repository

import (
	"context"
	"github.com/google/uuid"
	"go-pocket-link/internal/domain"
	"go-pocket-link/internal/repository/postgres"
	"go-pocket-link/pkg/storage"
)

type Repositories struct {
	Users    UsersRepository
	Links    LinksRepository
	Sessions SessionsRepository
}

func NewRepositories(db storage.DB) *Repositories {
	return &Repositories{
		Users:    postgres.NewUsersRepository(db),
		Links:    postgres.NewLinksRepository(db),
		Sessions: postgres.NewSessionsRepository(db),
	}
}

type UsersRepository interface {
	// Save user and store his ID in [domain.User.ID]
	Save(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, dest *domain.User) error
	GetByCredentials(ctx context.Context, dest *domain.User) error
	GetAll(ctx context.Context) ([]domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type LinksRepository interface {
	// Save link and store its ID in [domain.Link.ID]
	Save(ctx context.Context, link *domain.Link) error
	GetByID(ctx context.Context, dest *domain.Link) error
	GetByURL(ctx context.Context, dest *domain.Link) error
	GetAll(ctx context.Context) ([]domain.Link, error)
	GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Link, error)
	GetAllByTitle(ctx context.Context, userID uuid.UUID, title string) ([]domain.Link, error)
	Update(ctx context.Context, link *domain.Link) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type SessionsRepository interface {
	// Save link and store its ID in [domain.Session.ID]
	Save(ctx context.Context, session *domain.Session) error
	GetByID(ctx context.Context, dest *domain.Session) error
	GetByRefreshToken(ctx context.Context, dest *domain.Session) error
	GetAll(ctx context.Context) ([]domain.Session, error)
	GetAllByUserID(ctx context.Context, id uuid.UUID) ([]domain.Session, error)
	Update(ctx context.Context, session *domain.Session) error
	Delete(ctx context.Context, id uuid.UUID) error
	//DeleteByRefreshToken(ctx context.Context, token string) error
}
