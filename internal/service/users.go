package service

import (
	"go-pocket-link/internal/repository"
	"go-pocket-link/pkg/crypto/hash"
)

type UsersService struct {
	repo   repository.UsersRepository
	hasher hash.Hasher
}

func NewUsersService(repo repository.UsersRepository, hasher hash.Hasher) *UsersService {
	return &UsersService{repo: repo, hasher: hasher}
}

func (s *UsersService) Repository() repository.UsersRepository {
	return s.repo
}

func (s *UsersService) HashPassword(password string) string {
	return s.hasher.Hash(password)
}
