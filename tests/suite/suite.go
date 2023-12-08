package suite

import (
	"context"
	"net"
	"strconv"
	"testing"

	ssov1 "github.com/Gwinkamp/grpcauth-contracts/gen/go"
	"github.com/Gwinkamp/grpcauth-sso/internal/config"
	"github.com/Gwinkamp/grpcauth-sso/internal/storage/sqlite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Suite struct {
	*testing.T
	Cfg        *config.Config
	AuthClient ssov1.AuthClient
	Storage    *sqlite.Storage
}

const (
	grpcHost = "localhost"
)

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg := config.MustLoadByPath("./config/local_tests.yaml")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)
	t.Cleanup(func() {
		t.Helper()
		cancel()
	})

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		t.Fatalf("ошибка инициализации хранилища: %v", err)
	}

	cc, err := grpc.DialContext(
		ctx,
		grpcAddress(cfg),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("ошибка соединения с gRPC сервером: %v", err)
	}

	return ctx, &Suite{
		T:          t,
		Cfg:        cfg,
		Storage:    storage,
		AuthClient: ssov1.NewAuthClient(cc),
	}
}

func grpcAddress(cfg *config.Config) string {
	return net.JoinHostPort(grpcHost, strconv.Itoa(cfg.GRPC.Port))
}
