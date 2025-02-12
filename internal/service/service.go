package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/RVodassa/url-shortener/internal/lib/random"
	"github.com/RVodassa/url-shortener/internal/storage"
	"net/url"
)

type RandomProvider interface {
	RandomString(int) (string, error)
}

var (
	ErrNotFound = errors.New("ошибка: Url не найден")
	ErrBadUrl   = errors.New("ошибка: невалидный Url")
)

// TODO: в конфиг
const aliasLength = 10

type Service struct {
	Storage storage.Storage
	Random  RandomProvider
}

func New(storage storage.Storage) *Service {
	return &Service{
		Storage: storage,
		Random:  random.New(),
	}
}

// SaveUrl сохраняет Url и возвращает алиас.
func (s *Service) SaveUrl(ctx context.Context, urlStr string) (string, error) {
	const op = "service.SaveUrl"

	// Валидация Url
	parsedUrl, err := url.ParseRequestURI(urlStr)
	if err != nil || parsedUrl.Scheme == "" || parsedUrl.Host == "" {
		return "", ErrBadUrl
	}

	// Генерация алиаса
	var alias string

	for {
		alias, err = s.Random.RandomString(aliasLength)
		if err != nil {
			return "", fmt.Errorf("%s: %w", op, err)
		}

		err = s.Storage.SaveUrl(ctx, alias, urlStr)
		if err != nil {
			if errors.Is(err, storage.ErrExistAlias) {
				continue
			}
			return "", fmt.Errorf("%s: %w", op, err)
		}
		return alias, nil
	}
}

func (s *Service) GetUrl(ctx context.Context, alias string) (string, error) {
	const op = "service.GetUrl"

	getUrl, err := s.Storage.GetUrl(ctx, alias)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return "", ErrNotFound
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return getUrl, nil
}

func (s *Service) DeleteUrl(ctx context.Context, alias string) error {
	const op = "service.DeleteUrl"

	if err := s.Storage.DeleteUrl(ctx, alias); err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
