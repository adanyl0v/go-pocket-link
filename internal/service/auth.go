package service

import (
	jwt5 "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go-pocket-link/pkg/auth/jwt"
	"time"
)

type AuthService struct {
	jwtTm           jwt.TokenManager
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewAuthService(jwtTm jwt.TokenManager, accessTokenTTL, refreshTokenTTL time.Duration) *AuthService {
	return &AuthService{
		jwtTm:           jwtTm,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

func (s *AuthService) AccessTokenTTL() time.Duration {
	return s.accessTokenTTL
}

func (s *AuthService) RefreshTokenTTL() time.Duration {
	return s.refreshTokenTTL
}

func (s *AuthService) ParseAccessToken(token string) (jwt5.Claims, error) {
	return s.jwtTm.ParseAccessToken(token)
}

func (s *AuthService) ParseRefreshToken(token string) (jwt5.Claims, error) {
	return s.jwtTm.ParseRefreshToken(token)
}

func (s *AuthService) NewAccessToken(sessionID uuid.UUID) (string, error) {
	return s.jwtTm.NewAccessToken(sessionID, s.accessTokenTTL)
}

func (s *AuthService) NewRefreshToken() (string, error) {
	return s.jwtTm.NewRefreshToken(s.refreshTokenTTL)
}
