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

func (r *RedisStorage) SaveUrl(ctx context.Context, alias, UrlSave string) error {
	const op = "storage.RedisStorage.SaveUrl"

	if alias == "" {
		return storage.ErrAliasIsEmpty
	}

	if UrlSave == "" {
		return storage.ErrUrlIsEmpty
	}

	exists, err := r.client.Exists(ctx, alias).Result()
	if err != nil {
		return fmt.Errorf("%s: Url='%s', alias='%s'. %w", op, UrlSave, alias, err)
	}

	if exists > 0 {
		return storage.ErrExistAlias
	}

	cmd := r.client.Set(ctx, alias, UrlSave, 0)
	if cmd.Err() != nil {
		return fmt.Errorf("%s: Url='%s', alias='%s'. %w", op, UrlSave, alias, cmd.Err())
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
		return "", fmt.Errorf("%s: alias='%s'. %w", op, alias, cmd.Err())
	}

	return cmd.Val(), nil
}

func (r *RedisStorage) DeleteUrl(ctx context.Context, alias string) error {
	const op = "storage.RedisStorage.DeleteUrl"

	if alias == "" {
		return storage.ErrAliasIsEmpty
	}

	cmd := r.client.Del(ctx, alias)
	if cmd.Err() != nil {
		return fmt.Errorf("%s: alias='%s'. %w", op, alias, cmd.Err())
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
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
