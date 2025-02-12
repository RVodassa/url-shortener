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
	const op = "app.Run"

	// контекст для управления жизненным циклом
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// сtrl+c для мягкого завершения работы
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// инициализация модулей
	store, err := NewStorage(ctx) // хранилище
	if err != nil {
		log.Printf("%s: %v", op, err)
		os.Exit(1)
	}

	newService := service.New(store)          // сервис
	newHandler := grpchandler.New(newService) // handler

	// Слушатель на порту
	lis, err := net.Listen(a.cfg.Network, a.cfg.Port)
	if err != nil {
		log.Fatalf("%s: установка слушателя. Ошибка: %v", op, err)
	}

	defer func(lis net.Listener) {
		err = lis.Close()
		if err != nil {
			log.Printf("%s: закрытие слушателя. Ошибка: %v", op, err)
		}
	}(lis)

	// запуск gRPC сервера
	if a.cfg.Network == "" || a.cfg.Port == "" {
		log.Fatalf("%s: пустой port='%s' или network='%s'", op, a.cfg.Network, a.cfg.Port)
	}

	newGrpcServer := grpc.NewServer()
	genv1.RegisterUrlShortenerServer(newGrpcServer, newHandler)

	go func() {
		log.Printf("%s: запуск gRPC сервера. port=[%s], network=[%s]\n", op, a.cfg.Port, a.cfg.Network)
		if err = newGrpcServer.Serve(lis); err != nil {
			log.Printf("%s: запуск сервера. Ошибка: %v", op, err)
			signalChan <- syscall.SIGTERM
		}
	}()

	// ожидает сигнал завершения работы
	<-signalChan
	log.Printf("%s: завершение работы...", op)

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	newGrpcServer.GracefulStop()

	// TODO: мягкое завершение работы остальных частей приложения
}

func NewStorage(ctx context.Context) (storage.Storage, error) {
	const op = "app.NewStorage"

	storageType := os.Getenv("STORAGE_TYPE")
	log.Printf("%s: storageType='%s'", op, storageType)

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
			return nil, fmt.Errorf("%s: storageType='%s'. Ошибка: %v", op, storageType, errConn)
		}
		store = postgres.New(conn)
		return store, nil

	default:
		return nil, fmt.Errorf("%s: storageType='%s'. Ошибка: неизвестный тип: %s", op, storageType, err)
	}
}
