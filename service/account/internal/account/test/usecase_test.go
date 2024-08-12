package test

import (
	"context"
	"github.com/GCFactory/dbo-system/platform/config"
	"github.com/GCFactory/dbo-system/platform/pkg/logger"
	"github.com/GCFactory/dbo-system/service/account/internal/account/mock"
	"github.com/GCFactory/dbo-system/service/account/internal/account/usecase"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	testCfgUC = &config.Config{
		Env: "Development",
		Logger: config.Logger{
			Development: true,
			Level:       "Debug",
		},
	}
)

func TestAccountUC_ValidateAccStatus(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLogger := logger.NewServerLogger(testCfgUC)
	apiLogger.InitLogger()

	mockRepo := mock.NewMockRepository(ctrl)
	accUC := usecase.NewAccountUseCase(testCfgUC, mockRepo, apiLogger)

	t.Run("Success", func(t *testing.T) {
		err := accUC.ValidateAccStatus(context.Background(), usecase.AccStatusReserved)
		require.Nil(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		err := accUC.ValidateAccStatus(context.Background(), 12)
		require.Equal(t, err, usecase.ErrorWrongAccStatus)
	})
}

func TestAccountUC_ValidateAccMoneyValue(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLogger := logger.NewServerLogger(testCfgUC)
	apiLogger.InitLogger()

	mockRepo := mock.NewMockRepository(ctrl)
	accUC := usecase.NewAccountUseCase(testCfgUC, mockRepo, apiLogger)

	t.Run("Success", func(t *testing.T) {
		err := accUC.ValidateAccMoneyValue(context.Background(), usecase.AccMoneyValueRub)
		require.Nil(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		err := accUC.ValidateAccMoneyValue(context.Background(), 12)
		require.Equal(t, err, usecase.ErrorWrongMoneyValue)
	})
}

func TestAccountUC_ValidateOwner(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLogger := logger.NewServerLogger(testCfgUC)
	apiLogger.InitLogger()

	mockRepo := mock.NewMockRepository(ctrl)
	accUC := usecase.NewAccountUseCase(testCfgUC, mockRepo, apiLogger)

	t.Run("Success", func(t *testing.T) {
		err := accUC.ValidateOwner(context.Background(), "417")
		require.Nil(t, err)
	})

	t.Run("Wrong owner len", func(t *testing.T) {
		err := accUC.ValidateOwner(context.Background(), "1")
		require.Equal(t, err, usecase.ErrorWrongOwnerLen)
	})

	t.Run("Wrong owner", func(t *testing.T) {
		err := accUC.ValidateOwner(context.Background(), "999")
		require.Equal(t, err, usecase.ErrorWrongOwner)
	})
}

func TestAccountUC_ValidateActivity(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLogger := logger.NewServerLogger(testCfgUC)
	apiLogger.InitLogger()

	mockRepo := mock.NewMockRepository(ctrl)
	accUC := usecase.NewAccountUseCase(testCfgUC, mockRepo, apiLogger)

	t.Run("Success", func(t *testing.T) {
		err := accUC.ValidateActivity(context.Background(), "407", "02")
		require.Nil(t, err)
	})

	t.Run("Wrong owner validation", func(t *testing.T) {
		err := accUC.ValidateActivity(context.Background(), "", "02")
		require.Equal(t, err, usecase.ErrorWrongOwnerLen)
	})

	t.Run("Wrong activity len", func(t *testing.T) {
		err := accUC.ValidateActivity(context.Background(), "407", "")
		require.Equal(t, err, usecase.ErrorWrongActivityLen)
	})

	t.Run("No owner data", func(t *testing.T) {
		err := accUC.ValidateActivity(context.Background(), "301", "02")
		require.Equal(t, err, usecase.ErrorWrongActivity)
	})

	t.Run("Wrong activity", func(t *testing.T) {
		err := accUC.ValidateActivity(context.Background(), "407", "99")
		require.Equal(t, err, usecase.ErrorWrongActivity)
	})
}

func TestAccountUC_ValidateCurrency(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLogger := logger.NewServerLogger(testCfgUC)
	apiLogger.InitLogger()

	mockRepo := mock.NewMockRepository(ctrl)
	accUC := usecase.NewAccountUseCase(testCfgUC, mockRepo, apiLogger)

	t.Run("Success", func(t *testing.T) {
		currency, ok := usecase.PossibleAccMoneyValueStr["810"]
		require.Equal(t, ok, true)
		require.Equal(t, currency, usecase.AccMoneyValueRub)

		result, err := accUC.ValidateCurrency(context.Background(), "810")
		require.Nil(t, err)
		require.Equal(t, result, currency)
	})

	t.Run("Wrong currency len", func(t *testing.T) {
		result, err := accUC.ValidateCurrency(context.Background(), "")
		require.Equal(t, err, usecase.ErrorWrongCurrencyLen)
		require.Equal(t, result, usecase.AccMoneyValueError)
	})

	t.Run("Wrong currency", func(t *testing.T) {
		result, err := accUC.ValidateCurrency(context.Background(), "123")
		require.Equal(t, err, usecase.ErrorWrongCurrencyValue)
		require.Equal(t, result, usecase.AccMoneyValueUnknown)
	})
}

func TestAccountUC_ValidateCulcNumber(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLogger := logger.NewServerLogger(testCfgUC)
	apiLogger.InitLogger()

	mockRepo := mock.NewMockRepository(ctrl)
	accUC := usecase.NewAccountUseCase(testCfgUC, mockRepo, apiLogger)

	t.Run("Wrong culc number len", func(t *testing.T) {
		err := accUC.ValidateCulcNumber(context.Background(), "")
		require.Equal(t, err, usecase.ErrorWrongCulcNumberLen)
	})

	t.Run("Wrong owner", func(t *testing.T) {
		err := accUC.ValidateCulcNumber(context.Background(), "01234567890123456789")
		require.Equal(t, err, usecase.ErrorWrongOwner)
	})

	t.Run("Wrong activity", func(t *testing.T) {
		err := accUC.ValidateCulcNumber(context.Background(), "40799999990123456789")
		require.Equal(t, err, usecase.ErrorWrongActivity)
	})

	t.Run("Wrong currency", func(t *testing.T) {
		err := accUC.ValidateCulcNumber(context.Background(), "40704123990123456789")
		require.Equal(t, err, usecase.ErrorWrongCurrencyValue)
	})

	t.Run("Success", func(t *testing.T) {
		err := accUC.ValidateCulcNumber(context.Background(), "40705810990123456789")
		require.Nil(t, err)
	})
}

func TestAccountUC_ValidatePaymentSystem(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLogger := logger.NewServerLogger(testCfgUC)
	apiLogger.InitLogger()

	mockRepo := mock.NewMockRepository(ctrl)
	accUC := usecase.NewAccountUseCase(testCfgUC, mockRepo, apiLogger)

	t.Run("Wrong payment system len", func(t *testing.T) {
		err := accUC.ValidatePaymentSystem(context.Background(), "")
		require.Equal(t, err, usecase.ErrorWrongPaymentSystemLen)
	})

	t.Run("Wrong input data", func(t *testing.T) {
		err := accUC.ValidatePaymentSystem(context.Background(), "a")
		require.NotNil(t, err)
	})

	t.Run("Wrong payment system", func(t *testing.T) {
		err := accUC.ValidatePaymentSystem(context.Background(), "9")
		require.Equal(t, err, usecase.ErrorWrongPaymentSystem)
	})

	t.Run("Success", func(t *testing.T) {
		err := accUC.ValidatePaymentSystem(context.Background(), "2")
		require.Nil(t, err)
	})
}

func TestAccountUC_ValidateAccCountry(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLogger := logger.NewServerLogger(testCfgUC)
	apiLogger.InitLogger()

	mockRepo := mock.NewMockRepository(ctrl)
	accUC := usecase.NewAccountUseCase(testCfgUC, mockRepo, apiLogger)

	t.Run("Wrong country len", func(t *testing.T) {
		err := accUC.ValidateAccCountry(context.Background(), "")
		require.Equal(t, err, usecase.ErrorWrongAccCountryLen)
	})

	t.Run("Wrong input data", func(t *testing.T) {
		err := accUC.ValidateAccCountry(context.Background(), "a")
		require.NotNil(t, err)
	})

	t.Run("Wrong country", func(t *testing.T) {
		err := accUC.ValidateAccCountry(context.Background(), "9")
		require.Equal(t, err, usecase.ErrorWrongAccCountry)
	})

	t.Run("Success", func(t *testing.T) {
		err := accUC.ValidateAccCountry(context.Background(), "4")
		require.Nil(t, err)
	})
}

func TestAccountUC_ValidateAccCountryRegion(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLogger := logger.NewServerLogger(testCfgUC)
	apiLogger.InitLogger()

	mockRepo := mock.NewMockRepository(ctrl)
	accUC := usecase.NewAccountUseCase(testCfgUC, mockRepo, apiLogger)

	t.Run("Wrong country region len", func(t *testing.T) {
		err := accUC.ValidateAccCountryRegion(context.Background(), "", "")
		require.Equal(t, err, usecase.ErrorWrongAccCountryRegionLen)
	})

	t.Run("Wrong input data", func(t *testing.T) {
		err := accUC.ValidateAccCountryRegion(context.Background(), "4", "ab")
		require.NotNil(t, err)
	})

	t.Run("Wrong country validation", func(t *testing.T) {
		err := accUC.ValidateAccCountryRegion(context.Background(), "9", "52")
		require.Equal(t, err, usecase.ErrorWrongAccCountry)
	})

	t.Run("No country data", func(t *testing.T) {
		err := accUC.ValidateAccCountryRegion(context.Background(), "0", "12")
		require.Equal(t, err, usecase.ErrorWrongAccCountry)
	})

	t.Run("Wrong country region", func(t *testing.T) {
		err := accUC.ValidateAccCountryRegion(context.Background(), "4", "52")
		require.Equal(t, err, usecase.ErrorWrongAccCountryRegion)
	})

	t.Run("Success", func(t *testing.T) {
		err := accUC.ValidateAccCountryRegion(context.Background(), "4", "50")
		require.Nil(t, err)
	})
}

func TestAccountUC_ValidateAccMainOffice(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLogger := logger.NewServerLogger(testCfgUC)
	apiLogger.InitLogger()

	mockRepo := mock.NewMockRepository(ctrl)
	accUC := usecase.NewAccountUseCase(testCfgUC, mockRepo, apiLogger)

	t.Run("Wrong main office len", func(t *testing.T) {
		err := accUC.ValidateAccMainOffice(context.Background(), "", "", "")
		require.Equal(t, err, usecase.ErrorWrongAccMainOfficeLen)
	})

	t.Run("Wrong country", func(t *testing.T) {
		err := accUC.ValidateAccMainOffice(context.Background(), "", "41", "12")
		require.Equal(t, err, usecase.ErrorWrongAccCountryLen)
	})

	t.Run("Wrong country region len", func(t *testing.T) {
		err := accUC.ValidateAccMainOffice(context.Background(), "4", "", "12")
		require.Equal(t, err, usecase.ErrorWrongAccCountryRegionLen)
	})

	t.Run("Wrong main office data", func(t *testing.T) {
		err := accUC.ValidateAccMainOffice(context.Background(), "4", "50", "ab")
		require.NotNil(t, err)
	})

	t.Run("No region data", func(t *testing.T) {
		err := accUC.ValidateAccMainOffice(context.Background(),
			"4",
			"46",
			"12",
		)
		require.Equal(t, err, usecase.ErrorWrongAccCountryRegion)
	})

	t.Run("Wrong main office", func(t *testing.T) {
		err := accUC.ValidateAccMainOffice(context.Background(), "4", "50", "99")
		require.Equal(t, err, usecase.ErrorWrongAccMainOffice)
	})

	t.Run("Success", func(t *testing.T) {
		err := accUC.ValidateAccMainOffice(context.Background(), "4", "50", "25")
		require.Nil(t, err)
	})

}

func TestAccountUC_ValidateAccBankNumber(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLogger := logger.NewServerLogger(testCfgUC)
	apiLogger.InitLogger()

	mockRepo := mock.NewMockRepository(ctrl)
	accUC := usecase.NewAccountUseCase(testCfgUC, mockRepo, apiLogger)

	t.Run("Wrong bank number len", func(t *testing.T) {
		err := accUC.ValidateAccBankNumber(context.Background(), "")
		require.Equal(t, err, usecase.ErrorWrongAccBankNumberLen)
	})

	t.Run("Wrong bank number data", func(t *testing.T) {
		err := accUC.ValidateAccBankNumber(context.Background(), "abc")
		require.NotNil(t, err)
	})

	t.Run("Wrong bank number", func(t *testing.T) {
		err := accUC.ValidateAccBankNumber(context.Background(), "999")
		require.Equal(t, err, usecase.ErrorWrongAccBankNumber)
	})

	t.Run("Success", func(t *testing.T) {
		err := accUC.ValidateAccBankNumber(context.Background(), "025")
		require.Nil(t, err)
	})

}

func TestAccountUC_ValidateBIC(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLogger := logger.NewServerLogger(testCfgUC)
	apiLogger.InitLogger()

	mockRepo := mock.NewMockRepository(ctrl)
	accUC := usecase.NewAccountUseCase(testCfgUC, mockRepo, apiLogger)

	t.Run("Wrong BIC len", func(t *testing.T) {
		err := accUC.ValidateBIC(context.Background(), "")
		require.Equal(t, err, usecase.ErrorWrongBICLen)
	})

	t.Run("Wrong payment system", func(t *testing.T) {
		err := accUC.ValidateBIC(context.Background(), "987654321")
		require.Equal(t, err, usecase.ErrorWrongPaymentSystem)
	})

	t.Run("Wrong acc country", func(t *testing.T) {
		err := accUC.ValidateBIC(context.Background(), "299999999")
		require.Equal(t, err, usecase.ErrorWrongAccCountry)
	})

	t.Run("Wrong acc country region", func(t *testing.T) {
		err := accUC.ValidateBIC(context.Background(), "249999999")
		require.Equal(t, err, usecase.ErrorWrongAccCountryRegion)
	})

	t.Run("Wrong main office", func(t *testing.T) {
		err := accUC.ValidateBIC(context.Background(), "245099999")
		require.Equal(t, err, usecase.ErrorWrongAccMainOffice)
	})

	t.Run("Wrong bank number", func(t *testing.T) {
		err := accUC.ValidateBIC(context.Background(), "245025999")
		require.Equal(t, err, usecase.ErrorWrongAccBankNumber)
	})

	t.Run("Success", func(t *testing.T) {
		err := accUC.ValidateBIC(context.Background(), "245025025")
		require.Nil(t, err)
	})
}

//func TestAccountUC_ReservAcc(t *testing.T) {
//	t.Parallel()
//
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	apiLogger := logger.NewServerLogger(testCfgUC)
//	apiLogger.InitLogger()
//
//	mockRepo := mock.NewMockRepository(ctrl)
//	accUC := usecase.NewAccountUseCase(testCfgUC, mockRepo, apiLogger)
//
//	accUC = accUC
//
//	//acc_data := models.FullAccountData{
//	//	Acc_uuid: uuid.New(),
//	//	Acc_culc_number: Acc_culc_number,
//	//	acc.Acc_corr_number = acc_data.Acc_corr_number
//	//	acc.Acc_bic = acc_data.Acc_bic
//	//	acc.Acc_cio = acc_data.Acc_cio
//	//	acc.Acc_money_value = acc_data.Acc_money_value
//	//	acc.Acc_money_amount = 0.0
//	//}
//
//	t.Run("Success", func(t *testing.T) {
//		//	TODO: реализовать
//	})
//}
