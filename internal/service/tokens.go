package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/adanyl0v/go-pocket-link/internal/domain"
	"github.com/adanyl0v/go-pocket-link/internal/repository"
	"github.com/adanyl0v/go-pocket-link/pkg/auth/jwt"
	jwt5 "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

type TokensService struct {
	repo            repository.TokensRepository
	jwtTm           jwt.TokenManager
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

func NewTokensService(repo repository.TokensRepository, jwtTm jwt.TokenManager, accessTokenTTL, refreshTokenTTL time.Duration) *TokensService {
	return &TokensService{
		repo:            repo,
		jwtTm:           jwtTm,
		AccessTokenTTL:  accessTokenTTL,
		RefreshTokenTTL: refreshTokenTTL,
	}
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (s *TokensService) NewTokenPair(userID uuid.UUID) (TokenPair, error) {
	accessToken, err := s.jwtTm.NewToken(userID, s.AccessTokenTTL, jwt.AccessSecret)
	if err != nil {
		return TokenPair{}, fmt.Errorf("%w (generating access token)", err)
	}

	refreshToken, err := s.jwtTm.NewToken(userID, s.RefreshTokenTTL, jwt.RefreshSecret)
	if err != nil {
		return TokenPair{}, fmt.Errorf("%w (generating refresh token)", err)
	}

	return TokenPair{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (s *TokensService) SaveRefreshToken(ctx context.Context, token *domain.Token) error {
	return s.repo.Set(ctx, token, s.RefreshTokenTTL)
}

func (s *TokensService) SaveRefreshTokenFromString(ctx context.Context, token string) error {
	parsed, err := s.parseToken(token, jwt.RefreshSecret)
	if err != nil {
		return err
	}
	return s.repo.Set(ctx, &parsed, s.RefreshTokenTTL)
}

func (s *TokensService) ParseAccessToken(token string) (domain.Token, error) {
	return s.parseToken(token, jwt.AccessSecret)
}

func (s *TokensService) ParseRefreshToken(token string) (domain.Token, error) {
	return s.parseToken(token, jwt.RefreshSecret)
}

func (s *TokensService) ValidateAccessToken(token string) (jwt5.MapClaims, error) {
	claims, err := s.jwtTm.ParseToken(token, jwt.AccessSecret)
	if err != nil {
		return nil, fmt.Errorf("%w (validating access token)", err)
	}

	return claims, nil
}

func (s *TokensService) ValidateRefreshToken(token string) (jwt5.MapClaims, error) {
	claims, err := s.jwtTm.ParseToken(token, jwt.RefreshSecret)
	if err != nil {
		return nil, fmt.Errorf("%w (validating refresh token)", err)
	}

	return claims, nil
}

func (s *TokensService) InvalidateRefreshToken(ctx context.Context, tokenID uuid.UUID) error {
	if err := s.repo.DeleteByTokenID(ctx, tokenID); err != nil {
		return fmt.Errorf("%w (invalidating refresh token)", err)
	}
	return nil
}

func (s *TokensService) InvalidateUser(ctx context.Context, userID uuid.UUID) error {
	if err := s.repo.DeleteByUserID(ctx, userID); err != nil {
		return fmt.Errorf("%w (invalidating user)", err)
	}
	return nil
}

func (s *TokensService) parseToken(token string, secret jwt.Secret) (domain.Token, error) {
	var claims, err = jwt5.MapClaims(nil), error(nil)
	switch secret {
	case jwt.AccessSecret:
		claims, err = s.jwtTm.ParseToken(token, jwt.AccessSecret)
	case jwt.RefreshSecret:
		claims, err = s.jwtTm.ParseToken(token, jwt.RefreshSecret)
	default:
		return domain.Token{}, errors.New("invalid secret")
	}
	if err != nil {
		return domain.Token{}, err
	}

	tokenIDClaims, ok := claims[jwt.ClaimsJwtID]
	if !ok {
		return domain.Token{}, errMissedClaims(jwt.ClaimsJwtID)
	}

	tokenID, err := uuid.Parse(tokenIDClaims.(string))
	if err != nil {
		return domain.Token{}, errParsingClaims(jwt.ClaimsJwtID, err)
	}

	userIDClaims, ok := claims[jwt.ClaimsSubject]
	if !ok {
		return domain.Token{}, errMissedClaims(jwt.ClaimsSubject)
	}

	userID, err := uuid.Parse(userIDClaims.(string))
	if err != nil {
		return domain.Token{}, errParsingClaims(jwt.ClaimsSubject, err)
	}

	return domain.Token{
		ID:           tokenID,
		UserID:       userID,
		RefreshToken: token,
	}, nil
}

func errMissedClaims(key string) error {
	return fmt.Errorf("missed %s claims", key)
}

func errParsingClaims(key string, err error) error {
	return fmt.Errorf("%w (parsing %s claims)", err, key)
}
