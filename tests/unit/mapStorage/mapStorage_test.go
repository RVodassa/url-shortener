package mapStorage_test

import (
	"context"
	"github.com/RVodassa/url-shortener/internal/storage"
	"github.com/RVodassa/url-shortener/internal/storage/inMemory/mapStorage"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapStorage_SaveUrl(t *testing.T) {
	mapStore := mapStorage.New()

	tests := []struct {
		name        string
		alias       string
		url         string
		expectedErr error
	}{
		{
			name:        "успешное сохранение URL",
			alias:       "example-alias",
			url:         "http://google.com",
			expectedErr: nil,
		},
		{
			name:        "пустой алиас",
			alias:       "",
			url:         "http://google.com",
			expectedErr: storage.ErrAliasIsEmpty,
		},
		{
			name:        "пустой URL",
			alias:       "example-alias",
			url:         "",
			expectedErr: storage.ErrUrlIsEmpty,
		},
		{
			name:        "алиас уже существует",
			alias:       "existing-alias",
			url:         "http://example.com",
			expectedErr: nil,
		},
		{
			name:        "сохранение существующего алиаса",
			alias:       "existing-alias",
			url:         "http://another-url.com",
			expectedErr: storage.ErrExistAlias,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mapStore.SaveUrl(context.Background(), tt.alias, tt.url)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestMapStorage_GetUrl(t *testing.T) {
	mapStore := mapStorage.New()
	err := mapStore.SaveUrl(context.Background(), "example-alias", "http://google.com")
	if err != nil {
		t.Errorf("error saving url %v", err)
		return
	}

	tests := []struct {
		name        string
		alias       string
		expectedUrl string
		expectedErr error
	}{
		{
			name:        "успешное получение URL",
			alias:       "example-alias",
			expectedUrl: "http://google.com",
			expectedErr: nil,
		},
		{
			name:        "алиас не найден",
			alias:       "nonexistent-alias",
			expectedUrl: "",
			expectedErr: storage.ErrNotFound,
		},
		{
			name:        "пустой алиас",
			alias:       "",
			expectedUrl: "",
			expectedErr: storage.ErrAliasIsEmpty,
		},
	}

	var url string
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url, err = mapStore.GetUrl(context.Background(), tt.alias)
			assert.Equal(t, tt.expectedUrl, url)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestMapStorage_DeleteUrl(t *testing.T) {
	mapStore := mapStorage.New()
	err := mapStore.SaveUrl(context.Background(), "example-alias", "http://google.com")
	if err != nil {
		t.Errorf("error saving url %v", err)
		return
	}

	tests := []struct {
		name        string
		alias       string
		expectedErr error
	}{
		{
			name:        "успешное удаление URL",
			alias:       "example-alias",
			expectedErr: nil,
		},
		{
			name:        "алиас не найден",
			alias:       "nonexistent-alias",
			expectedErr: storage.ErrNotFound,
		},
		{
			name:        "пустой алиас",
			alias:       "",
			expectedErr: storage.ErrAliasIsEmpty,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err = mapStore.DeleteUrl(context.Background(), tt.alias)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
