package auth

import (
	"context"

	ssov1 "github.com/Gwinkamp/grpcauth-contracts/gen/go"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServerAPI struct {
	ssov1.UnimplementedAuthServer
	validate *validator.Validate
}

func NewAuthServer() *AuthServerAPI {
	validate := validator.New(validator.WithRequiredStructEnabled())
	registerAuthServerValidationRules(validate)

	return &AuthServerAPI{
		validate: validate,
	}
}

func Register(gRPC *grpc.Server) {
	ssov1.RegisterAuthServer(gRPC, NewAuthServer())
}

func (s *AuthServerAPI) Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	if err := s.validate.Struct(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &ssov1.LoginResponse{AccessToken: "dummy access token"}, nil
}

func (s *AuthServerAPI) Register(ctx context.Context, req *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	if err := s.validate.Struct(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &ssov1.RegisterResponse{UserId: "bb493e9b-ee5b-4c28-818c-657bcf099460"}, nil
}

func (s *AuthServerAPI) IsAdmin(ctx context.Context, req *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {
	if err := s.validate.Struct(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &ssov1.IsAdminResponse{IsAdmin: false}, nil
}
