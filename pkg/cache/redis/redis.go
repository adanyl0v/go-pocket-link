package redis

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"time"
)

type DB struct {
	client *redis.Client
}

func Connect(dsn string) (*DB, error) {
	opt, err := redis.ParseURL(dsn)
	if err != nil {
		return nil, errConnecting(err)
	}
	client := redis.NewClient(opt)
	if _, err = client.Ping(context.Background()).Result(); err != nil {
		return nil, errConnecting(err)
	}
	return &DB{client: client}, nil
}

func (c *DB) Close() error {
	err := c.client.Close()
	if err != nil {
		return errClosing(err)
	}
	return nil
}

func (c *DB) Get(ctx context.Context, key string) (string, error) {
	val, err := c.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", errKeyDoesNotExist(key)
	} else if err != nil {
		return "", err
	}
	return val, nil
}

func (c *DB) ScanKeys(ctx context.Context, regexp string, count int64) ([]string, error) {
	keys := make([]string, 0, count)
	iter := c.client.Scan(ctx, 0, regexp, count).Iterator()
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	if err := iter.Err(); err != nil {
		return nil, err
	}
	return keys, nil
}

func (c *DB) ScanValues(ctx context.Context, keys []string) ([]any, error) {
	values, err := c.client.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}
	return values, nil
}

func (c *DB) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	return c.client.Set(ctx, key, value, expiration).Err()
}

func (c *DB) Delete(ctx context.Context, keys ...string) error {
	return c.client.Del(ctx, keys...).Err()
}
