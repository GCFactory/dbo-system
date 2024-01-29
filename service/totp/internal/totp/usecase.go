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
	Enable(ctx context.Context, id uuid.UUID) (*models.TOTPEnable, error)
	Disable(ctx context.Context, totpId uuid.UUID, userId uuid.UUID) (*models.TOTPDisable, error)
}
