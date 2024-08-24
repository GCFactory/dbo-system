package grpc_handlers

import (
	"fmt"
	"github.com/GCFactory/dbo-system/platform/pkg/logger"
	"github.com/GCFactory/dbo-system/service/account/config"
	acc_proto_api "github.com/GCFactory/dbo-system/service/account/gen_proto/proto/api"
	acc_proto_platform "github.com/GCFactory/dbo-system/service/account/gen_proto/proto/platform"
	"github.com/GCFactory/dbo-system/service/account/internal/account"
	"github.com/GCFactory/dbo-system/service/account/internal/account/usecase"
	"github.com/GCFactory/dbo-system/service/account/internal/models"
	"github.com/GCFactory/dbo-system/service/account/pkg/kafka"
	"github.com/IBM/sarama"
	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"golang.org/x/net/context"
)

type AccountGRPCHandlers struct {
	cfg       *config.Config
	kProducer *kafka.ProducerProvider
	accUC     account.UseCase
	accLog    logger.Logger
}

func (accGRPCH AccountGRPCHandlers) ReserveAccount(ctx context.Context, saga_uuid string, acc_data *acc_proto_platform.AccountDetails, kProducer *kafka.ProducerProvider) error {

	fmt.Printf("Message claimed: value = %s\n", acc_data)
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "accGRPCH.ReserveAccount")
	defer span.Finish()

	acc_model := &models.FullAccountData{
		Acc_uuid:        uuid.New(),
		Acc_corr_number: acc_data.GetCorrNumber(),
		Acc_culc_number: acc_data.GetCulcNumber(),
		Acc_cio:         acc_data.GetCio(),
		Acc_bic:         acc_data.GetBic(),
		Reason:          acc_data.GetReserveReason(),
	}

	flag_error := false
	var err error
	answer_topic := TopicResult

	if err = accGRPCH.accUC.ReservAcc(ctxWithTrace, acc_model); err != nil {
		accGRPCH.accLog.Error(err)
		flag_error = true
		answer_topic = TopicError
	}
	answer := &acc_proto_api.EventStatus{
		SagaUuid: saga_uuid,
	}
	if flag_error {
		switch err {
		case usecase.ErrorCreateAcc:
			{
				answer.Status = 500
			}
		default:
			answer.Status = 412
		}
		answer.Info = err.Error()
	} else {
		answer.Status = 201
		answer.Info = "Account reserved"
	}

	answer_data, err := proto.Marshal(answer)
	if err != nil {
		accGRPCH.accLog.Error(err)
		return err
	}

	err = kProducer.ProduceRecord(answer_topic, sarama.ByteEncoder(answer_data))
	if err != nil {
		accGRPCH.accLog.Error(err)
		return err
	}
	accGRPCH.accLog.Info("Success to send answer!")

	return nil
}

func NewAccountGRPCHandlers(cfg *config.Config, kProducer *kafka.ProducerProvider, accUC account.UseCase, accLog logger.Logger) account.GRPCHandlers {
	return &AccountGRPCHandlers{cfg: cfg, kProducer: kProducer, accUC: accUC, accLog: accLog}
}
