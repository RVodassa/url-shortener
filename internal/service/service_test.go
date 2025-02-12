package service_test

import (
	"context"
	"fmt"
	mockRand "github.com/RVodassa/url-shortener/internal/lib/random/mock"
	"github.com/RVodassa/url-shortener/internal/service"
	"github.com/RVodassa/url-shortener/internal/storage"
	mockStore "github.com/RVodassa/url-shortener/internal/storage/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"testing"
)

// TODO: в конфиг
const aliasLength = 10

func TestService_SaveUrl(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mockStore.NewMockStorage(ctrl)
	mockRandom := mockRand.NewMockRandomProvider(ctrl)
	s := service.New(mockStorage, mockRandom)

	tests := []struct {
		name           string
		url            string
		mockRandom     func()
		mockSaveUrl    func()
		expectedResult string
		expectedErr    error
	}{
		{
			name: "успешное сохранение url",
			url:  "http://google.com",
			mockRandom: func() {
				mockRandom.EXPECT().
					RandomString(aliasLength).
					Return("example-alias", nil)
			},
			mockSaveUrl: func() {
				mockStorage.EXPECT().
					SaveUrl(gomock.Any(), "example-alias", "http://google.com").
					Return(nil)
			},
			expectedResult: "example-alias",
			expectedErr:    nil,
		},
		{
			name: "невалидный url",
			url:  "invalid-url",
			mockRandom: func() {
				// Нет вызова RandomString, так как валидация URL происходит раньше
			},
			mockSaveUrl: func() {
				// Нет вызова SaveUrl, так как валидация URL происходит раньше
			},
			expectedResult: "",
			expectedErr:    service.ErrBadUrl,
		},
		{
			name: "алиас уже существует",
			url:  "http://google.com",
			mockRandom: func() {
				// Первая попытка генерации алиаса
				mockRandom.EXPECT().
					RandomString(aliasLength).
					Return("existing-alias", nil)
				// Вторая попытка генерации алиаса
				mockRandom.EXPECT().
					RandomString(aliasLength).
					Return("new-alias", nil)
			},
			mockSaveUrl: func() {
				// Первая попытка сохранения (алиас уже существует)
				mockStorage.EXPECT().
					SaveUrl(gomock.Any(), "existing-alias", "http://google.com").
					Return(storage.ErrExistAlias)
				// Вторая попытка сохранения (успешно)
				mockStorage.EXPECT().
					SaveUrl(gomock.Any(), "new-alias", "http://google.com").
					Return(nil)
			},
			expectedResult: "new-alias",
			expectedErr:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Настраиваем моки
			if tt.mockRandom != nil {
				tt.mockRandom()
			}
			if tt.mockSaveUrl != nil {
				tt.mockSaveUrl()
			}

			// Вызываем метод SaveUrl
			result, err := s.SaveUrl(context.Background(), tt.url)

			// Проверяем результат
			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestService_GetUrl(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mockStore.NewMockStorage(ctrl)
	mockRandom := mockRand.NewMockRandomProvider(ctrl)
	s := service.New(mockStorage, mockRandom)

	tests := []struct {
		name        string
		alias       string
		mockGetUrl  func()
		expectedUrl string
		expectedErr error
	}{
		{
			name:  "успешное получение URL",
			alias: "example-alias",
			mockGetUrl: func() {
				mockStorage.EXPECT().
					GetUrl(gomock.Any(), "example-alias").
					Return("http://google.com", nil)
			},
			expectedUrl: "http://google.com",
			expectedErr: nil,
		},
		{
			name:  "alias не найден",
			alias: "not-exist-alias",
			mockGetUrl: func() {
				mockStorage.EXPECT().
					GetUrl(gomock.Any(), "not-exist-alias").
					Return("", storage.ErrNotFound)
			},
			expectedUrl: "",
			expectedErr: service.ErrNotFound,
		},
		{
			name:  "ошибка при получении URL",
			alias: "error-get-url",
			mockGetUrl: func() {
				mockStorage.EXPECT().
					GetUrl(gomock.Any(), "error-get-url").
					Return("", fmt.Errorf("internal errror"))
			},
			expectedUrl: "",
			expectedErr: fmt.Errorf("service.GetUrl: %w", fmt.Errorf("internal errror")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockGetUrl != nil {
				tt.mockGetUrl()
			}

			result, err := s.GetUrl(context.Background(), tt.alias)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedUrl, result)
		})
	}
}

func TestService_DeleteUrl(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mockStore.NewMockStorage(ctrl)
	mockRandom := mockRand.NewMockRandomProvider(ctrl)
	s := service.New(mockStorage, mockRandom)

	tests := []struct {
		name        string
		alias       string
		mockDelete  func()
		expectedErr error
	}{
		{
			name:  "успешное удаление URL",
			alias: "QWERTY1234",
			mockDelete: func() {
				mockStorage.EXPECT().
					DeleteUrl(gomock.Any(), "QWERTY1234").
					Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:  "alias не найден",
			alias: "not-exist-alias",
			mockDelete: func() {
				mockStorage.EXPECT().
					DeleteUrl(gomock.Any(), "not-exist-alias").
					Return(storage.ErrNotFound)
			},
			expectedErr: service.ErrNotFound,
		},
		{
			name:  "ошибка при удалении URL",
			alias: "error-alias",
			mockDelete: func() {
				mockStorage.EXPECT().
					DeleteUrl(gomock.Any(), "error-alias").
					Return(fmt.Errorf("internal errror"))
			},
			expectedErr: fmt.Errorf("service.DeleteUrl: %w", fmt.Errorf("internal errror")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Настраиваем моки
			if tt.mockDelete != nil {
				tt.mockDelete()
			}

			// Вызываем метод DeleteUrl
			err := s.DeleteUrl(context.Background(), tt.alias)

			// Проверяем результат
			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
