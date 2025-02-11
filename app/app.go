package app

import (
	"context"
	"fmt"
	"github.com/RVodassa/url-shortener/internal/config"
	grpchandler "github.com/RVodassa/url-shortener/internal/handler/grpc"
	"github.com/RVodassa/url-shortener/internal/service"
	"github.com/RVodassa/url-shortener/internal/storage"
	"github.com/RVodassa/url-shortener/internal/storage/inMemory/redisStorage"
	"github.com/RVodassa/url-shortener/internal/storage/sql/postgres"
	"github.com/RVodassa/url-shortener/protos/genv1"
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
	Redis    = "redis"
	InMemory = "in-memory"
	Postgres = "postgres"
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
	store, err := NewStorage(ctx)
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
	defer func(lis net.Listener) {
		err = lis.Close()
		if err != nil {
			log.Printf("ошибка при отложенном вызове listener.close: %v", err)
		}
	}(lis)

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

	newGrpcServer.GracefulStop()
	// TODO: мягкое завершение работы остальных частей приложения
}

func NewStorage(ctx context.Context) (storage.Storage, error) {

	storageType := os.Getenv("STORAGE_TYPE")
	log.Printf("инициализация хранилища типа: %s", storageType)

	var store storage.Storage
	var err error

	switch storageType {

	case Redis:
		store, err = redisStorage.Connect(ctx)
		if err != nil {
			return nil, err
		}
		return store, nil

	case Postgres:
		conn, errConn := postgres.ConnectDB(ctx)
		if errConn != nil {
			log.Fatalf("ошибка при подключении к PostgreSQL: %v", errConn)
		}
		store = postgres.New(conn)
		return store, nil

	default:
		return nil, fmt.Errorf("неизвестный тип хранилища: %s", storageType)
	}

}
