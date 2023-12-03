package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Gwinkamp/grpcauth-sso/internal/app"
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
		"запуск приложения",
		slog.String("env", cfg.Env),
	)

	application := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)

	go application.AuthApp.MustRun()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.AuthApp.Stop()
	log.Info("приложение остановлено")
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
		panic("отсутствует конфигурация логгера для окружения: " + env)
	}

	return log
}
