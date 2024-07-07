package registration

import (
	"github.com/GCFactory/dbo-system/service/account/internal/models"
	"github.com/google/uuid"
	"golang.org/x/net/context"
)

type UseCase interface {
	ReservAcc(ctx context.Context, acc_data *models.FullAccountData) error
	CreateAcc(ctx context.Context, acc_uuid uuid.UUID) error
	OpenAcc(ctx context.Context, acc_uuid uuid.UUID) error
	CloseAcc(ctx context.Context, acc_uuid uuid.UUID) error
	BlockAcc(ctx context.Context, acc_uuid uuid.UUID) error
	GetAccInfo(ctx context.Context, acc_uuid uuid.UUID) (*models.FullAccountData, error)
	AddingAcc(ctx context.Context, acc_uuid uuid.UUID, add_value float64) error
	WidthAcc(ctx context.Context, acc_uuid uuid.UUID, width_value float64) error
}
