package usecase

import (
	"github.com/GCFactory/dbo-system/platform/config"
	"github.com/GCFactory/dbo-system/platform/pkg/logger"
	"github.com/GCFactory/dbo-system/service/account/internal/account"
	"github.com/GCFactory/dbo-system/service/account/internal/models"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"golang.org/x/net/context"
	"strconv"
)

type accountUC struct {
	cfg         *config.Config
	accountRepo registration.Repository
	logger      logger.Logger
}

func (UC *accountUC) ReservAcc(ctx context.Context, acc_data *models.FullAccountData) error {
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "accountUC.ReservAcc")
	defer span.Finish()

	acc, _ := UC.accountRepo.GetAccountData(ctxWithTrace, acc_data.Acc_uuid)
	if acc != nil {
		return ErrorAccIsExisting
	}

	reason, _ := UC.accountRepo.GetReserveReason(ctxWithTrace, acc_data.Acc_uuid)
	if reason != nil {
		return ErrorReasonIsExisting
	}

	acc = &models.Account{
		Acc_uuid:         acc_data.Acc_uuid,
		Acc_status:       AccStatusReserved,
		Acc_culc_number:  acc_data.Acc_culc_number,
		Acc_corr_number:  acc_data.Acc_corr_number,
		Acc_bic:          acc_data.Acc_bic,
		Acc_cio:          acc_data.Acc_cio,
		Acc_money_value:  acc_data.Acc_money_value,
		Acc_money_amount: 0.0,
	}

	if err := UC.ValidateCulcNumber(ctxWithTrace, acc_data.Acc_culc_number); err != nil {
		return err
	} else if err = UC.ValidateBIC(ctxWithTrace, acc_data.Acc_bic); err != nil {
		return err
	} else if err = UC.ValidateCorrNumber(ctxWithTrace, acc_data.Acc_corr_number, acc_data.Acc_bic); err != nil {
		return err
	} else if err = UC.ValidateKPP(ctxWithTrace, acc_data.Acc_cio, acc_data.Acc_bic); err != nil {
		return err
	}

	reason = &models.ReserverReason{
		Acc_uuid: acc_data.Acc_uuid,
		Reason:   acc_data.Reason,
	}

	if UC.accountRepo.CreateAccount(ctxWithTrace, acc, reason) != nil {
		return ErrorCreateAcc
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

	err = UC.accountRepo.UpdateAccountStatus(ctxWithTrace, acc_uuid, AccStatusBlocked)
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

	result := models.FullAccountData{
		Acc_uuid:         acc_uuid,
		Acc_money_amount: acc.Acc_money_amount,
		Acc_money_value:  acc.Acc_money_value,
		Acc_cio:          acc.Acc_cio,
		Acc_bic:          acc.Acc_bic,
		Acc_corr_number:  acc.Acc_corr_number,
		Acc_culc_number:  acc.Acc_culc_number,
		Acc_status:       acc.Acc_status,
	}

	reason, err := UC.accountRepo.GetReserveReason(ctxWithTrace, acc_uuid)
	if err != nil {
		result.Reason = ""
	} else {
		result.Reason = reason.Reason
	}

	return &result, nil

}

func (UC *accountUC) AddingAcc(ctx context.Context, acc_uuid uuid.UUID, add_value float64) error {
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "accountUC.AddingAcc")
	defer span.Finish()

	acc, err := UC.accountRepo.GetAccountData(ctxWithTrace, acc_uuid)
	if err != nil {
		return ErrorNoFoundAcc
	} else if acc.Acc_status != AccStatusOpen {
		return ErrorWrongAccOpenStatus
	}

	new_value := acc.Acc_money_amount + add_value

	if new_value < acc.Acc_money_amount || new_value == acc.Acc_money_amount {
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
	} else if acc.Acc_status != AccStatusOpen {
		return ErrorWrongAccOpenStatus
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

// Валидирует статус счёта
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

// Валидирует денежную величину
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

// Валидирует расчётный номер счёта
func (UC *accountUC) ValidateCulcNumber(ctx context.Context, culc_number string) error {
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "accountUC.ValidateCulcNumber")
	defer span.Finish()

	if len(culc_number) != 20 {
		return ErrorWrongCulcNumberLen
	}

	owner := culc_number[:3]
	activity_field := culc_number[3:5]
	currency := culc_number[5:8]

	if err := UC.ValidateOwner(ctxWithTrace, owner); err != nil {
		return err
	}

	if err := UC.ValidateActivity(ctxWithTrace, owner, activity_field); err != nil {
		return err
	}

	if _, err := UC.ValidateCurrency(ctxWithTrace, currency); err != nil {
		return err
	}

	return nil
}

