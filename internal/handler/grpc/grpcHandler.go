package grpchandler

import (
	"context"
	"errors"
	"github.com/RVodassa/url-shortener/internal/service"
	"github.com/RVodassa/url-shortener/protos/genv1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

//go:generate mockgen -source=grpcHandler.go -destination=./../../service/mock/service_mock.go
type ServiceProvider interface {
	SaveUrl(ctx context.Context, UrlStr string) (string, error)
	GetUrl(ctx context.Context, alias string) (string, error)
	DeleteUrl(ctx context.Context, alias string) error
}

var (
	ErrUrlEmpty   = errors.New("ошибка: пустой url")
	ErrBadUrl     = errors.New("ошибка: невалидный url")
	ErrAliasEmpty = errors.New("ошибка: пустой alias")
	ErrNotFound   = errors.New("ошибка: url не найден")
	ErrInternal   = errors.New("ошибка: внутренняя ошибка")
)

type GrpcHandler struct {
	Service ServiceProvider
	genv1.UnimplementedUrlShortenerServer
}

func New(service ServiceProvider) *GrpcHandler {
	return &GrpcHandler{
		Service: service,
	}
}

func (g *GrpcHandler) SaveUrl(ctx context.Context, req *genv1.SaveUrlRequest) (*genv1.SaveUrlResponse, error) {
	const op = "grpchandler.SaveUrl"

	// Первичная валидация входных данных
	if req.Url == "" {
		log.Printf("%s: url='%s'. %v", op, req.Url, ErrUrlEmpty.Error())
		return nil, status.Error(codes.InvalidArgument, ErrUrlEmpty.Error())
	}

	// Вызов сервиса для сохранения Url
	alias, err := g.Service.SaveUrl(ctx, req.Url)
	if err != nil {
		log.Printf("%s: url='%s'. %v", op, req.Url, err)

		if errors.Is(err, service.ErrBadUrl) {
			return nil, status.Error(codes.InvalidArgument, ErrBadUrl.Error())
		}
		return nil, status.Error(codes.Internal, ErrInternal.Error())
	}

	// Успешный ответ
	response := &genv1.SaveUrlResponse{
		Alias: alias,
	}

	log.Printf("%s: alias='%s'. сохранен Url", op, alias)
	return response, nil
}

func (g *GrpcHandler) GetUrl(ctx context.Context, req *genv1.GetUrlRequest) (*genv1.GetUrlResponse, error) {
	const op = "grpchandler.GetUrl"

	if req.Alias == "" {
		log.Printf("%s: alias='%s'. %v", op, req.Alias, ErrAliasEmpty)
		return nil, status.Error(codes.InvalidArgument, ErrAliasEmpty.Error())
	}

	Url, err := g.Service.GetUrl(ctx, req.Alias)
	if err != nil {
		log.Printf("%s: alias='%s'. %v", op, req.Alias, err)
		if errors.Is(err, service.ErrNotFound) {
			return nil, status.Error(codes.NotFound, ErrNotFound.Error())
		}
		return nil, status.Error(codes.Internal, ErrInternal.Error())
	}

	log.Printf("%s: alias='%s'. получен Url", op, req.Alias)
	return &genv1.GetUrlResponse{Url: Url}, nil
}

func (g *GrpcHandler) DeleteUrl(ctx context.Context, req *genv1.DeleteUrlRequest) (*genv1.DeleteUrlResponse, error) {
	const op = "grpchandler.DeleteUrl"

	if req.Alias == "" {
		log.Printf("%s: alias='%s'. %v", op, req.Alias, ErrAliasEmpty)
		return nil, status.Error(codes.InvalidArgument, ErrAliasEmpty.Error())
	}

	err := g.Service.DeleteUrl(ctx, req.Alias)
	if err != nil {
		log.Printf("%s: alias='%s'. %v", op, req.Alias, err)
		if errors.Is(err, service.ErrNotFound) {
			return nil, status.Error(codes.NotFound, ErrNotFound.Error())
		}
		return nil, status.Error(codes.Internal, ErrInternal.Error())
	}

	response := &genv1.DeleteUrlResponse{
		Status: "OK", // eng
	}

	log.Printf("%s: alias='%s'. удален Url", op, req.Alias)
	return response, nil
}
