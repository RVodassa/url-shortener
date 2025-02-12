package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/RVodassa/url-shortener/internal/storage"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type IPGX interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Close()
}
type Postgres struct {
	pool IPGX
}

func New(pool IPGX) *Postgres {
	return &Postgres{pool: pool}
}

// SaveUrl сохраняет Url в базе данных.
func (p *Postgres) SaveUrl(ctx context.Context, alias, urlsave string) error {
	const op = "storage.Postgres.SaveUrl"

	if alias == "" {
		return storage.ErrAliasIsEmpty
	}
	if urlsave == "" {
		return storage.ErrUrlIsEmpty
	}

	query := `INSERT INTO urls (alias, Url) VALUES ($1, $2)`

	_, err := p.pool.Exec(ctx, query, alias, urlsave)
	if err != nil {
		// Проверка на ошибку уникальности
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" { // Код ошибки для нарушения уникальности
				return storage.ErrExistAlias
			}
		}
		return fmt.Errorf("%s: url='%s', alias='%s'. %w", op, urlsave, alias, err)
	}

	return nil
}

// GetUrl возвращает Url по его alias.
func (p *Postgres) GetUrl(ctx context.Context, alias string) (string, error) {
	const op = "storage.Postgres.GetUrl"

	if alias == "" {
		return "", storage.ErrAliasIsEmpty
	}

	var Url string
	query := `SELECT Url FROM urls WHERE alias = $1`

	err := p.pool.QueryRow(ctx, query, alias).Scan(&Url)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", storage.ErrNotFound
		}
		return "", fmt.Errorf("%s: alias='%s'. %w", op, alias, err)
	}

	return Url, nil
}

// DeleteUrl удаляет Url по его alias.
func (p *Postgres) DeleteUrl(ctx context.Context, alias string) error {
	const op = "storage.Postgres.DeleteUrl"

	if alias == "" {
		return storage.ErrAliasIsEmpty
	}

	query := `DELETE FROM urls WHERE alias = $1`

	result, err := p.pool.Exec(ctx, query, alias)
	if err != nil {
		return fmt.Errorf("%s: alias='%s'. %w", op, alias, err)
	}

	// Проверяем количество затронутых строк
	if result.RowsAffected() == 0 {
		return storage.ErrNotFound
	}

	return nil
}

// Disconnect закрывает соединение с базой данных.
func (p *Postgres) Disconnect(ctx context.Context) error {
	p.pool.Close()
	return nil
}