// Проверяет владельца счёта
func (UC *accountUC) ValidateOwner(ctx context.Context, owner string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "accountUC.ValidateOwner")
	defer span.Finish()

	if len(owner) != 3 {
		return ErrorWrongOwnerLen
	}

	values, ok := PossibleOwner[string(owner[0])]
	if ok {
		for i := 0; i < len(values); i++ {
			if owner == values[i] {
				return nil
			}
		}
	}

	return ErrorWrongOwner
}

// Проверяет тип деятелности владельца счёта
func (UC *accountUC) ValidateActivity(ctx context.Context, owner string, activity string) error {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "accountUC.ValidateOwner")
	defer span.Finish()

	if len(activity) != 2 {
		return ErrorWrongActivityLen
	}

	values, ok := PossibleAccActivity[owner]
	if err := UC.ValidateOwner(ctxWithTrace, owner); err != nil {
		return err
	} else if !ok {
		return ErrorWrongActivity
	} else {
		for i := 0; i < len(values); i++ {
			if activity == values[i] {
				return nil
			}
		}
	}

	return ErrorWrongActivity
}

// Проверят тип валюты
func (UC *accountUC) ValidateCurrency(ctx context.Context, currency string) (uint8, error) {

	span, _ := opentracing.StartSpanFromContext(ctx, "accountUC.ValidateCurrency")
	defer span.Finish()

	if len(currency) != 3 {
		return AccMoneyValueError, ErrorWrongCurrencyLen
	}

	local_currency, ok := PossibleAccMoneyValueStr[currency]

	if !ok {
		return AccMoneyValueUnknown, ErrorWrongCurrencyValue
	}

	return local_currency, nil

}

// Валидирует БИК счёта
func (UC *accountUC) ValidateBIC(ctx context.Context, BIC string) error {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "accountUC.ValidateBIC")
	defer span.Finish()

	if len(BIC) != 9 {
		return ErrorWrongBICLen
	}

	if err := UC.ValidatePaymentSystem(ctxWithTrace, BIC[:1]); err != nil {
		return err
	} else if err = UC.ValidateAccCountry(ctxWithTrace, BIC[1:2]); err != nil {
		return err
	} else if err = UC.ValidateAccCountryRegion(ctxWithTrace, BIC[1:2], BIC[2:4]); err != nil {
		return err
	} else if err = UC.ValidateAccMainOffice(ctxWithTrace, BIC[1:2], BIC[2:4], BIC[4:6]); err != nil {
		return err
	} else if err = UC.ValidateAccBankNumber(ctxWithTrace, BIC[6:9]); err != nil {
		return err
	}

	return nil
}

// Валидирует участника платёжной системы
func (UC *accountUC) ValidatePaymentSystem(ctx context.Context, payment_system string) error {

	span, _ := opentracing.StartSpanFromContext(ctx, "accountUC.ValidatePaymentSystem")
	defer span.Finish()

	if len(payment_system) != 1 {
		return ErrorWrongPaymentSystemLen
	}

	local_payment_system, err := strconv.Atoi(payment_system)
	if err != nil {
		return err
	}

	for i := 0; i < len(PossiblePaymentSystems); i++ {
		if uint8(local_payment_system) == PossiblePaymentSystems[i] {
			return nil
		}
	}

	return ErrorWrongPaymentSystem
}

// Валидирует страну счёта
func (UC *accountUC) ValidateAccCountry(ctx context.Context, acc_country string) error {

	span, _ := opentracing.StartSpanFromContext(ctx, "accountUC.ValidateAccCountry")
	defer span.Finish()

	if len(acc_country) != 1 {
		return ErrorWrongAccCountryLen
	}

	local_acc_country, err := strconv.Atoi(acc_country)
	if err != nil {
		return err
	}

	for i := 0; i < len(PossibleAccCountry); i++ {
		if uint8(local_acc_country) == PossibleAccCountry[i] {
			return nil
		}
	}

	return ErrorWrongAccCountry

}

// Валидирует регион страны
func (UC *accountUC) ValidateAccCountryRegion(ctx context.Context, acc_country string, acc_country_region string) error {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "accountUC.ValidateAccCountryRegion")
	defer span.Finish()

	if len(acc_country_region) != 2 {
		return ErrorWrongAccCountryRegionLen
	}

	if err := UC.ValidateAccCountry(ctxWithTrace, acc_country); err != nil {
		return err
	}

	local_country_region, err := strconv.Atoi(acc_country_region)
	if err != nil {
		return err
	}

	local_country, _ := strconv.Atoi(acc_country)

	list_of_regions, ok := PossibleAccCountryRegion[uint8(local_country)]
	if !ok {
		return ErrorWrongAccCountry
	} else {
		for i := 0; i < len(list_of_regions); i++ {
			if uint8(local_country_region) == list_of_regions[i] {
				return nil
			}
		}
	}

	return ErrorWrongAccCountryRegion
}

