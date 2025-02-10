package storage

import (
	"context"
	"errors"
)

var (
	ErrNotFound   = errors.New("ошибка: url с таким alias не найден")
	ErrExistAlias = errors.New("ошибка: url с таким alias уже существует")
)

type Storage interface {
	SaveURL(ctx context.Context, alias, urlSave string) error
	GetUrl(ctx context.Context, alias string) (string, error)
	DeleteURL(ctx context.Context, alias string) error
}
