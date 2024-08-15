//go:generate mockgen -source pg_repository.go -destination mock/pg_repository_mock.go -package mock
package registration

import (
	"context"
	"github.com/GCFactory/dbo-system/service/account/internal/models"
	"github.com/google/uuid"
)

type Repository interface {
	CreateAccount(ctx context.Context, account *models.Account, reason *models.ReserverReason) error
	DeleteAccount(ctx context.Context, acc_uuid uuid.UUID) error
	GetAccountData(ctx context.Context, acc_uuid uuid.UUID) (*models.Account, error)
	UpdateAccountStatus(ctx context.Context, acc_uuid uuid.UUID, new_status uint8) error
	GetAccountStatus(ctx context.Context, acc_uuid uuid.UUID) (uint8, error)
	UpdateAccountAmount(ctx context.Context, acc_uuid uuid.UUID, acc_new_amount float64) error
	GetAccountAmount(ctx context.Context, acc_uuid uuid.UUID) (float64, error)
	DeleteReserveReason(ctx context.Context, acc_uuid uuid.UUID) error
	GetReserveReason(ctx context.Context, acc_uuid uuid.UUID) (*models.ReserverReason, error)
}
