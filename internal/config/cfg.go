package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env        string `yaml:"env" env-required:"true"`
	HTTPServer `yaml:"http_server"`
}

type HTTPServer struct {
	Host        string        `yaml:"host" env-required:"true"`
	Port        string        `yaml:"port" env-required:"true"`
	Network     string        `yaml:"network" env-required:"true"`
	ReqTimeout  time.Duration `yaml:"request_timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func MustLoad(configPath string) *Config {
	_, err := os.Stat(configPath)
	if os.IsNotExist(err) {
		log.Fatalf("не найден конфиг файл. cfg_path=%s", configPath)
	}

	var cfg Config
	err = cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("ошибка при чтении файла конфиг.: %v", err)
	}
	return &cfg
}
