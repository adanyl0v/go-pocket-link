package service

import (
	"context"
	"github.com/google/uuid"
	"go-pocket-link/internal/domain"
	"go-pocket-link/internal/repository"
	"time"
)

type SessionsService struct {
	repo repository.SessionsRepository
}

func NewSessionsService(repo repository.SessionsRepository) *SessionsService {
	return &SessionsService{repo: repo}
}

func (s *SessionsService) Save(ctx context.Context, session *domain.Session) error {
	session.CreatedAt = time.Now()
	return s.repo.Save(ctx, session)
}

func (s *SessionsService) Get(ctx context.Context, id uuid.UUID) (domain.Session, error) {
	return s.repo.Get(ctx, id)
}

func (s *SessionsService) GetByUserID(ctx context.Context, id uuid.UUID) (domain.Session, error) {
	return s.repo.GetByUserID(ctx, id)
}

func (s *SessionsService) GetByRefreshToken(ctx context.Context, token string) (domain.Session, error) {
	return s.repo.GetByRefreshToken(ctx, token)
}

func (s *SessionsService) Update(ctx context.Context, session *domain.Session) error {
	return s.repo.Update(ctx, session)
}

func (s *SessionsService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *SessionsService) Invalidate(ctx context.Context, id uuid.UUID) error {
	session := domain.Session{ID: id, RefreshToken: "0", IsInvoked: true}
	return s.repo.Update(ctx, &session)
}
