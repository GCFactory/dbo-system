//go:generate mockgen -source pg_repository.go -destination mock/pg_repository_mock.go -package mock
package totp

import (
	"context"
	"github.com/GCFactory/dbo-system/service/totp/internal/models"
	"github.com/google/uuid"
)

type Repository interface {
	CreateConfig(ctx context.Context, totp models.TOTPConfig) error
	GetURL(ctx context.Context, id uuid.UUID) (string, error)
	UpdateActive(ctx context.Context, id uuid.UUID, status bool) error
}
