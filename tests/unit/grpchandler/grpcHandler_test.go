package grpchandler_test

import (
	"context"
	"errors"
	"github.com/RVodassa/url-shortener/internal/handler/grpc"
	"github.com/RVodassa/url-shortener/internal/service"
	mockService "github.com/RVodassa/url-shortener/internal/service/mock"
	"github.com/RVodassa/url-shortener/protos/genv1"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func TestGrpcHandler_SaveUrl(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServiceProvider := mockService.NewMockServiceProvider(ctrl)
	handler := grpchandler.New(mockServiceProvider)

	tests := []struct {
		name            string
		req             *genv1.SaveUrlRequest
		mockSaveUrl     func()
		expectedResp    *genv1.SaveUrlResponse
		expectedErr     error
		expectedErrCode codes.Code
	}{
		{
			name: "Успешное сохранение Url",
			req:  &genv1.SaveUrlRequest{Url: "https://example.com"},
			mockSaveUrl: func() {
				mockServiceProvider.EXPECT().
					SaveUrl(gomock.Any(), "https://example.com").
					Return("example-alias", nil)
			},
			expectedResp: &genv1.SaveUrlResponse{Alias: "example-alias"},
			expectedErr:  nil,
		},
		{
			name: "Пустой Url",
			req:  &genv1.SaveUrlRequest{Url: ""},
			mockSaveUrl: func() {
				// Нет вызова SaveUrl, так как валидация происходит до вызова сервиса
			},
			expectedErr:     status.Error(codes.InvalidArgument, grpchandler.ErrUrlEmpty.Error()),
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "Невалидный Url",
			req:  &genv1.SaveUrlRequest{Url: "invalid-Url"},
			mockSaveUrl: func() {
				mockServiceProvider.EXPECT().
					SaveUrl(gomock.Any(), "invalid-Url").
					Return("", service.ErrBadUrl)
			},
			expectedErr:     status.Error(codes.InvalidArgument, grpchandler.ErrBadUrl.Error()),
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "Внутренняя ошибка сервиса",
			req:  &genv1.SaveUrlRequest{Url: "https://example.com"},
			mockSaveUrl: func() {
				mockServiceProvider.EXPECT().
					SaveUrl(gomock.Any(), "https://example.com").
					Return("", errors.New("internal error"))
			},
			expectedErr:     status.Error(codes.Internal, grpchandler.ErrInternal.Error()),
			expectedErrCode: codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Настраиваем мок
			if tt.mockSaveUrl != nil {
				tt.mockSaveUrl()
			}

			resp, err := handler.SaveUrl(context.Background(), tt.req)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErrCode, status.Code(err))
				assert.Contains(t, err.Error(), tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResp, resp)
			}
		})
	}
}

func TestGrpcHandler_GetUrl(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServiceProvider := mockService.NewMockServiceProvider(ctrl)
	handler := grpchandler.New(mockServiceProvider)

	tests := []struct {
		name            string
		req             *genv1.GetUrlRequest
		mockGetUrl      func()
		expectedResp    *genv1.GetUrlResponse
		expectedErr     error
		expectedErrCode codes.Code
	}{
		{
			name: "Успешное получение Url",
			req:  &genv1.GetUrlRequest{Alias: "QWERTY1234"},
			mockGetUrl: func() {
				mockServiceProvider.EXPECT().
					GetUrl(gomock.Any(), "QWERTY1234").
					Return("https://example.com", nil)
			},
			expectedResp: &genv1.GetUrlResponse{Url: "https://example.com"},
			expectedErr:  nil,
		},
		{
			name: "Пустой alias",
			req:  &genv1.GetUrlRequest{Alias: ""},
			mockGetUrl: func() {
				// Нет вызова SaveUrl, так как валидация происходит до вызова сервиса
			},
			expectedErr:     status.Error(codes.InvalidArgument, grpchandler.ErrAliasEmpty.Error()),
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "Не найден url",
			req:  &genv1.GetUrlRequest{Alias: "QWERTY1234"},
			mockGetUrl: func() {
				mockServiceProvider.EXPECT().
					GetUrl(gomock.Any(), "QWERTY1234").
					Return("", service.ErrNotFound)
			},
			expectedErr:     status.Error(codes.NotFound, grpchandler.ErrNotFound.Error()),
			expectedErrCode: codes.NotFound,
		},
		{
			name: "Внутренняя ошибка сервиса",
			req:  &genv1.GetUrlRequest{Alias: "QWERTY1234"},
			mockGetUrl: func() {
				mockServiceProvider.EXPECT().
					GetUrl(gomock.Any(), "QWERTY1234").
					Return("", errors.New("internal error"))
			},
			expectedErr:     status.Error(codes.Internal, grpchandler.ErrInternal.Error()),
			expectedErrCode: codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockGetUrl != nil {
				tt.mockGetUrl()
			}
			resp, err := handler.GetUrl(context.Background(), tt.req)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErrCode, status.Code(err))
				assert.Contains(t, err.Error(), tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResp, resp)
			}
		})
	}
}
func TestGrpcHandler_DeleteUrl(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServiceProvider := mockService.NewMockServiceProvider(ctrl)
	handler := grpchandler.New(mockServiceProvider)

	tests := []struct {
		name            string
		req             *genv1.DeleteUrlRequest
		mockDeleteUrl   func()
		expectedResp    *genv1.DeleteUrlResponse
		expectedErr     error
		expectedErrCode codes.Code
	}{
		{
			name: "Успешное удаление url",
			req:  &genv1.DeleteUrlRequest{Alias: "QWERTY1234"},
			mockDeleteUrl: func() {
				mockServiceProvider.EXPECT().
					DeleteUrl(gomock.Any(), "QWERTY1234").
					Return(nil)
			},
			expectedResp: &genv1.DeleteUrlResponse{Status: "OK"},
			expectedErr:  nil,
		},
		{
			name: "Пустой alias",
			req:  &genv1.DeleteUrlRequest{Alias: ""},
			mockDeleteUrl: func() {
			},
			expectedErr:     status.Error(codes.InvalidArgument, grpchandler.ErrAliasEmpty.Error()),
			expectedErrCode: codes.InvalidArgument,
		},
		{
			name: "Не найден url",
			req:  &genv1.DeleteUrlRequest{Alias: "QWERTY1234"},
			mockDeleteUrl: func() {
				mockServiceProvider.EXPECT().
					DeleteUrl(gomock.Any(), "QWERTY1234").
					Return(service.ErrNotFound)
			},
			expectedErr:     status.Error(codes.NotFound, grpchandler.ErrNotFound.Error()),
			expectedErrCode: codes.NotFound,
		},
		{
			name: "Внутренняя ошибка сервиса",
			req:  &genv1.DeleteUrlRequest{Alias: "QWERTY1234"},
			mockDeleteUrl: func() {
				mockServiceProvider.EXPECT().
					DeleteUrl(gomock.Any(), "QWERTY1234").
					Return(errors.New("internal error"))
			},
			expectedErr:     status.Error(codes.Internal, grpchandler.ErrInternal.Error()),
			expectedErrCode: codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockDeleteUrl != nil {
				tt.mockDeleteUrl()
			}
			resp, err := handler.DeleteUrl(context.Background(), tt.req)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErrCode, status.Code(err))
				assert.Contains(t, err.Error(), tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResp, resp)
			}
		})
	}
}
