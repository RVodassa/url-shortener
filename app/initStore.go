package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/RVodassa/url-shortener/internal/config"
	"github.com/RVodassa/url-shortener/internal/storage"
	"github.com/RVodassa/url-shortener/internal/storage/inMemory/redisStorage"
	"github.com/go-redis/redis/v8"
	"log"
)

// Доступные хранилища
const (
	InMemoryStorage = "in-memory"
	sqlStorage      = "sql"
)

func initStore(ctx context.Context, cfg *config.Config, storageType string) (storage.Storage, error) {
	log.Printf("инициализация хранилища типа: %s", storageType)

	var store storage.Storage

	switch storageType {
	case InMemoryStorage:
		// Redis для in-memory storage
		redisClient := redis.NewClient(&redis.Options{
			Addr: cfg.Redis.Address,
		})
		if err := redisClient.Ping(ctx).Err(); err != nil {
			return nil, fmt.Errorf("ошибка подключения к Redis: %v", err)
		}
		store = redisStorage.New(redisClient)
		log.Printf("хранилище доступно по адресу: %s\n", cfg.Redis.Address)

	case sqlStorage:
		// Postgres для sql storage
		return nil, errors.New("postgres не реализован")
	default:
		return nil, errors.New("неизвестный тип хранилища")
	}

	return store, nil
}
