package jwt

import (
	jwt5 "github.com/golang-jwt/jwt/v5"
	"go-pocket-link/pkg/errb"
	"time"
)

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type TokenManager interface {
	// CreateAccessToken with "sub" claims used to store a user ID
	CreateAccessToken(sub string, ttl time.Duration) (string, error)
	// CreateRefreshToken with "sub" claims used to store a user ID
	CreateRefreshToken(sub string, ttl time.Duration) (string, error)
	// Parse access token and return "sub" claim
	Parse(accessToken string) (string, error)
}

type tokenManagerImpl struct {
	secret []byte
}

func NewTokenManager(secret []byte) TokenManager {
	return &tokenManagerImpl{secret: secret}
}

const (
	claimSubject   = "sub"
	claimExpiresAt = "exp"
)

func (m *tokenManagerImpl) CreateAccessToken(sub string, ttl time.Duration) (string, error) {
	tok, err := m.newToken(jwt5.MapClaims{
		claimSubject:   sub,
		claimExpiresAt: time.Now().Add(ttl).Unix(),
	})
	if err != nil {
		return "", errb.Errorf("create access token: %v", err)
	}
	return tok, nil
}

func (m *tokenManagerImpl) CreateRefreshToken(sub string, ttl time.Duration) (string, error) {
	tok, err := m.newToken(jwt5.MapClaims{
		claimSubject:   sub,
		claimExpiresAt: time.Now().Add(ttl).Unix(),
	})
	if err != nil {
		return "", errb.Errorf("create refresh token: %v", err)
	}
	return tok, nil
}

func (m *tokenManagerImpl) Parse(accessToken string) (string, error) {
	token, err := jwt5.Parse(accessToken, func(token *jwt5.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt5.SigningMethodHMAC); !ok {
			return nil, errb.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.secret, nil
	})
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", errb.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt5.MapClaims)
	if !ok {
		return "", errb.Errorf("unexpected claims: %v", token.Claims)
	}
	return claims[claimSubject].(string), nil
}

func (m *tokenManagerImpl) newToken(claims jwt5.Claims) (string, error) {
	token := jwt5.NewWithClaims(jwt5.SigningMethodHS256, claims)
	signed, err := token.SignedString(m.secret)
	if err != nil {
		return "", err
	}
	return signed, nil
}
