package providers

import (
	"context"

	"github.com/Gwinkamp/grpcauth-sso/internal/domain/models"
	"github.com/google/uuid"
)

type UserProvider interface {
	GetUser(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userId uuid.UUID) (bool, error)
	CreateUser(ctx context.Context, email string, passHash []byte) (id uuid.UUID, err error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}
