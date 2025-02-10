package app

import (
	"context"
	"fmt"
	"github.com/RVodassa/url-shortener/internal/config"
	grpchandler "github.com/RVodassa/url-shortener/internal/handler/grpc"
	"github.com/RVodassa/url-shortener/internal/service"
	"github.com/RVodassa/url-shortener/protos/genv1"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
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
	store, err := initStore(ctx, a.cfg, a.StorageType)
	if err != nil {
		log.Fatalf("ошибка: хранилище не готово к работе: %v", err)
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
