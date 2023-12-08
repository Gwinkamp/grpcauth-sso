package auth

import (
	"context"
	"errors"
	"log/slog"

	ssov1 "github.com/Gwinkamp/grpcauth-contracts/gen/go"
	"github.com/Gwinkamp/grpcauth-sso/internal/domain/app"
	"github.com/Gwinkamp/grpcauth-sso/internal/services/auth"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServerAPI struct {
	ssov1.UnimplementedAuthServer
	validate *validator.Validate
	auth     app.Auth
	log      *slog.Logger
}

func NewAuthServer(auth app.Auth, logger *slog.Logger) *AuthServerAPI {
	validate := validator.New(validator.WithRequiredStructEnabled())
	registerAuthServerValidationRules(validate)

	return &AuthServerAPI{
		validate: validate,
		auth:     auth,
		log:      logger,
	}
}

func Register(gRPC *grpc.Server, auth app.Auth, logger *slog.Logger) {
	ssov1.RegisterAuthServer(gRPC, NewAuthServer(auth, logger))
}

func (s *AuthServerAPI) Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	if err := s.validate.Struct(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	token, err := s.auth.Login(
		ctx,
		req.GetEmail(),
		req.GetPassword(),
		uuid.MustParse(req.GetServiceId()),
	)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "Неверный логин или пароль")
		}
		if errors.Is(err, auth.ErrServiceNotFound) {
			return nil, status.Error(codes.InvalidArgument, "Сервис не найден")
		}

		s.log.Error(
			"внутренняя ошибка авторизации",
			slog.String("email", req.GetEmail()),
			slog.String("serviceId", req.GetServiceId()),
		)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.LoginResponse{AccessToken: token}, nil
}

func (s *AuthServerAPI) Register(ctx context.Context, req *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	if err := s.validate.Struct(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	userId, err := s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, auth.ErrUserAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, "пользователь уже зарегистрирован")
		}

		s.log.Error(
			"внутренняя ошибка регистрации нового пользователя",
			slog.String("email", req.GetEmail()),
		)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.RegisterResponse{UserId: userId.String()}, nil
}

func (s *AuthServerAPI) IsAdmin(ctx context.Context, req *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {
	if err := s.validate.Struct(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	isAdmin, err := s.auth.IsAdmin(ctx, uuid.MustParse(req.GetUserId()))
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.NotFound, "пользователь не найден")
		}

		s.log.Error(
			"внутренняя идентификации администратора",
			slog.String("email", req.GetUserId()),
		)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.IsAdminResponse{IsAdmin: isAdmin}, nil
}
