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

var (
	ErrUrlEmpty   = errors.New("ошибка: в запросе отсутствует url")
	ErrBadUrl     = errors.New("ошибка: недопустимая ссылка")
	ErrAliasEmpty = errors.New("ошибка: в запросе отсутствует alias")
	ErrNotFound   = errors.New("ошибка: url с таким alias не найден")
	ErrInternal   = errors.New("ошибка: внутренняя ошибка сервера")
)

type GrpcHandler struct {
	Service *service.Service
	genv1.UnimplementedUrlShortenerServer
}

func New(service *service.Service) *GrpcHandler {
	return &GrpcHandler{
		Service: service,
	}
}

func (g *GrpcHandler) SaveUrl(ctx context.Context, req *genv1.SaveUrlRequest) (*genv1.SaveUrlResponse, error) {
	const op = "grpchandler.SaveUrl"

	// Первичная валидация входных данных
	if req.Url == "" {
		log.Printf("%s: ошибка при сохранении url=%s: %v", op, req.Url, ErrUrlEmpty.Error())
		return nil, status.Error(codes.InvalidArgument, ErrUrlEmpty.Error())
	}

	// Вызов сервиса для сохранения URL
	alias, err := g.Service.SaveURL(ctx, req.Url)
	if err != nil {
		log.Printf("%s: ошибка при сохранении url=%s: %v", op, req.Url, err)

		if errors.Is(err, service.ErrBadUrl) {
			return nil, status.Error(codes.InvalidArgument, ErrBadUrl.Error())
		}
		return nil, status.Error(codes.Internal, ErrInternal.Error())
	}

	// Успешный ответ
	response := &genv1.SaveUrlResponse{
		Alias: alias,
	}

	log.Printf("%s: успешно сохранен URL с алиасом=%s", op, alias)
	return response, nil
}

func (g *GrpcHandler) GetUrl(ctx context.Context, req *genv1.GetUrlRequest) (*genv1.GetUrlResponse, error) {
	const op = "grpchandler.GetUrl"

	if req.Alias == "" {
		log.Printf("%s: ошибка при получении url с алиасом=%s: %v", op, req.Alias, ErrAliasEmpty)
		return nil, status.Error(codes.InvalidArgument, ErrAliasEmpty.Error())
	}

	url, err := g.Service.GetURL(ctx, req.Alias)
	if err != nil {
		log.Printf("%s: ошибка при получении url с алиасом=%s: %v", op, req.Alias, err)
		if errors.Is(err, service.ErrNotFound) {
			return nil, status.Error(codes.NotFound, ErrNotFound.Error())
		}
		return nil, status.Error(codes.Internal, ErrInternal.Error())
	}

	log.Printf("%s: успешно получен URL с алиасом=%s", op, req.Alias)
	return &genv1.GetUrlResponse{Url: url}, nil
}

func (g *GrpcHandler) DeleteUrl(ctx context.Context, req *genv1.DeleteUrlRequest) (*genv1.DeleteUrlResponse, error) {
	const op = "grpchandler.DeleteUrl"

	if req.Alias == "" {
		log.Printf("%s: ошибка при получении url с алиасом=%s: %v", op, req.Alias, ErrAliasEmpty)
		return nil, status.Error(codes.InvalidArgument, ErrAliasEmpty.Error())
	}

	err := g.Service.DeleteURL(ctx, req.Alias)
	if err != nil {
		log.Printf("%s: ошибка при удалении url с алиасом=%s: %v", op, req.Alias, err)
		if errors.Is(err, service.ErrNotFound) {
			return nil, status.Error(codes.NotFound, ErrNotFound.Error())
		}
		return nil, status.Error(codes.Internal, ErrInternal.Error())
	}

	response := &genv1.DeleteUrlResponse{
		Status: "success",
	}

	log.Printf("%s: успешно удален URL с алиасом=%s", op, req.Alias)
	return response, nil
}
