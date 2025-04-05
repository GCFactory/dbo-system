package grpc_handlers

import (
	"github.com/GCFactory/dbo-system/platform/config"
	"github.com/GCFactory/dbo-system/platform/pkg/logger"
	acc_proto_api "github.com/GCFactory/dbo-system/service/account/gen_proto/proto/api/account"
	acc_proto_platform "github.com/GCFactory/dbo-system/service/account/gen_proto/proto/platform"
	"github.com/GCFactory/dbo-system/service/account/internal/account"
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

func (accGRPCH AccountGRPCHandlers) ReserveAccount(ctx context.Context, saga_uuid string, event_uuid string, acc_data *acc_proto_platform.AccountDetails, kProducer *kafka.ProducerProvider) error {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "accGRPCH.ReserveAccount")
	defer span.Finish()

	acc_model := &models.FullAccountData{
		Acc_uuid:        uuid.New(),
		Acc_name:        acc_data.GetAccountName(),
		Acc_corr_number: acc_data.GetCorrNumber(),
		Acc_culc_number: acc_data.GetCulcNumber(),
		Acc_cio:         acc_data.GetCio(),
		Acc_bic:         acc_data.GetBic(),
		Reason:          acc_data.GetReserveReason(),
	}

	flag_error := false
	var err error
	answer_topic := TopicResult

	answer := &acc_proto_api.EventStatus{
		SagaUuid:      saga_uuid,
		EventUuid:     event_uuid,
		OperationName: ReserveAccount,
	}

	error_answer := &acc_proto_api.EventError{
		SagaUuid:      saga_uuid,
		EventUuid:     event_uuid,
		OperationName: ReserveAccount,
	}

	if err = accGRPCH.accUC.ReservAcc(ctxWithTrace, acc_model); err != nil {
		accGRPCH.accLog.Error(err)
		flag_error = true
		answer_topic = TopicError

		error_answer.Info = err.Error()
		error_answer.Status = GetErrorCode(err)
	} else {
		answer.Result = &acc_proto_api.EventStatus_Info{acc_model.Acc_uuid.String()}
	}

	var answer_data []byte

	if flag_error {
		answer_data, err = proto.Marshal(error_answer)
	} else {
		answer_data, err = proto.Marshal(answer)
	}

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

func (accGRPCH AccountGRPCHandlers) ChangeAccountStatus(ctx context.Context, saga_uuid string, event_uuid string, operation_type string, acc_data *acc_proto_api.OperationDetails, kProducer *kafka.ProducerProvider) error {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "accGRPCH.ChangeAccountStatus")
	defer span.Finish()

	flag_error := false
	var err error
	answer_topic := TopicError

	answer := &acc_proto_api.EventStatus{
		SagaUuid:      saga_uuid,
		EventUuid:     event_uuid,
		OperationName: operation_type,
	}

	answer_error := &acc_proto_api.EventError{
		SagaUuid:      saga_uuid,
		EventUuid:     event_uuid,
		OperationName: operation_type,
	}

	account_uuid_str := acc_data.GetAccUuid()

	var account_uuid uuid.UUID
	account_uuid, err = uuid.Parse(account_uuid_str)
	if err != nil {
		answer_error.Info = err.Error()
		answer_error.Status = GetErrorCode(ErrorInvalidInputData)

		flag_error = true
	}

	if !flag_error {
		switch operation_type {
		case CreateAccount:
			err = accGRPCH.accUC.CreateAcc(ctxWithTrace, account_uuid)
		case OpenAccount:
			err = accGRPCH.accUC.OpenAcc(ctxWithTrace, account_uuid)
		case CloseAccount:
			err = accGRPCH.accUC.CloseAcc(ctxWithTrace, account_uuid)
		case BlockAccount:
			err = accGRPCH.accUC.BlockAcc(ctxWithTrace, account_uuid)
		}
		if err == nil {
			answer_topic = TopicResult
			answer.Result = &acc_proto_api.EventStatus_Info{"Success"}
		} else {
			flag_error = true
			accGRPCH.accLog.Error(err)
			answer_error.Info = err.Error()
			answer_error.Status = GetErrorCode(err)
		}
	}

	var answer_data []byte
	if flag_error {
		answer_data, err = proto.Marshal(answer_error)
	} else {
		answer_data, err = proto.Marshal(answer)
	}
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

func (accGRPCH AccountGRPCHandlers) GetAccountData(ctx context.Context, saga_uuid string, event_uuid string, acc_data *acc_proto_api.OperationDetails, kProducer *kafka.ProducerProvider) error {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "accGRPCH.GetAccountData")
	defer span.Finish()

	flag_error := false
	var err error
	answer_topic := TopicError

	answer := &acc_proto_api.EventStatus{
		SagaUuid:      saga_uuid,
		EventUuid:     event_uuid,
		OperationName: GetAccountData,
	}

	answer_error := &acc_proto_api.EventError{
		SagaUuid:      saga_uuid,
		EventUuid:     event_uuid,
		OperationName: GetAccountData,
	}

	account_uuid_str := acc_data.GetAccUuid()

	var account_uuid uuid.UUID
	account_uuid, err = uuid.Parse(account_uuid_str)
	if err != nil {
		answer_error.Info = err.Error()
		answer_error.Status = GetErrorCode(ErrorInvalidInputData)

		flag_error = true
	}

	if !flag_error {
		result, err := accGRPCH.accUC.GetAccInfo(ctxWithTrace, account_uuid)
		if err != nil {
			accGRPCH.accLog.Error(err)
			flag_error = true
		}

		if flag_error {
			answer_error.Status = GetErrorCode(err)
			answer_error.Info = err.Error()
		} else {
			answer_topic = TopicResult
			answer.Result = &acc_proto_api.EventStatus_AccData{
				AccData: &acc_proto_platform.FullAccountData{
					AccDetails: &acc_proto_platform.AccountDetails{
						CulcNumber:    result.Acc_culc_number,
						CorrNumber:    result.Acc_corr_number,
						Bic:           result.Acc_bic,
						Cio:           result.Acc_cio,
						ReserveReason: result.Reason,
						AccountName:   result.Acc_name,
					},
					AccStatus:      uint64(result.Acc_status),
					AccMoneyValue:  uint64(result.Acc_money_value),
					AccMoneyAmount: float32(result.Acc_money_amount),
				},
			}
		}
	}

	var answer_data []byte
	if flag_error {
		answer_data, err = proto.Marshal(answer_error)
	} else {
		answer_data, err = proto.Marshal(answer)
	}
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

func (accGRPCH AccountGRPCHandlers) OperationWithAccAmount(ctx context.Context, saga_uuid string, event_uuid string, operation_type string, acc_data *acc_proto_api.OperationDetails, kProducer *kafka.ProducerProvider) error {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "accGRPCH.OperationWithAccAmount")
	defer span.Finish()

	flag_error := false
	var err error
	answer_topic := TopicError

	answer := &acc_proto_api.EventStatus{
		SagaUuid:      saga_uuid,
		EventUuid:     event_uuid,
		OperationName: operation_type,
	}

	answer_error := &acc_proto_api.EventError{
		SagaUuid:      saga_uuid,
		EventUuid:     event_uuid,
		OperationName: operation_type,
	}

	account_uuid_str := acc_data.GetAccUuid()
	var account_uuid uuid.UUID
	account_uuid, err = uuid.Parse(account_uuid_str)
	if err != nil {
		answer_error.Info = err.Error()
		answer_error.Status = GetErrorCode(ErrorInvalidInputData)

		flag_error = true
	}

	if !flag_error {
		switch operation_type {
		case AddingAcc:
			err = accGRPCH.accUC.AddingAcc(ctxWithTrace, account_uuid, float64(acc_data.GetAdditionalData()))
		case WidthAcc:
			err = accGRPCH.accUC.WidthAcc(ctxWithTrace, account_uuid, float64(acc_data.GetAdditionalData()))
		}

		if err == nil {
			answer_topic = TopicResult
			answer.Result = &acc_proto_api.EventStatus_Info{"Success"}
		} else {
			flag_error = true
			accGRPCH.accLog.Error(err)
			answer_error.Status = GetErrorCode(err)
			answer_error.Info = err.Error()
		}
	}

	var answer_data []byte
	if flag_error {
		answer_data, err = proto.Marshal(answer_error)
	} else {
		answer_data, err = proto.Marshal(answer)
	}
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

func (accGRPCH AccountGRPCHandlers) RemoveAccount(ctx context.Context, saga_uuid string, event_uuid string, acc_data *acc_proto_api.OperationDetails, kProducer *kafka.ProducerProvider) error {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "accGRPCH.RemoveAccount")
	defer span.Finish()

	flag_error := false
	var err error
	answer_topic := TopicError

	answer := &acc_proto_api.EventStatus{
		SagaUuid:      saga_uuid,
		EventUuid:     event_uuid,
		OperationName: RemoveAcc,
	}

	answer_error := &acc_proto_api.EventError{
		SagaUuid:      saga_uuid,
		EventUuid:     event_uuid,
		OperationName: RemoveAcc,
	}

	account_uuid_str := acc_data.GetAccUuid()
	var account_uuid uuid.UUID
	account_uuid, err = uuid.Parse(account_uuid_str)
	if err != nil {
		answer_error.Info = err.Error()
		answer_error.Status = GetErrorCode(ErrorInvalidInputData)

		flag_error = true
	}

	if !flag_error {
		err = accGRPCH.accUC.RemoveAccount(ctxWithTrace, account_uuid)

		if err == nil {
			answer_topic = TopicResult
			answer.Result = &acc_proto_api.EventStatus_Info{"Success"}
		} else {
			flag_error = true
			accGRPCH.accLog.Error(err)
			answer_error.Status = GetErrorCode(err)
			answer_error.Info = err.Error()
		}
	}

	var answer_data []byte
	if flag_error {
		answer_data, err = proto.Marshal(answer_error)
	} else {
		answer_data, err = proto.Marshal(answer)
	}
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
