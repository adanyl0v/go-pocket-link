package cache

import (
	"context"
	"time"
)

type DB interface {
	Close() error
	Get(ctx context.Context, key string) (string, error)
	ScanKeys(ctx context.Context, regexp string, count int64) ([]string, error)
	ScanValues(ctx context.Context, keys []string) ([]any, error)
	Set(ctx context.Context, key string, value any, expiration time.Duration) error
	Delete(ctx context.Context, keys ...string) error
}
