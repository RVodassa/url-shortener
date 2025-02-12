package mapStorage

import (
	"context"
	"github.com/RVodassa/url-shortener/internal/storage"
	"sync"
)

type MapStorage struct {
	mu    sync.RWMutex
	store map[string]string
}

func New() storage.Storage {
	return &MapStorage{
		store: make(map[string]string),
	}
}

func (s *MapStorage) SaveURL(ctx context.Context, alias, urlSave string) error {
	const op = "storage.MapStorage.SaveURL"

	if alias == "" {
		return storage.ErrAliasIsEmpty
	}
	if urlSave == "" {
		return storage.ErrUrlIsEmpty
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.store[alias]; exists {
		return storage.ErrExistAlias
	}

	s.store[alias] = urlSave
	return nil
}

func (s *MapStorage) GetUrl(ctx context.Context, alias string) (string, error) {
	const op = "storage.MapStorage.GetUrl"

	if alias == "" {
		return "", storage.ErrAliasIsEmpty
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	url, exists := s.store[alias]
	if !exists {
		return "", storage.ErrNotFound
	}

	return url, nil
}

func (s *MapStorage) DeleteURL(ctx context.Context, alias string) error {
	const op = "storage.MapStorage.DeleteURL"

	if alias == "" {
		return storage.ErrAliasIsEmpty
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.store[alias]; !exists {
		return storage.ErrNotFound
	}

	delete(s.store, alias)
	return nil
}

func (s *MapStorage) Disconnect(ctx context.Context) error {
	return nil
}
