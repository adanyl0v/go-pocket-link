package service

import "go-pocket-link/internal/repository"

type LinksService struct {
	repo repository.LinksRepository
}

func NewLinksService(repo repository.LinksRepository) *LinksService {
	return &LinksService{repo: repo}
}

func (s *LinksService) Repository() repository.LinksRepository {
	return s.repo
}
