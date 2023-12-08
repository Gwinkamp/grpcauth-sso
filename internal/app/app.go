package app

import (
	"log/slog"
	"time"

	"github.com/Gwinkamp/grpcauth-sso/internal/app/authapp"
	"github.com/Gwinkamp/grpcauth-sso/internal/services/auth"
	"github.com/Gwinkamp/grpcauth-sso/internal/storage/sqlite"
)

type App struct {
	AuthApp *authapp.AuthApp
}

func New(
	log *slog.Logger,
	grpcPort int,
	storagePath string,
	tokenTTL time.Duration,
) *App {
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}

	authService := auth.New(log, storage, storage, tokenTTL)

	grpcApp := authapp.New(log, authService, grpcPort)

	return &App{
		AuthApp: grpcApp,
	}
}
