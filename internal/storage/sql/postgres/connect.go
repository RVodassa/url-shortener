package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
)

// GetDBConfig возвращает карту с переменными окружения для подключения к базе данных.
func GetDBConfig() map[string]string {
	return map[string]string{
		"host":     GetEnv("DB_HOST", "localhost"),
		"port":     GetEnv("DB_PORT", "5432"),
		"user":     GetEnv("DB_USER", ""),
		"password": GetEnv("DB_PASSWORD", ""),
		"name":     GetEnv("DB_NAME", ""),
		"ssl":      GetEnv("DB_SSL", "disable"),
	}
}

// GetEnv возвращает значение переменной окружения или значение по умолчанию.
func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func ConnectDB(ctx context.Context) (*pgxpool.Pool, error) {
	const op = "postgres.ConnectDB"

	// Получаем конфигурацию базы данных
	config := GetDBConfig()

	// Проверяем обязательные переменные
	if config["user"] == "" {
		return nil, fmt.Errorf("%s: пустой DB_USER", op)
	}
	if config["password"] == "" {
		return nil, fmt.Errorf("%s: пустой DB_PASSWORD", op)
	}
	if config["name"] == "" {
		return nil, fmt.Errorf("%s: пустой DB_NAME", op)
	}

	// Формируем строку подключения
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		config["user"], config["password"], config["host"], config["port"], config["name"], config["ssl"])

	conn, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("%s: New. %w", op, err)
	}

	// Проверяем соединение
	err = conn.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: Ping. %w", op, err)
	}

	log.Printf("%s: запуск миграций", op)
	err = runMigrations(connStr)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	log.Printf("%s: база данных готова к работе", op)

	return conn, nil
}

func runMigrations(connStr string) error {
	const op = "postgres.runMigrations"

	m, err := migrate.New("file://migrations", connStr)
	if err != nil {
		return fmt.Errorf("%s: New. %v", op, err)
	}
	defer func() {
		if m != nil {
			errSource, errDB := m.Close()
			if errSource != nil {
				log.Printf("%s: m.Close. %v", op, errSource)
			}
			if errDB != nil {
				log.Printf("%s: m.Close. %v", op, errDB)
			}
			return
		}
	}()

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("%s: не удалось применить миграции. Ошибка: %v", op, err)
	}

	return nil
}
