package service

import (
	"context"
	"github.com/adanyl0v/go-pocket-link/internal/domain"
	"github.com/adanyl0v/go-pocket-link/internal/repository"
	"github.com/adanyl0v/go-pocket-link/pkg/crypto/hash"
	"github.com/adanyl0v/go-pocket-link/pkg/validator"
	"github.com/google/uuid"
)

type UsersService struct {
	repo      repository.UsersRepository
	hasher    hash.Hasher
	validator *validator.CredentialsValidator
}

func NewUsersService(repo repository.UsersRepository, hasher hash.Hasher, validator *validator.CredentialsValidator) *UsersService {
	return &UsersService{
		repo:      repo,
		hasher:    hasher,
		validator: validator,
	}
}

func (s *UsersService) Save(ctx context.Context, user *domain.User) error {
	user.Password = s.hasher.Hash(user.Password)
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

func (s *UsersService) ComparePasswordAndHash(password, hashed string) bool {
	return s.hasher.Hash(password) == hashed
}

func (s *UsersService) ValidateName(name string) error {
	return s.validator.ValidateName(name)
}

func (s *UsersService) ValidateEmail(email string) error {
	return s.validator.ValidateEmail(email)
}

func (s *UsersService) ValidatePassword(password string) error {
	return s.validator.ValidatePassword(password)
}
