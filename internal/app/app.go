package app

import (
	"log/slog"
	"time"

	"github.com/Gwinkamp/grpcauth-sso/internal/app/authapp"
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
	// TODO: инициализировать харнилище (storage)

	authApp := authapp.New(log, grpcPort)

	return &App{
		AuthApp: authApp,
	}
}
