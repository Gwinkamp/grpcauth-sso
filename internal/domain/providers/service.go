package providers

import (
	"context"

	"github.com/Gwinkamp/grpcauth-sso/internal/domain/models"
	"github.com/google/uuid"
)

type ServiceProvider interface {
	GetService(ctx context.Context, serviceId uuid.UUID) (models.Service, error)
}
