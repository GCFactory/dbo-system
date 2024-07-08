package usecase

import (
	"github.com/GCFactory/dbo-system/platform/config"
	"github.com/GCFactory/dbo-system/platform/pkg/logger"
	"github.com/GCFactory/dbo-system/service/account/internal/account"
	"github.com/GCFactory/dbo-system/service/account/internal/models"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"golang.org/x/net/context"
)

type accountUC struct {
	cfg         *config.Config
	accountRepo registration.Repository
	logger      logger.Logger
}

func (UC *accountUC) ReservAcc(ctx context.Context, acc_data *models.FullAccountData) error {
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "accountUC.ReservAcc")
	defer span.Finish()

	if UC.ValidateAccMoneyValue(ctxWithTrace, acc_data.Acc_money_value) != nil {
		return ErrorWrongMoneyValue
	}

	acc, _ := UC.accountRepo.GetAccountData(ctxWithTrace, acc_data.Acc_uuid)
	if acc != nil {
		return ErrorAccIsExisting
	}

	reason, _ := UC.accountRepo.GetReserveReason(ctxWithTrace, acc_data.Acc_uuid)
	if reason != nil {
		return ErrorReasonIsExisting
	}

	acc.Acc_uuid = acc_data.Acc_uuid
	acc.Acc_status = AccStatusReserved
	acc.Acc_culc_number = acc_data.Acc_culc_number
	acc.Acc_corr_number = acc_data.Acc_corr_number
	acc.Acc_bic = acc_data.Acc_bic
	acc.Acc_cio = acc_data.Acc_cio
	acc.Acc_money_value = acc_data.Acc_money_value
	acc.Acc_money_amount = 0.0

	if UC.accountRepo.CreateAccount(ctxWithTrace, acc) != nil {
		return ErrorCreateAcc
	}

	reason.Acc_uuid = acc_data.Acc_uuid
	reason.Reason = acc_data.Reason

	if UC.accountRepo.AddReserveReason(ctxWithTrace, reason) != nil {
		if UC.accountRepo.DeleteAccount(ctxWithTrace, acc_data.Acc_uuid) != nil {
			return ErrorFatal
		}
		return ErrorAddReason
	}

	return nil
}

func (UC *accountUC) CreateAcc(ctx context.Context, acc_uuid uuid.UUID) error {
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "accountUC.CreateAcc")
	defer span.Finish()

	acc, err := UC.accountRepo.GetAccountData(ctxWithTrace, acc_uuid)
	if err != nil {
		return ErrorNoFoundAcc
	}

	if acc.Acc_status != AccStatusReserved {
		return ErrorWrongAccReservedStatus
	}

	err = UC.accountRepo.UpdateAccountStatus(ctxWithTrace, acc_uuid, AccStatusCreated)
	if err != nil {
		return ErrorUpdateAccStatus
	}

	return nil

}

func (UC *accountUC) OpenAcc(ctx context.Context, acc_uuid uuid.UUID) error {
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "accountUC.OpenAcc")
	defer span.Finish()

	acc, err := UC.accountRepo.GetAccountData(ctxWithTrace, acc_uuid)
	if err != nil {
		return ErrorNoFoundAcc
	}

	if acc.Acc_status != AccStatusCreated {
		return ErrorWrongAccOpenStatus
	}

	err = UC.accountRepo.UpdateAccountStatus(ctxWithTrace, acc_uuid, AccStatusOpen)
	if err != nil {
		return ErrorUpdateAccStatus
	}

	return nil

}

func (UC *accountUC) CloseAcc(ctx context.Context, acc_uuid uuid.UUID) error {
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "accountUC.CloseAcc")
	defer span.Finish()

	acc, err := UC.accountRepo.GetAccountData(ctxWithTrace, acc_uuid)
	if err != nil {
		return ErrorNoFoundAcc
	}

	if acc.Acc_status != AccStatusOpen {
		return ErrorWrongAccCloseStatus
	}

	err = UC.accountRepo.UpdateAccountStatus(ctxWithTrace, acc_uuid, AccStatusClose)
	if err != nil {
		return ErrorUpdateAccStatus
	}

	return nil

}

