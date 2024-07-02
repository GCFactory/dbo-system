//go:generate mockgen -source pg_repository.go -destination mock/pg_repository_mock.go -package mock
package registration

import (
	"context"
	"github.com/GCFactory/dbo-system/service/account/internal/models"
	"github.com/google/uuid"
)

type Repository interface {
	CreateAccount(ctx context.Context, account models.Account) error
	DeleteAccount(ctx context.Context, acc_uuid uuid.UUID) error
	GetAccountData(ctx context.Context, acc_uuid uuid.UUID) (*models.Account, error)
	UpdateAccountStatus(ctx context.Context, acc_uuid uuid.UUID) error
	GetAccountStatus(ctx context.Context, acc_uuid uuid.UUID) (uint8, error)
	UpdateAccountAmount(ctx context.Context, acc_uuid uuid.UUID, acc_new_amount float32) error
	GetAccountAmount(ctx context.Context, acc_uuid uuid.UUID) (float64, error)
	AddReserveReason(ctx context.Context, acc_uuid uuid.UUID, reason string) error
	GetReserveReason(ctx context.Context, acc_uuid uuid.UUID) (string, error)
}
