package service

import (
	"context"
	"github.com/google/uuid"
	"go-pocket-link/internal/domain"
	"go-pocket-link/internal/repository"
	"go-pocket-link/pkg/crypto/hash"
	"time"
)

type UsersService struct {
	repo   repository.UsersRepository
	hasher hash.Hasher
}

func NewUsersService(repo repository.UsersRepository, hasher hash.Hasher) *UsersService {
	return &UsersService{repo: repo, hasher: hasher}
}

func (s *UsersService) Save(ctx context.Context, user *domain.User) error {
	user.Password = s.hasher.Hash(user.Password)
	user.CreatedAt = time.Now()
	return s.repo.Save(ctx, user)
}

func (s *UsersService) Get(ctx context.Context, id uuid.UUID) (domain.User, error) {
	return s.repo.Get(ctx, id)
}

func (s *UsersService) GetByCredentials(ctx context.Context, email, password string) (domain.User, error) {
	return s.repo.GetByCredentials(ctx, email, s.hasher.Hash(password))
}

func (s *UsersService) Update(ctx context.Context, user *domain.User) error {
	user.Password = s.hasher.Hash(user.Password)
	return s.repo.Update(ctx, user)
}

func (s *UsersService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
