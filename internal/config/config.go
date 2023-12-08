package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string        `yaml:"env" env-default:"debug"`
	StoragePath string        `yaml:"storage_path" env-required:"true"`
	TokenTTL    time.Duration `yaml:"token_ttl" env-default:"1h"`
	GRPC        GRPCConfig    `yaml:"grpc"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("Не задан путь до файла с конфигурацией")
	}

	return MustLoadByPath(path)
}

func MustLoadByPath(pathToConfig string) *Config {
	if _, err := os.Stat(pathToConfig); os.IsNotExist(err) {
		panic("Не найден файл конфигурации: " + pathToConfig)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(pathToConfig, &cfg); err != nil {
		panic("Ошибка чтения файла конфигурации: " + err.Error())
	}

	return &cfg
}

// fetchConfigPath получает путь до файла с конфигурацией из командной строки или из переменной окружения
// Приоритет: flag > env > default
// default: пустая строка
func fetchConfigPath() string {
	var res string

	// --config="path/to/config.yaml"
	flag.StringVar(&res, "config", "", "Путь до файла с конфигурацией")
	flag.Parse()

	if res == "" {
		os.Getenv("CONFIG_PATH")
	}

	return res
}
