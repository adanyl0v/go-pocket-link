package jwt

import (
	"fmt"
	jwt5 "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

type TokenManager interface {
	NewToken(id uuid.UUID, ttl time.Duration, secret Secret) (string, error)
	ParseToken(token string, secret Secret) (jwt5.MapClaims, error)
}

type Secret int

const (
	AccessSecret Secret = iota
	RefreshSecret
)

const (
	ClaimsJwtID     = "jti"
	ClaimsIssuer    = "iss"
	ClaimsSubject   = "sub"
	ClaimsAudience  = "aud"
	ClaimsIssuedAt  = "iat"
	ClaimsExpiresAt = "exp"
)

type StaticClaims struct {
	Issuer   string `json:"iss"`
	Audience string `json:"aud"`
}

type tokenManagerImpl struct {
	accessSecret  []byte
	refreshSecret []byte
	staticClaims  StaticClaims
}

func NewTokenManager(accessSecret, refreshSecret string, claims StaticClaims) TokenManager {
	return &tokenManagerImpl{
		accessSecret:  []byte(accessSecret),
		refreshSecret: []byte(refreshSecret),
		staticClaims:  claims,
	}
}

func (tm *tokenManagerImpl) ParseToken(token string, secret Secret) (jwt5.MapClaims, error) {
	var s []byte
	switch secret {
	case AccessSecret:
		s = tm.accessSecret
	case RefreshSecret:
		s = tm.refreshSecret
	default:
		return nil, fmt.Errorf("invalid secret")
	}
	return tm.parseToken(token, s)
}

func (tm *tokenManagerImpl) NewToken(id uuid.UUID, ttl time.Duration, secret Secret) (string, error) {
	var s []byte
	switch secret {
	case AccessSecret:
		s = tm.accessSecret
	case RefreshSecret:
		s = tm.refreshSecret
	default:
		return "", fmt.Errorf("invalid secret")
	}
	return tm.newToken(id, ttl, s)
}

func (tm *tokenManagerImpl) newToken(id uuid.UUID, ttl time.Duration, secret []byte) (string, error) {
	token := jwt5.NewWithClaims(jwt5.SigningMethodHS256, jwt5.MapClaims{
		ClaimsJwtID:     uuid.New(),
		ClaimsIssuer:    tm.staticClaims.Issuer,
		ClaimsAudience:  tm.staticClaims.Audience,
		ClaimsSubject:   id.String(),
		ClaimsIssuedAt:  time.Now().Unix(),
		ClaimsExpiresAt: time.Now().Add(ttl).Unix(),
	})

	signed, err := token.SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("%w (signing token)", err)
	}

	return signed, nil
}

func (tm *tokenManagerImpl) parseToken(token string, secret []byte) (jwt5.MapClaims, error) {
	parsed, err := jwt5.Parse(token, func(token *jwt5.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt5.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("%w (parsing token)", err)
	} else if !parsed.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	mapClaims := parsed.Claims.(jwt5.MapClaims)
	if mapClaims[ClaimsIssuer] != tm.staticClaims.Issuer {
		return nil, fmt.Errorf("invalid %s claims", ClaimsIssuer)
	} else if mapClaims[ClaimsAudience] != tm.staticClaims.Audience {
		return nil, fmt.Errorf("invalid %s claims", ClaimsAudience)
	}

	return mapClaims, nil
}
