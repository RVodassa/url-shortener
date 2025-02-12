package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/RVodassa/url-shortener/internal/storage"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresDB struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *PostgresDB {
	return &PostgresDB{pool: pool}
}

// SaveURL сохраняет URL в базе данных.
func (p *PostgresDB) SaveURL(ctx context.Context, alias, urlSave string) error {
	const op = "storage.PostgresDB.SaveURL"

	if alias == "" {
		return storage.ErrAliasIsEmpty
	}
	if urlSave == "" {
		return storage.ErrUrlIsEmpty
	}

	query := `INSERT INTO urls (alias, url) VALUES ($1, $2)`

	_, err := p.pool.Exec(ctx, query, alias, urlSave)
	if err != nil {
		// Проверка на ошибку уникальности
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" { // Код ошибки для нарушения уникальности
				return storage.ErrExistAlias
			}
		}
		return fmt.Errorf("%s: url='%s', alias='%s'. %w", op, urlSave, alias, err)
	}

	return nil
}

// GetUrl возвращает URL по его alias.
func (p *PostgresDB) GetUrl(ctx context.Context, alias string) (string, error) {
	const op = "storage.PostgresDB.GetUrl"

	if alias == "" {
		return "", storage.ErrAliasIsEmpty
	}

	var url string
	query := `SELECT url FROM urls WHERE alias = $1`

	err := p.pool.QueryRow(ctx, query, alias).Scan(&url)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", storage.ErrNotFound
		}
		return "", fmt.Errorf("%s: alias='%s'. %w", op, alias, err)
	}

	return url, nil
}

// DeleteURL удаляет URL по его alias.
func (p *PostgresDB) DeleteURL(ctx context.Context, alias string) error {
	const op = "storage.PostgresDB.DeleteURL"

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
func (p *PostgresDB) Disconnect(ctx context.Context) error {
	p.pool.Close()
	return nil
}
