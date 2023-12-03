package authapp

import (
	"fmt"
	"log/slog"
	"net"

	authgrpc "github.com/Gwinkamp/grpcauth-sso/internal/grpc/auth"
	"google.golang.org/grpc"
)

type AuthApp struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

// New здает новый gRPC сервер приложения AuthServer
func New(log *slog.Logger, port int) *AuthApp {
	gRPCServer := grpc.NewServer()
	authgrpc.Register(gRPCServer)

	return &AuthApp{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *AuthApp) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *AuthApp) Run() error {
	const operation = "authapp.Run"

	log := a.log.With(
		slog.String("operation", operation),
		slog.Int("port", a.port),
	)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}

	log.Info(
		"gRPC сервер запущен",
		slog.String("addr", listener.Addr().String()),
	)

	if err := a.gRPCServer.Serve(listener); err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}

	return nil
}

func (a *AuthApp) Stop() {
	const operation = "authapp.Stop"

	a.gRPCServer.GracefulStop()

	a.log.With(slog.String("operation", operation)).Info("gRPC сервер остановлен")
}
