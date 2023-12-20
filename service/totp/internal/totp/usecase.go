//go:generate mockgen -source usecase.go -destination mock/usecase_mock.go -package mock
package totp

import (
	"context"
	"github.com/GCFactory/dbo-system/service/totp/internal/models"
)

type UseCase interface {
	Enroll(ctx context.Context, totp models.TOTPConfig) (*models.TOTPEnroll, error)
	Verify(ctx context.Context, request *models.TOTPRequest) (*models.TOTPConfig, error)
	Validate(ctx context.Context, request *models.TOTPRequest) (*models.TOTPBase, error)
	Enable(ctx context.Context, request *models.TOTPRequest) (*models.TOTPBase, error)
	Disable(ctx context.Context, request *models.TOTPRequest) (*models.TOTPBase, error)
}
