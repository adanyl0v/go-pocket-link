package service

import (
	"context"
	"github.com/google/uuid"
	"go-pocket-link/internal/domain"
	"go-pocket-link/internal/repository"
	"go-pocket-link/pkg/auth/jwt"
	"time"
)

type AuthService struct {
	repository      repository.SessionsRepository
	jwtTokenManager jwt.TokenManager
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewAuthService(repo repository.SessionsRepository, jwtTokenManager jwt.TokenManager, accessTokenTTL, refreshTokenTTL time.Duration) *AuthService {
	return &AuthService{
		repository:      repo,
		jwtTokenManager: jwtTokenManager,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

func (s *AuthService) Access(accessToken string) (uuid.UUID, error) {
	sub, err := s.jwtTokenManager.Parse(accessToken)
	if err != nil {
		return uuid.Nil, err
	}
	return uuid.Parse(sub)
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (jwt.TokenPair, error) {
	session := domain.Session{RefreshToken: refreshToken}
	err := s.repository.GetByRefreshToken(ctx, &session)
	if err != nil {
		return jwt.TokenPair{}, err
	}
	accessTok, err := s.jwtTokenManager.CreateAccessToken(session.UserID.String(), s.accessTokenTTL)
	if err != nil {
		return jwt.TokenPair{}, err
	}
	refreshTok, err := s.jwtTokenManager.CreateRefreshToken(session.UserID.String(), s.refreshTokenTTL)
	if err != nil {
		return jwt.TokenPair{}, err
	}
	session.RefreshToken = refreshTok
	session.ExpiresAt = time.Now().Add(s.refreshTokenTTL)
	return jwt.TokenPair{AccessToken: accessTok, RefreshToken: refreshTok}, nil
}

func (s *AuthService) Login(ctx context.Context, userID uuid.UUID) (jwt.TokenPair, error) {
	session := domain.Session{UserID: userID}
	err := s.repository.GetByUserID(ctx, &session)
	if err != nil {
		return jwt.TokenPair{}, err
	}
	tokenPair := jwt.TokenPair{RefreshToken: session.RefreshToken}
	accessToken, err := s.jwtTokenManager.CreateAccessToken(userID.String(), s.accessTokenTTL)
	if err != nil {
		return jwt.TokenPair{}, err
	}
	tokenPair.AccessToken = accessToken
	return tokenPair, nil
}

func (s *AuthService) Signup(ctx context.Context, userID uuid.UUID) (jwt.TokenPair, error) {
	refreshToken, err := s.jwtTokenManager.CreateRefreshToken(userID.String(), s.refreshTokenTTL)
	if err != nil {
		return jwt.TokenPair{}, err
	}
	accessToken, err := s.jwtTokenManager.CreateAccessToken(userID.String(), s.accessTokenTTL)
	if err != nil {
		return jwt.TokenPair{}, err
	}
	session := domain.Session{
		UserID:       userID,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(s.refreshTokenTTL),
		CreatedAt:    time.Now(),
	}
	err = s.repository.Save(ctx, &session)
	if err != nil {
		return jwt.TokenPair{}, err
	}
	return jwt.TokenPair{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}
