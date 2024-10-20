package service

import "go-pocket-link/internal/repository"

type SessionsService struct {
	repo repository.SessionsRepository
}

func NewSessionsService(repo repository.SessionsRepository) *SessionsService {
	return &SessionsService{repo: repo}
}

func (s *SessionsService) Repository() repository.SessionsRepository {
	return s.repo
}
