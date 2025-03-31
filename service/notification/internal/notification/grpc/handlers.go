package grpc

import (
	"context"
	"github.com/GCFactory/dbo-system/platform/config"
	"github.com/GCFactory/dbo-system/platform/pkg/logger"
	api "github.com/GCFactory/dbo-system/service/notification/gen_proto/proto/notification_api"
	"github.com/GCFactory/dbo-system/service/notification/internal/models"
	"github.com/GCFactory/dbo-system/service/notification/internal/notification"
	"github.com/GCFactory/dbo-system/service/notification/pkg/kafka"
	"github.com/IBM/sarama"
	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
)

type NotificationGrpcHandlers struct {
	cfg       *config.Config
	kProducer *kafka.ProducerProvider
	useCase   notification.UseCase
	accLog    logger.Logger
}

func NewNotificationGRPCHandlers(cfg *config.Config, kProducer *kafka.ProducerProvider, useCase notification.UseCase, accLog logger.Logger) notification.GRPCHandlers {
	return &NotificationGrpcHandlers{cfg: cfg, kProducer: kProducer, useCase: useCase, accLog: accLog}
}

func (gh NotificationGrpcHandlers) AddUserSettings(ctx context.Context, saga_uuid string, event_uuid string, users_data *api.AdditionalInfo, kProducer *kafka.ProducerProvider) error {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "NotificationGrpcHandlers.AddUserSettings")
	defer span.Finish()

	var answer_data []byte
	flag_error := false
	var err error
	answer_topic := TopicResult

	answer := &api.EventSuccess{
		SagaUuid:      saga_uuid,
		EventUuid:     event_uuid,
		OperationName: OperationAddUserSettings,
	}

	error_answer := &api.EventError{
		SagaUuid:      saga_uuid,
		EventUuid:     event_uuid,
		OperationName: OperationAddUserSettings,
	}

	userIdStr := users_data.GetUserId()
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		return err
	}
	userEmail := users_data.GetEmail()
	if userEmail == "" {
		return ErrorNoUserEmail
	}

	userInfo := &models.UserNotificationInfo{
		UserUuid:   userId,
		Email:      userEmail,
		EmailUsage: true,
	}
	if err = gh.useCase.AddUserSettings(ctxWithTrace, userInfo); err != nil {
		error_answer.Info = err.Error()
		error_answer.Status = GetErrorCode(err)
		flag_error = true
	} else {
		answer.Info = "Ok"
	}

	if flag_error {
		answer_data, err = proto.Marshal(error_answer)
	} else {
		answer_data, err = proto.Marshal(answer)
	}

	if err != nil {
		gh.accLog.Error(err)
		return err
	}

	err = kProducer.ProduceRecord(answer_topic, sarama.ByteEncoder(answer_data))
	if err != nil {
		gh.accLog.Error(err)
		return err
	}
	gh.accLog.Info("Success to send answer!")

	return nil
}

func (gh NotificationGrpcHandlers) RemoveUserSettings(ctx context.Context, saga_uuid string, event_uuid string, users_data *api.AdditionalInfo, kProducer *kafka.ProducerProvider) error {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "NotificationGrpcHandlers.AddUserSettings")
	defer span.Finish()

	var answer_data []byte
	flag_error := false
	var err error
	answer_topic := TopicResult

	answer := &api.EventSuccess{
		SagaUuid:      saga_uuid,
		EventUuid:     event_uuid,
		OperationName: OperationRemoveUserSettings,
	}

	error_answer := &api.EventError{
		SagaUuid:      saga_uuid,
		EventUuid:     event_uuid,
		OperationName: OperationRemoveUserSettings,
	}

	userIdStr := users_data.GetUserId()
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		return err
	}

	if err = gh.useCase.DeleteUserSettings(ctxWithTrace, userId); err != nil {
		error_answer.Info = err.Error()
		error_answer.Status = GetErrorCode(err)
		flag_error = true
	} else {
		answer.Info = "Ok"
	}

	if flag_error {
		answer_data, err = proto.Marshal(error_answer)
	} else {
		answer_data, err = proto.Marshal(answer)
	}

	if err != nil {
		gh.accLog.Error(err)
		return err
	}

	err = kProducer.ProduceRecord(answer_topic, sarama.ByteEncoder(answer_data))
	if err != nil {
		gh.accLog.Error(err)
		return err
	}
	gh.accLog.Info("Success to send answer!")

	return nil
}
