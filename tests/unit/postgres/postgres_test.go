package postgres_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/RVodassa/url-shortener/internal/storage"
	"github.com/jackc/pgx/v5"
	"testing"

	"github.com/RVodassa/url-shortener/internal/storage/sql/postgres"
	mockPGX "github.com/RVodassa/url-shortener/internal/storage/sql/postgres/mock"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

//go:generate mockgen -source=postgres_test.go -destination=./mock/pgx_mock.go

func TestSaveUrl(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pgxmock := mockPGX.NewMockIPGX(ctrl)
	store := postgres.New(pgxmock)

	tests := []struct {
		name    string
		alias   string
		url     string
		mock    func()
		wantErr error
	}{
		{
			name:  "Success",
			alias: "alias1",
			url:   "http://example.com",
			mock: func() {
				pgxmock.EXPECT().
					Exec(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(pgconn.NewCommandTag("INSERT 1"), nil)
			},
			wantErr: nil,
		},
		{
			name:    "Empty Alias",
			alias:   "",
			url:     "http://example.com",
			mock:    func() {},
			wantErr: storage.ErrAliasIsEmpty,
		},
		{
			name:    "Empty Url",
			alias:   "alias1",
			url:     "",
			mock:    func() {},
			wantErr: storage.ErrUrlIsEmpty,
		},
		{
			name:  "Duplicate Alias",
			alias: "alias1",
			url:   "http://example.com",
			mock: func() {
				pgxmock.EXPECT().
					Exec(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(pgconn.CommandTag{}, &pgconn.PgError{Code: "23505"})
			},
			wantErr: storage.ErrExistAlias,
		},
		{
			name:  "Internal Error",
			alias: "alias1",
			url:   "http://example.com",
			mock: func() {
				pgxmock.EXPECT().
					Exec(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(pgconn.CommandTag{}, errors.New("internal error"))
			},
			wantErr: fmt.Errorf("storage.Postgres.SaveUrl: url='http://example.com', alias='alias1'. internal error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := store.SaveUrl(context.Background(), tt.alias, tt.url)

			if tt.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.wantErr.Error())
			}
		})
	}
}
func TestGetUrl(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pgxmock := mockPGX.NewMockIPGX(ctrl)
	store := postgres.New(pgxmock)

	tests := []struct {
		name    string
		alias   string
		mock    func()
		want    string
		wantErr error
	}{
		{
			name:  "Success",
			alias: "alias1",
			mock: func() {
				pgxmock.EXPECT().QueryRow(gomock.Any(), gomock.Any(), gomock.Any()).Return(pgxmock)
				pgxmock.EXPECT().Scan(gomock.Any()).SetArg(0, "http://example.com").Return(nil)
			},
			want:    "http://example.com",
			wantErr: nil,
		},
		{
			name:    "Empty Alias",
			alias:   "",
			mock:    func() {},
			want:    "",
			wantErr: storage.ErrAliasIsEmpty,
		},
		{
			name:  "Not Found",
			alias: "alias1",
			mock: func() {
				pgxmock.EXPECT().QueryRow(gomock.Any(), gomock.Any(), gomock.Any()).Return(pgxmock)
				pgxmock.EXPECT().Scan(gomock.Any()).Return(pgx.ErrNoRows)
			},
			want:    "",
			wantErr: storage.ErrNotFound,
		},
		{
			name:  "Internal Error",
			alias: "alias1",
			mock: func() {
				pgxmock.EXPECT().QueryRow(gomock.Any(), gomock.Any(), gomock.Any()).Return(pgxmock)
				pgxmock.EXPECT().Scan(gomock.Any()).Return(errors.New("internal error"))
			},
			want:    "",
			wantErr: fmt.Errorf("storage.Postgres.GetUrl: alias='alias1'. internal error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := store.GetUrl(context.Background(), tt.alias)

			if tt.wantErr == nil {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			} else {
				assert.EqualError(t, err, tt.wantErr.Error())
			}
		})
	}
}
func TestDeleteUrl(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pgxmock := mockPGX.NewMockIPGX(ctrl)
	store := postgres.New(pgxmock)

	tests := []struct {
		name    string
		alias   string
		mock    func()
		wantErr error
	}{
		{
			name:  "Success",
			alias: "alias1",
			mock: func() {
				pgxmock.EXPECT().Exec(gomock.Any(), gomock.Any(), gomock.Any()).Return(pgconn.NewCommandTag("DELETE 1"), nil)
			},
			wantErr: nil,
		},
		{
			name:    "Empty Alias",
			alias:   "",
			mock:    func() {},
			wantErr: storage.ErrAliasIsEmpty,
		},
		{
			name:  "Not Found",
			alias: "alias1",
			mock: func() {
				pgxmock.EXPECT().Exec(gomock.Any(), gomock.Any(), gomock.Any()).Return(pgconn.NewCommandTag("DELETE 0"), nil)
			},
			wantErr: storage.ErrNotFound,
		},
		{
			name:  "Internal Error",
			alias: "alias1",
			mock: func() {
				pgxmock.EXPECT().Exec(gomock.Any(), gomock.Any(), gomock.Any()).Return(pgconn.CommandTag{}, errors.New("internal error"))
			},
			wantErr: fmt.Errorf("storage.Postgres.DeleteUrl: alias='alias1'. internal error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := store.DeleteUrl(context.Background(), tt.alias)

			if tt.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.wantErr.Error())
			}
		})
	}
}

func TestDisconnect(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pgxmock := mockPGX.NewMockIPGX(ctrl)
	store := postgres.New(pgxmock)

	tests := []struct {
		name    string
		mock    func()
		wantErr error
	}{
		{
			name: "Success",
			mock: func() {
				pgxmock.EXPECT().Close().Return()
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := store.Disconnect(context.Background())

			if tt.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.wantErr.Error())
			}
		})
	}
}
