package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/RVodassa/url-shortener/internal/config"
	grpchandler "github.com/RVodassa/url-shortener/internal/handler/grpc"
	"github.com/RVodassa/url-shortener/internal/service"
	"github.com/RVodassa/url-shortener/internal/storage"
	"github.com/RVodassa/url-shortener/internal/storage/inMemory/redisStorage"
	"github.com/RVodassa/url-shortener/protos/genv1"
	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Доступные хранилища
const (
	InMemoryStorage = "in-memory"
	SqlStorage      = "sql"
)

type App struct {
	cfg         *config.Config
	StorageType string
}

func New(cfg *config.Config, storageType string) *App {
	return &App{
		cfg:         cfg,
		StorageType: storageType,
	}
}

func (a *App) Run() {
	// контекст для управления жизненным циклом
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// инициализация хранилища
	store, err := initStorage(ctx, a.cfg, a.StorageType)
	if err != nil {
		log.Printf("ошибка: хранилище не готово к работе: %v", err)
		return
	}

	// инстанс сервиса
	newService := service.New(store)

	// инстанс handler
	newHandler := grpchandler.New(newService)

	// запуск gRPC сервера
	newGrpcServer := grpc.NewServer()
	genv1.RegisterUrlShortenerServer(newGrpcServer, newHandler)

	// Установка слушателя на порту
	lis, err := net.Listen(a.cfg.Network, fmt.Sprintf(":%s", a.cfg.Port))
	if err != nil {
		log.Fatalf("ошибка при создании слушателя: %v", err)
	}

	// Канал для обработки сигналов
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// Запуск сервера в горутине
	go func() {
		log.Printf("Запуск gRPC сервера [%s]\n", lis.Addr().String())
		if err = newGrpcServer.Serve(lis); err != nil {
			log.Fatalf("ошибка при запуске сервера: %v", err)
		}
	}()

	// Ожидание сигнала завершения
	<-signalChan
	log.Println("Получен сигнал завершения, начинаем мягкое завершение работы...")

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	newGrpcServer.GracefulStop() // Остановка сервера

	// TODO: мягкое завершение работы остальных частей приложения
}

func initStorage(ctx context.Context, cfg *config.Config, storageType string) (storage.Storage, error) {
	log.Printf("инициализация хранилища типа: %s", storageType)

	var store storage.Storage

	switch storageType {
	case InMemoryStorage:
		// Redis для in-memory storage
		redisClient := redis.NewClient(&redis.Options{
			Addr: cfg.Redis.Address,
		})
		if err := redisClient.Ping(ctx).Err(); err != nil {
			return nil, fmt.Errorf("ошибка при подключении к Redis: %v", err)
		}
		store = redisStorage.New(redisClient)
		log.Printf("хранилище доступно по адресу: %s\n", cfg.Redis.Address)

	case SqlStorage:
		// Postgres для sql storage
		return nil, errors.New("postgres не реализован")
	default:
		return nil, errors.New("неизвестный тип хранилища")
	}

	return store, nil
}
