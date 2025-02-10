package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/RVodassa/url-shortener/internal/storage"
	"net/url"
)

var (
	ErrNotFound = errors.New("ошибка: url с таким alias не найден")
	ErrBadUrl   = errors.New("ошибка: недопустимая ссылка")
)

// TODO: перенести в конфиг
const aliasLength = 10

type Service struct {
	Storage storage.Storage
}

func New(s storage.Storage) *Service {
	return &Service{s}
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
	alias := NewRandomString(10)
	if alias == "" {
		return "", fmt.Errorf("%s: не удалось сгенерировать алиас", op)
	}

	err = s.Storage.SaveURL(ctx, alias, urlStr)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return alias, nil
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

func (s *Service) GetURL(ctx context.Context, alias string) (string, error) {
	const op = "service.DeleteURL"

	getUrl, err := s.Storage.GetUrl(ctx, alias)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return "", ErrNotFound
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return getUrl, nil
}
