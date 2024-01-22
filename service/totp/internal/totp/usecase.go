//go:generate mockgen -source usecase.go -destination mock/usecase_mock.go -package mock
package totp

import (
	"context"
	"github.com/GCFactory/dbo-system/service/totp/internal/models"
	"github.com/google/uuid"
)

type UseCase interface {
	Enroll(ctx context.Context, totp models.TOTPConfig) (*models.TOTPEnroll, error)
	Verify(ctx context.Context, url string) (*models.TOTPVerify, error)
	Validate(ctx context.Context, id uuid.UUID, code string) (*models.TOTPValidate, error)
	Enable(ctx context.Context, request *models.TOTPRequest) (*models.TOTPBase, error)
	Disable(ctx context.Context, request *models.TOTPRequest) (*models.TOTPBase, error)
}
