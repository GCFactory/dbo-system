//go:generate mockgen -source usecase.go -destination mock/usecase_mock.go -package mock
package account

import (
	"github.com/GCFactory/dbo-system/service/account/internal/models"
	"github.com/google/uuid"
	"golang.org/x/net/context"
)

type UseCase interface {
	ValidateAccStatus(ctx context.Context, status uint8) error
	ValidateAccMoneyValue(ctx context.Context, money_value uint8) error
	ValidateCulcNumber(ctx context.Context, culc_number string) error
	ValidateOwner(ctx context.Context, owner string) error
	ValidateActivity(ctx context.Context, owner string, activity string) error
	ValidateCurrency(ctx context.Context, currency string) (uint8, error)
	ValidateBIC(ctx context.Context, BIC string) error
	ValidatePaymentSystem(ctx context.Context, payment_system string) error
	ValidateAccCountry(ctx context.Context, acc_country string) error
	ValidateAccCountryRegion(ctx context.Context, acc_country string, acc_country_region string) error
	ValidateAccMainOffice(ctx context.Context, acc_country string, acc_country_region string, acc_main_office string) error
	ValidateAccBankNumber(ctx context.Context, acc_bank_number string) error
	ValidateCorrNumber(ctx context.Context, acc_corr_number string, acc_bic string) error
	ValidateAccCorrNumberOwner(ctx context.Context, acc_corr_number_owner string) error
	ValidateKPP(ctx context.Context, acc_kpp string, acc_bic string) error
	ReservAcc(ctx context.Context, acc_data *models.FullAccountData) error
	CreateAcc(ctx context.Context, acc_uuid uuid.UUID) error
	OpenAcc(ctx context.Context, acc_uuid uuid.UUID) error
	CloseAcc(ctx context.Context, acc_uuid uuid.UUID) error
	BlockAcc(ctx context.Context, acc_uuid uuid.UUID) error
	GetAccInfo(ctx context.Context, acc_uuid uuid.UUID) (*models.FullAccountData, error)
	AddingAcc(ctx context.Context, acc_uuid uuid.UUID, add_value float64) error
	WidthAcc(ctx context.Context, acc_uuid uuid.UUID, width_value float64) error
}
