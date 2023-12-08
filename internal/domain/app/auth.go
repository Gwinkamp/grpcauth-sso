package app

import (
	"context"

	"github.com/google/uuid"
)

type Auth interface {
	Login(ctx context.Context, email, password string, serviceId uuid.UUID) (token string, err error)
	RegisterNewUser(ctx context.Context, email, password string) (userId uuid.UUID, err error)
	IsAdmin(ctx context.Context, userId uuid.UUID) (bool, error)
}
