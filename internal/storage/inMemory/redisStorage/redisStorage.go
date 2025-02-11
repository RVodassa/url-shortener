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

func (r *RedisStorage) SaveURL(ctx context.Context, alias, urlSave string) error {
	const op = "storage.RedisStorage.SaveURL"

	if alias == "" {
		return storage.ErrAliasIsEmpty
	}

	if urlSave == "" {
		return storage.ErrUrlIsEmpty
	}

	exists, err := r.client.Exists(ctx, alias).Result()
	if err != nil {
		return fmt.Errorf("%s: %v: алиас=%s", op, err, alias)
	}

	if exists > 0 {
		return storage.ErrExistAlias
	}

	cmd := r.client.Set(ctx, alias, urlSave, 0)
	if cmd.Err() != nil {
		return fmt.Errorf("%s: %v: url=%s alias=%s", op, cmd.Err(), urlSave, alias)
	}

	return nil
}

func (r *RedisStorage) GetUrl(ctx context.Context, alias string) (string, error) {
	const op = "storage.RedisStorage.GetUrl"

	if alias == "" {
		return "", storage.ErrAliasIsEmpty
	}

	cmd := r.client.Get(ctx, alias)
	if cmd.Err() != nil {
		if errors.Is(cmd.Err(), redis.Nil) {
			return "", storage.ErrNotFound
		}
		return "", fmt.Errorf("%s: %v: алиас=%s", op, cmd.Err(), alias)
	}

	return cmd.Val(), nil
}

func (r *RedisStorage) DeleteURL(ctx context.Context, alias string) error {
	const op = "storage.RedisStorage.DeleteURL"

	if alias == "" {
		return storage.ErrAliasIsEmpty
	}

	cmd := r.client.Del(ctx, alias)
	if cmd.Err() != nil {
		return fmt.Errorf("%s: алиас=%s: %w", op, alias, cmd.Err())
	}

	if cmd.Val() == 0 {
		return storage.ErrNotFound
	}

	return nil
}

func (r *RedisStorage) Disconnect(ctx context.Context) error {
	const op = "storage.RedisStorage.Disconnect"

	err := r.client.Close()
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	return nil
}
