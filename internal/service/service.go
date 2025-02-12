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
	ErrNotFound = errors.New("ошибка: url не найден")
	ErrBadUrl   = errors.New("ошибка: неправильный url")
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

// SaveURL сохраняет URL и возвращает алиас.
func (s *Service) SaveURL(ctx context.Context, urlStr string) (string, error) {
	const op = "service.SaveURL"

	// Валидация URL
	parsedURL, err := url.ParseRequestURI(urlStr)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return "", ErrBadUrl
	}

	// Генерация алиаса
	var alias string

	for {
		alias, err = s.Random.RandomString(aliasLength)
		if err != nil {
			return "", fmt.Errorf("%s: %w", op, err)
		}

		err = s.Storage.SaveURL(ctx, alias, urlStr)
		if err != nil {
			if errors.Is(err, storage.ErrExistAlias) {
				continue
			}
			return "", fmt.Errorf("%s: %w", op, err)
		}
		return alias, nil
	}
}

func (s *Service) GetURL(ctx context.Context, alias string) (string, error) {
	const op = "service.GetURL"

	getUrl, err := s.Storage.GetUrl(ctx, alias)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return "", ErrNotFound
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return getUrl, nil
}

func (s *Service) DeleteURL(ctx context.Context, alias string) error {
	const op = "service.DeleteURL"

	if err := s.Storage.DeleteURL(ctx, alias); err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