func (UC *accountUC) BlockAcc(ctx context.Context, acc_uuid uuid.UUID) error {
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "accountUC.BlockAcc")
	defer span.Finish()

	acc, err := UC.accountRepo.GetAccountData(ctxWithTrace, acc_uuid)
	if err != nil {
		return ErrorNoFoundAcc
	}

	if acc.Acc_status != AccStatusOpen {
		return ErrorWrongAccBlockStatus
	}

	err = UC.accountRepo.UpdateAccountStatus(ctxWithTrace, acc_uuid, AccStatusClose)
	if err != nil {
		return ErrorUpdateAccStatus
	}

	return nil

}

func (UC *accountUC) GetAccInfo(ctx context.Context, acc_uuid uuid.UUID) (*models.FullAccountData, error) {
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "accountUC.GetAccInfo")
	defer span.Finish()

	acc, err := UC.accountRepo.GetAccountData(ctxWithTrace, acc_uuid)
	if err != nil {
		return nil, ErrorNoFoundAcc
	}

	reason, err := UC.accountRepo.GetReserveReason(ctxWithTrace, acc_uuid)
	if err != nil {
		return nil, ErrorGetAccReservedReason
	}

	result := models.FullAccountData{
		Acc_uuid:         acc_uuid,
		Reason:           reason.Reason,
		Acc_money_amount: acc.Acc_money_amount,
		Acc_money_value:  acc.Acc_money_value,
		Acc_cio:          acc.Acc_cio,
		Acc_bic:          acc.Acc_bic,
		Acc_corr_number:  acc.Acc_corr_number,
		Acc_culc_number:  acc.Acc_culc_number,
		Acc_status:       acc.Acc_status,
	}

	return &result, nil

}

func (UC *accountUC) AddingAcc(ctx context.Context, acc_uuid uuid.UUID, add_value float64) error {
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "accountUC.AddingAcc")
	defer span.Finish()

	acc, err := UC.accountRepo.GetAccountData(ctxWithTrace, acc_uuid)
	if err != nil {
		return ErrorNoFoundAcc
	}

	new_value := acc.Acc_money_amount + add_value

	if new_value < acc.Acc_money_amount {
		return ErrorOverflowAmount
	}

	err = UC.accountRepo.UpdateAccountAmount(ctxWithTrace, acc_uuid, new_value)
	if err != nil {
		return ErrorUpdateAmountValue
	}

	return nil

}

func (UC *accountUC) WidthAcc(ctx context.Context, acc_uuid uuid.UUID, width_value float64) error {
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "accountUC.AddingAcc")
	defer span.Finish()

	acc, err := UC.accountRepo.GetAccountData(ctxWithTrace, acc_uuid)
	if err != nil {
		return ErrorNoFoundAcc
	}

	new_value := acc.Acc_money_amount - width_value
	if new_value < 0 {
		return ErrorNotEnoughMoneyAmount
	}

	err = UC.accountRepo.UpdateAccountAmount(ctxWithTrace, acc_uuid, new_value)
	if err != nil {
		return ErrorUpdateAmountValue
	}

	return nil

}

func (UC *accountUC) ValidateAccStatus(ctx context.Context, status uint8) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "accountUC.ValidateAccStatus")
	defer span.Finish()

	for i := 0; i < len(PossibleAccStatus); i++ {
		if status == PossibleAccStatus[i] {
			return nil
		}
	}

	return ErrorWrongAccStatus
}

func (UC *accountUC) ValidateAccMoneyValue(ctx context.Context, money_value uint8) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "accountUC.ValidateAccMoneyValue")
	defer span.Finish()

	for i := 0; i < len(PossibleAccMoneyValue); i++ {
		if money_value == PossibleAccMoneyValue[i] {
			return nil
		}
	}

	return ErrorWrongMoneyValue
}

func NewAccountUseCase(cfg *config.Config, account_repo registration.Repository, log logger.Logger) registration.UseCase {
	return &accountUC{cfg: cfg, accountRepo: account_repo, logger: log}
}
