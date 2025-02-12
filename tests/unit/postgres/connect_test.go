package postgres_test

import (
	"github.com/RVodassa/url-shortener/internal/storage/sql/postgres"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue string
		setEnv       map[string]string
		want         string
	}{
		{
			name:         "Env variable is set",
			key:          "TEST_KEY",
			defaultValue: "default",
			setEnv:       map[string]string{"TEST_KEY": "value"},
			want:         "value",
		},
		{
			name:         "Env variable is not set",
			key:          "TEST_KEY",
			defaultValue: "default",
			setEnv:       map[string]string{},
			want:         "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Устанавливаем переменные окружения
			for k, v := range tt.setEnv {
				t.Setenv(k, v)
			}

			got := postgres.GetEnv(tt.key, tt.defaultValue)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetDBConfig(t *testing.T) {
	tests := []struct {
		name   string
		setEnv map[string]string
		want   map[string]string
	}{
		{
			name: "All env variables are set",
			setEnv: map[string]string{
				"DB_HOST":     "localhost",
				"DB_PORT":     "5432",
				"DB_USER":     "user",
				"DB_PASSWORD": "password",
				"DB_NAME":     "dbname",
				"DB_SSL":      "disable",
			},
			want: map[string]string{
				"host":     "localhost",
				"port":     "5432",
				"user":     "user",
				"password": "password",
				"name":     "dbname",
				"ssl":      "disable",
			},
		},
		{
			name:   "No env variables are set",
			setEnv: map[string]string{},
			want: map[string]string{
				"host":     "localhost",
				"port":     "5432",
				"user":     "",
				"password": "",
				"name":     "",
				"ssl":      "disable",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Устанавливаем переменные окружения
			for k, v := range tt.setEnv {
				t.Setenv(k, v)
			}

			got := postgres.GetDBConfig()
			assert.Equal(t, tt.want, got)
		})
	}
}
