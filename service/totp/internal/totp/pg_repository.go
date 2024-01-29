//go:generate mockgen -source pg_repository.go -destination mock/pg_repository_mock.go -package mock
package totp

import (
	"context"
	"github.com/GCFactory/dbo-system/service/totp/internal/models"
	"github.com/google/uuid"
)

type Repository interface {
	CreateConfig(ctx context.Context, totp models.TOTPConfig) error
	GetActiveConfig(ctx context.Context, userId uuid.UUID) (*models.TOTPConfig, error)
	GetConfigByTotpId(ctx context.Context, totpId uuid.UUID) (*models.TOTPConfig, error)
	GetConfigByUserId(ctx context.Context, userId uuid.UUID) (*models.TOTPConfig, error)
	GetLastDisabledConfig(ctx context.Context, userId uuid.UUID) (*models.TOTPConfig, error)
	UpdateTotpActivityByTotpId(ctx context.Context, totpId uuid.UUID, status bool) error
}
