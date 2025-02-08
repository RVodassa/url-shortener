package main

import (
	"flag"
	"log"
)

// go run cmd/url-shortener/main.go -storage postgres -cfg-path ./path/to/config.yaml
func main() {
	const op = "main.main"

	var storageType string
	var configPath string

	flag.StringVar(&storageType, "storage", "in-memory", "postgres or default in-memory")
	flag.StringVar(&configPath, "cfg-path", "./config/cfg.yaml", "path to config file")
	flag.Parse()

	log.Printf("%s тип хранилища: %s", op, storageType)
	log.Printf("%s путь к файлу конфигурации: %s", op, configPath)

	//cfg := config.MustLoad(configPath)
	//fmt.Println(cfg)
}
