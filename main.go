package main

import (
	"github.com/RVodassa/url-shortener/app"
	"github.com/RVodassa/url-shortener/internal/config"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
		return
	}
	configPath := os.Getenv("CFG_PATH")
	storageType := os.Getenv("STORAGE_TYPE")

	// загрузка конфиг.
	cfg := config.MustLoad(configPath)
	if cfg == nil {
		log.Fatal("ошибка: конфиг. не готов к работе")
	}

	// запуск
	newApp := app.New(cfg, storageType)
	newApp.Run()

}
