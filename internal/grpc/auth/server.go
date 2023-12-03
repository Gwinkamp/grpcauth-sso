package auth

import (
	"context"

	ssov1 "github.com/Gwinkamp/grpcauth-contracts/gen/go"
	"google.golang.org/grpc"
)

type AuthServerAPI struct {
	ssov1.UnimplementedAuthServer
}

func Register(gRPC *grpc.Server) {
	ssov1.RegisterAuthServer(gRPC, &AuthServerAPI{})
}

func (s *AuthServerAPI) Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	panic("Не реализовано")
}

func (s *AuthServerAPI) Register(ctx context.Context, req *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	panic("Не реализовано")
}

func (s *AuthServerAPI) IsAdmin(ctx context.Context, req *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {
	panic("Не реализовано")
}
