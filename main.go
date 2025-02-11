package main

import (
	"flag"
	"github.com/RVodassa/url-shortener/app"
	"github.com/RVodassa/url-shortener/internal/config"
	"github.com/joho/godotenv"
	"log"
)

// TODO: postgres, tests, map for in-memory

func main() {
	// Параметры запуска
	var storageType string // тип хранилища (default: inMemoryStorage)
	flag.StringVar(&storageType, "storage", app.Redis, "postgres or default redis")

	var configPath string // путь к файлу конф.
	flag.StringVar(&configPath, "cfg_path", "", "path to config file")

	flag.Parse()

	// загрузка конфиг.
	cfg := config.MustLoad(configPath)
	if cfg == nil {
		log.Fatal("ошибка: конфиг. не готов к работе")
	}
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
		return
	}

	// новый инстанс приложения и запуск
	newApp := app.New(cfg, storageType)
	newApp.Run()
}
