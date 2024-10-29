package redis

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"go-pocket-link/internal/domain"
	"go-pocket-link/pkg/cache/redis"
	"strings"
	"time"
)

type TokensRepository struct {
	cache *redis.DB
}

func NewTokensRepository(cache *redis.DB) *TokensRepository {
	return &TokensRepository{cache}
}

func (r *TokensRepository) Get(ctx context.Context, userID, tokenID uuid.UUID) (domain.Token, error) {
	token, err := r.cache.Get(ctx, fmt.Sprintf("%s:%s", userID, tokenID.String()))
	if err != nil {
		return domain.Token{}, err
	}

	return domain.Token{
		ID:           userID,
		UserID:       tokenID,
		RefreshToken: token,
	}, nil
}

func (r *TokensRepository) GetByKey(ctx context.Context, key string) (domain.Token, error) {
	token, err := r.cache.Get(ctx, key)
	if err != nil {
		return domain.Token{}, err
	}

	keyParts := strings.Split(key, ":")

	userID, err := uuid.Parse(keyParts[0])
	if err != nil {
		return domain.Token{}, err
	}

	tokenID, err := uuid.Parse(keyParts[1])
	if err != nil {
		return domain.Token{}, err
	}

	return domain.Token{
		ID:           tokenID,
		UserID:       userID,
		RefreshToken: token,
	}, nil
}

func (r *TokensRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Token, error) {
	keys, err := r.cache.ScanKeys(ctx, fmt.Sprintf("%s:*", userID.String()), 0)
	if err != nil {
		return nil, err
	}

	values, err := r.cache.ScanValues(ctx, keys)
	if err != nil {
		return nil, err
	}

	tokenValues := make([]string, 0, len(values))
	for _, value := range values {
		tokenValues = append(tokenValues, value.(string))
	}

	tokens := make([]domain.Token, 0, len(tokenValues))
	for i, value := range tokenValues {
		token := domain.Token{
			UserID:       userID,
			RefreshToken: value,
		}
		if token.ID, err = uuid.Parse(strings.Split(keys[i], ":")[1]); err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}
	return tokens, nil
}

func (r *TokensRepository) GetByTokenID(ctx context.Context, tokenID uuid.UUID) (domain.Token, error) {
	keys, err := r.cache.ScanKeys(ctx, fmt.Sprintf("*:%s", tokenID.String()), 1)
	if err != nil {
		return domain.Token{}, err
	}

	values, err := r.cache.ScanValues(ctx, keys)
	if err != nil {
		return domain.Token{}, err
	}

	tokenValue := values[0].(string)
	token := domain.Token{
		ID:           tokenID,
		RefreshToken: tokenValue,
	}
	if token.UserID, err = uuid.Parse(strings.Split(keys[0], ":")[0]); err != nil {
		return domain.Token{}, err
	}
	return token, nil
}

func (r *TokensRepository) Set(ctx context.Context, token *domain.Token, ttl time.Duration) error {
	return r.cache.Set(ctx, token.Key(), token.RefreshToken, ttl)
}

func (r *TokensRepository) Delete(ctx context.Context, key string) error {
	return r.cache.Delete(ctx, key)
}

func (r *TokensRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	keys, err := r.cache.ScanKeys(ctx, fmt.Sprintf("%s:*", userID.String()), 0)
	if err != nil {
		return err
	}

	return r.cache.Delete(ctx, keys...)
}

func (r *TokensRepository) DeleteByTokenID(ctx context.Context, tokenID uuid.UUID) error {
	keys, err := r.cache.ScanKeys(ctx, fmt.Sprintf("*:%s", tokenID.String()), 1)
	if err != nil {
		return err
	}

	return r.cache.Delete(ctx, keys[0])
}
