package main

import (
	"log/slog"
	"os"

	"github.com/Gwinkamp/grpcauth-sso/internal/config"
)

const (
	envDebug = "debug"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info(
		"Запуск приложения",
		slog.String("env", cfg.Env),
	)

	// TODO: инициализировать приложение (app)

	// TODO: запустить gRPC-сурвер приложения
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envDebug:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		panic("Отсутствует конфигурация логгера для окружения: " + env)
	}

	return log
}