// Валидирует главный офис региона
func (UC *accountUC) ValidateAccMainOffice(
	ctx context.Context,
	acc_country string,
	acc_country_region string,
	acc_main_office string) error {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "accountUC.ValidateAccMainOffice")
	defer span.Finish()

	if len(acc_main_office) != 2 {
		return ErrorWrongAccMainOfficeLen
	}

	if err := UC.ValidateAccCountryRegion(ctxWithTrace, acc_country, acc_country_region); err != nil {
		return err
	}

	local_main_office, err := strconv.Atoi(acc_main_office)
	if err != nil {
		return err
	}
	local_country_region, _ := strconv.Atoi(acc_country_region)

	data_main_office, ok := PossibleAccMainOffice[uint8(local_country_region)]
	if !ok {
		return ErrorWrongAccCountryRegion
	}
	if uint8(local_main_office) == data_main_office {
		return nil
	} else {
		return ErrorWrongAccMainOffice
	}

}

// Валидирует код банка
func (UC *accountUC) ValidateAccBankNumber(ctx context.Context, acc_bank_number string) error {

	span, _ := opentracing.StartSpanFromContext(ctx, "accountUC.ValidateAccCountryRegion")
	defer span.Finish()

	if len(acc_bank_number) != 3 {
		return ErrorWrongAccBankNumberLen
	}
	local_bank_number, err := strconv.Atoi(acc_bank_number)
	if err != nil {
		return err
	}

	for i := 0; i < len(PossibleAccBankNumber); i++ {
		if uint16(local_bank_number) == PossibleAccBankNumber[i] {
			return nil
		}
	}

	return ErrorWrongAccBankNumber
}

// Валидирует корр счет
func (UC *accountUC) ValidateCorrNumber(ctx context.Context, acc_corr_number string, acc_bic string) error {

	span, contextWithTrace := opentracing.StartSpanFromContext(ctx, "accountUC.ValidateCorrNumber")
	defer span.Finish()

	if err := UC.ValidateBIC(contextWithTrace, acc_bic); err != nil {
		return err
	}

	if len(acc_corr_number) != 20 {
		return ErrorWrongAccCorrNumberLen
	}

	if acc_corr_number[:3] != "301" {
		return ErrorWrongAccCorrFirstGroupNumber
	}

	if err := UC.ValidateAccCorrNumberOwner(contextWithTrace, acc_corr_number[3:5]); err != nil {
		return err
	} else if _, err = UC.ValidateCurrency(contextWithTrace, acc_corr_number[5:8]); err != nil {
		return err
	} else if acc_corr_number[9:13] != acc_bic[3:7] {
		return ErrorWrongAccCorrBankNumber
	}

	if acc_corr_number[len(acc_corr_number)-3] != acc_bic[len(acc_bic)-3] {
		return ErrorWrongAccCorrNumber
	}

	return nil
}

// Валидирует владельца корр счёта
func (UC *accountUC) ValidateAccCorrNumberOwner(ctx context.Context, acc_corr_number_owner string) error {

	span, _ := opentracing.StartSpanFromContext(ctx, "accountUC.ValidateAccCurrNumberOwner")
	defer span.Finish()

	if len(acc_corr_number_owner) != 2 {
		return ErrorWrongAccCorrOwnerLen
	}

	for i := 0; i < len(PossibleAccCurrNumberOwner); i++ {
		if acc_corr_number_owner == PossibleAccCurrNumberOwner[i] {
			return nil
		}
	}

	return ErrorWrongAccCorrOwner

}

// Валидирует КПП
func (UC *accountUC) ValidateKPP(ctx context.Context, acc_kpp string, acc_bic string) error {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "accountUC.ValidateKPP")
	defer span.Finish()

	if len(acc_kpp) != 9 {
		return ErrorWrongAccKPPLen
	} else if err := UC.ValidateBIC(ctxWithTrace, acc_bic); err != nil {
		return err
	} else if err = UC.ValidateAccCountryRegion(ctxWithTrace, acc_bic[1:2], acc_kpp[:2]); err != nil {
		return err
	}

	local_reason, err := strconv.Atoi(acc_kpp[4:6])
	if err != nil || local_reason > 99 || local_reason < 1 {
		return ErrorWrongAccKPP
	}

	return nil
}

func NewAccountUseCase(cfg *config.Config, account_repo registration.Repository, log logger.Logger) registration.UseCase {
	return &accountUC{cfg: cfg, accountRepo: account_repo, logger: log}
}
