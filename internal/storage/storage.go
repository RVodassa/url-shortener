package storage

import (
	"context"
	"errors"
)

var (
	ErrUrlIsEmpty   = errors.New("ошибка: пустой Url")
	ErrAliasIsEmpty = errors.New("ошибка: пустой alias")
	ErrNotFound     = errors.New("ошибка: Url не найден")
	ErrExistAlias   = errors.New("ошибка: alias занят")
)

type Storage interface {
	SaveUrl(ctx context.Context, alias, UrlSave string) error
	GetUrl(ctx context.Context, alias string) (string, error)
	DeleteUrl(ctx context.Context, alias string) error
	Disconnect(ctx context.Context) error
}
