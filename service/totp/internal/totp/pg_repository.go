//go:generate mockgen -source pg_repository.go -destination mock/pg_repository_mock.go -package mock
package totp

import (
	"context"
	"github.com/GCFactory/dbo-system/service/totp/internal/models"
)

type Repository interface {
	CreateConfig(ctx context.Context, totp models.TOTPConfig) error
	Verify(ctx context.Context, request models.TOTPRequest) (*models.TOTPConfig, error)
	Validate(ctx context.Context, request models.TOTPRequest) (*models.TOTPBase, error)
	Enable(ctx context.Context, request models.TOTPRequest) (*models.TOTPBase, error)
	Disable(ctx context.Context, request models.TOTPRequest) (*models.TOTPBase, error)
}
