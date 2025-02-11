package redisStorage

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"os"
)

func Connect(ctx context.Context) (*RedisStorage, error) {

	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		return nil, fmt.Errorf("ошибка: пустой REDIS_ADDR в переменной окр")
	}

	r := &RedisStorage{
		client: redis.NewClient(&redis.Options{
			Addr: addr,
		}),
	}
	if err := r.client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("ошибка при подключении к Redis: %v", err)
	}

	return r, nil
}
