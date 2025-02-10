package redisStorage

import (
	"context"
	"errors"
	"fmt"
	"github.com/RVodassa/url-shortener/internal/storage"
	"github.com/go-redis/redis/v8"
)

type RedisStorage struct {
	client *redis.Client
}

func New(client *redis.Client) *RedisStorage {
	return &RedisStorage{client: client}
}

func (r *RedisStorage) SaveURL(ctx context.Context, alias, urlSave string) error {
	const op = "storage.RedisStorage.SaveURL"

	exists, err := r.client.Exists(ctx, alias).Result()
	if err != nil {
		return err
	}

	if exists > 0 {
		return fmt.Errorf("%s: %s aliace=%s", op, storage.ErrExistAlias, alias)
	}

	cmd := r.client.Set(ctx, alias, urlSave, 0)
	if cmd.Err() != nil {
		return cmd.Err()
	}

	return nil
}

func (r *RedisStorage) GetUrl(ctx context.Context, alias string) (string, error) {
	const op = "storage.RedisStorage.GetUrl"

	cmd := r.client.Get(ctx, alias)
	if cmd.Err() != nil {
		if errors.Is(cmd.Err(), redis.Nil) {
			return "", storage.ErrNotFound
		}
		return "", fmt.Errorf("%s: failed to get URL: %w", op, cmd.Err())
	}

	return cmd.Val(), nil
}

func (r *RedisStorage) DeleteURL(ctx context.Context, alias string) error {
	const op = "storage.RedisStorage.DeleteURL"

	if alias == "" {
		return fmt.Errorf("%s: пустой алиас, alias=%s", op, alias)
	}

	cmd := r.client.Del(ctx, alias)
	if cmd.Err() != nil {
		return fmt.Errorf("%s: alias=%s: %w", op, alias, cmd.Err())
	}
	
	if cmd.Val() == 0 {
		return storage.ErrNotFound
	}

	return nil
}
