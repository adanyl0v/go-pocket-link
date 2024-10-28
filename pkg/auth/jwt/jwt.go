package jwt

import (
	"fmt"
	jwt5 "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

type TokenManager interface {
	ParseAccessToken(token string) (jwt5.Claims, error)
	ParseRefreshToken(token string) (jwt5.Claims, error)
	NewAccessToken(id uuid.UUID, ttl time.Duration) (string, error)
	NewRefreshToken(ttl time.Duration) (string, error)
}

const (
	ClaimsSubject   = "sub"
	ClaimsIssuedAt  = "iat"
	ClaimsExpiresAt = "exp"
)

type tokenManagerImpl struct {
	accessSecret  []byte
	refreshSecret []byte
}

func NewTokenManager(accessSecret, refreshSecret string) TokenManager {
	return &tokenManagerImpl{
		accessSecret:  []byte(accessSecret),
		refreshSecret: []byte(refreshSecret),
	}
}

func (tm *tokenManagerImpl) ParseAccessToken(token string) (jwt5.Claims, error) {
	return tm.parseToken(token, tm.accessSecret)
}

func (tm *tokenManagerImpl) ParseRefreshToken(token string) (jwt5.Claims, error) {
	return tm.parseToken(token, tm.refreshSecret)
}

func (tm *tokenManagerImpl) NewAccessToken(id uuid.UUID, ttl time.Duration) (string, error) {
	return tm.newToken(jwt5.MapClaims{
		ClaimsSubject:   id.String(),
		ClaimsIssuedAt:  time.Now().Unix(),
		ClaimsExpiresAt: time.Now().Add(ttl).Unix(),
	}, tm.accessSecret)
}

func (tm *tokenManagerImpl) NewRefreshToken(ttl time.Duration) (string, error) {
	return tm.newToken(jwt5.MapClaims{
		ClaimsIssuedAt:  time.Now().Unix(),
		ClaimsExpiresAt: time.Now().Add(ttl).Unix(),
	}, tm.refreshSecret)
}

func (tm *tokenManagerImpl) newToken(claims jwt5.Claims, secret []byte) (string, error) {
	token := jwt5.NewWithClaims(jwt5.SigningMethodHS256, claims)
	signed, err := token.SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("signing token: %w", err)
	}
	return signed, nil
}

func (tm *tokenManagerImpl) parseToken(token string, secret []byte) (jwt5.Claims, error) {
	parsed, err := jwt5.Parse(token, func(token *jwt5.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt5.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("parsing token: %w", err)
	} else if !parsed.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return parsed.Claims, nil
}
