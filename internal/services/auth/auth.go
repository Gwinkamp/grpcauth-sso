package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/Gwinkamp/grpcauth-sso/internal/domain/providers"
	"github.com/Gwinkamp/grpcauth-sso/internal/lib/jwt"
	"github.com/Gwinkamp/grpcauth-sso/internal/storage"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("неправильный логин или пароль")
	ErrServiceNotFound    = errors.New("сервис не найден")
	ErrUserAlreadyExists  = errors.New("пользователь уже существует")
)

type Auth struct {
	log             *slog.Logger
	userProvider    providers.UserProvider
	serviceProvider providers.ServiceProvider
	tokenTTL        time.Duration
}

func New(
	log *slog.Logger,
	userProvider providers.UserProvider,
	serviceProvider providers.ServiceProvider,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		log:             log,
		userProvider:    userProvider,
		serviceProvider: serviceProvider,
		tokenTTL:        tokenTTL,
	}
}

// Login проверяет данные пользователя в системе, авторизует его и выдает токен доступа
func (a *Auth) Login(
	ctx context.Context,
	email string,
	passwod string,
	serviceId uuid.UUID,
) (string, error) {
	const operation = "Auth.Login"

	log := a.log.With(
		slog.String("operation", operation),
		slog.String("email", email),
		slog.String("service_id", serviceId.String()),
	)

	user, err := a.userProvider.GetUser(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("не найден пользователь", slog.String("error", err.Error()))
			return "", fmt.Errorf("%s: %w", operation, ErrInvalidCredentials)
		}

		return "", fmt.Errorf("%s: %w", operation, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(passwod)); err != nil {
		log.Warn("неправильный пароль", slog.String("error", err.Error()))
		return "", fmt.Errorf("%s: %w", operation, ErrInvalidCredentials)
	}

	service, err := a.serviceProvider.GetService(ctx, serviceId)
	if err != nil {
		if errors.Is(err, storage.ErrServiceNotFound) {
			log.Warn("не найден сервис", slog.String("error", err.Error()))
			return "", fmt.Errorf("%s: %w", operation, ErrServiceNotFound)
		}
		return "", fmt.Errorf("%s: %w", operation, err)
	}

	token, err := jwt.NewToken(user, service, a.tokenTTL)
	if err != nil {
		log.Error("не удалось сгенерировать токен", slog.String("error", err.Error()))
		return "", fmt.Errorf("%s: %w", operation, err)
	}

	return token, nil
}

// RegisterNewUser проверяет данные нового пользователя и создает его
func (a *Auth) RegisterNewUser(
	ctx context.Context,
	email string,
	password string,
) (uuid.UUID, error) {
	const operation = "Auth.RegisterNewUser"

	log := a.log.With(
		slog.String("operation", operation),
		slog.String("email", email),
	)

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("ошибка генерации хэша", slog.String("error", err.Error()))
		return uuid.Nil, fmt.Errorf("%s: %w", operation, err)
	}

	id, err := a.userProvider.CreateUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserAlreadyExists) {
			log.Warn("пользователь уже существует", slog.String("error", err.Error()))
			return uuid.Nil, fmt.Errorf("%s: %w", operation, ErrUserAlreadyExists)
		}
		log.Error("ошибка сохранения пользователя в хранилище", slog.String("error", err.Error()))
		return uuid.Nil, fmt.Errorf("%s: %w", operation, err)
	}

	log.Info("зарегистрирован новый пользователь", slog.String("user_id", id.String()))
	return id, nil
}

// IsAdmin определяет, является ли пользователь администратором
func (a *Auth) IsAdmin(
	ctx context.Context,
	userId uuid.UUID,
) (bool, error) {
	const operation = "Auth.IsAdmin"

	isAdmin, err := a.userProvider.IsAdmin(ctx, userId)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return false, fmt.Errorf("%s: %w", operation, ErrInvalidCredentials)
		}
		return false, fmt.Errorf("%s: %w", operation, err)
	}

	return isAdmin, nil
}
