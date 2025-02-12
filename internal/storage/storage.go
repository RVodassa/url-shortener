package storage

import (
	"context"
	"errors"
)

var (
	ErrUrlIsEmpty   = errors.New("ошибка: пустой url")
	ErrAliasIsEmpty = errors.New("ошибка: пустой alias")
	ErrNotFound     = errors.New("ошибка: url не найден")
	ErrExistAlias   = errors.New("ошибка: alias занят")
)

type Storage interface {
	SaveURL(ctx context.Context, alias, urlSave string) error
	GetUrl(ctx context.Context, alias string) (string, error)
	DeleteURL(ctx context.Context, alias string) error
	Disconnect(ctx context.Context) error
}
