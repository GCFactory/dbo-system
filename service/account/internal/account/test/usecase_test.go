package test

import (
	"github.com/GCFactory/dbo-system/platform/config"
	"github.com/GCFactory/dbo-system/platform/pkg/logger"
	"github.com/GCFactory/dbo-system/service/account/internal/account/mock"
	"github.com/GCFactory/dbo-system/service/account/internal/account/usecase"
	"github.com/golang/mock/gomock"
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

func TestAccountUC_ReservAcc(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLogger := logger.NewServerLogger(testCfgUC)
	apiLogger.InitLogger()

	mockRepo := mock.NewMockRepository(ctrl)
	accUC := usecase.NewAccountUseCase(testCfgUC, mockRepo, apiLogger)

	accUC = accUC

	//acc_data := models.FullAccountData{
	//	Acc_uuid: uuid.New(),
	//	Acc_culc_number: Acc_culc_number,
	//	acc.Acc_corr_number = acc_data.Acc_corr_number
	//	acc.Acc_bic = acc_data.Acc_bic
	//	acc.Acc_cio = acc_data.Acc_cio
	//	acc.Acc_money_value = acc_data.Acc_money_value
	//	acc.Acc_money_amount = 0.0
	//}

	t.Run("Success", func(t *testing.T) {
		//	TODO: реализовать
	})
}
