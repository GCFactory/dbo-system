package grpc

import (
	"context"
	"github.com/GCFactory/dbo-system/platform/pkg/logger"
	"github.com/GCFactory/dbo-system/service/registration/config"
	"github.com/GCFactory/dbo-system/service/registration/gen_proto/proto/api"
	"github.com/GCFactory/dbo-system/service/registration/gen_proto/proto/platform"
	"github.com/GCFactory/dbo-system/service/registration/internal/models"
	"github.com/GCFactory/dbo-system/service/registration/internal/registration"
	"github.com/GCFactory/dbo-system/service/registration/internal/registration/usecase"
	"github.com/GCFactory/dbo-system/service/registration/pkg/kafka"
	"github.com/IBM/sarama"
	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type GRPCRegistrationHandlers struct {
	cfg            *config.Config
	kProducer      *kafka.ProducerProvider
	registrationUC registration.UseCase
	regLog         logger.Logger
}

func (h *GRPCRegistrationHandlers) StartOperation(ctx context.Context, operation_type uint8, operation_data map[string]interface{}) (err error) {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "GRPCRegistrationHandlers.StartOperation")
	defer span.Finish()

	var list_of_events []*models.Event
	switch operation_type {
	case usecase.OperationCreateUser:
		{
			list_of_events, err = h.registrationUC.StartOperation(ctxWithTrace, operation_type, operation_data)
			if err != nil {
				return err
			}
			if len(list_of_events) == 0 {
				return ErrorEmptyStartEventList
			}

		}
	default:
		{
			h.regLog.Debug("Unknown type of operation!")
		}
	}

	for _, event := range list_of_events {
		err = h.Process(ctxWithTrace, event.Saga_uuid, nil, event.Event_uuid, event, operation_data, true)
		if err != nil {
			h.regLog.Error(err)
		}
	}
	return nil
}

func (h *GRPCRegistrationHandlers) Process(ctx context.Context, saga_uuid uuid.UUID, saga *models.Saga, event_uuid uuid.UUID, event *models.Event, data map[string]interface{}, is_success bool) (err error) {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "GRPCRegistrationHandlers.Process")
	defer span.Finish()

	if uuid.Nil == saga_uuid || uuid.Nil == event_uuid {
		return ErrorInvalidUUIDs
	}

	if event != nil {
		err = h.SendRequest(ctxWithTrace, GetServerByOperation(event.Event_name), event.Event_name, saga_uuid, event_uuid, data)
		if err != nil {
			h.regLog.Error(err)
			return err
		}
	} else {
		list_of_events, local_err := h.registrationUC.ProcessingSagaAndEvents(ctxWithTrace,
			saga_uuid,
			event_uuid,
			is_success,
			data,
		)
		if local_err != nil {
			h.regLog.Error(local_err)
			return local_err
		}
		for _, new_event := range list_of_events {
			err = h.Process(ctxWithTrace,
				new_event.Saga_uuid,
				nil,
				new_event.Event_uuid,
				new_event,
				data,
				true)
			if err != nil {
				h.regLog.Error(err)
				return err
			}
		}
	}
	return nil
}

func (h *GRPCRegistrationHandlers) SendRequest(ctx context.Context, server uint8, operation_name string, saga_uuid uuid.UUID, event_uuid uuid.UUID, data map[string]interface{}) (err error) {

	span, _ := opentracing.StartSpanFromContext(ctx, "GRPCRegistrationHandlers.SendRequest")
	defer span.Finish()

	if !ValidateServer(server) {
		return ErrorInvalidServer
	} else if !ValidateServersOperation(server, operation_name) {
		return ErrorInvalidServersOperation
	}

	switch server {
	case ServerTypeUsers:
		{
			//  TODO: complete
			switch operation_name {
			case OperationCreateUser:
				{

					if !ValidateOperationsData(operation_name, data) {
						return ErrorInvalidOperationsData
					}

					if !ValidateServerTopic(server, ServerTopicUsersConsumer) {
						return ErrorInvalidServersTopic
					}

					birth_date, err := time.Parse("02-01-2006 15:04:05", data["birth_date"].(string))
					if err != nil {
						return ErrorConvertTimestamp
					}
					birth_date_proto := timestamppb.New(birth_date)

					authority_date, err := time.Parse("02-01-2006 15:04:05", data["authority_date"].(string))
					if err != nil {
						return ErrorConvertTimestamp
					}
					authority_date_proto := timestamppb.New(authority_date)

					kafka_data := &api.EventData{
						SagaUuid:      saga_uuid.String(),
						EventUuid:     event_uuid.String(),
						OperationName: OperationCreateUser,
						Data: &api.EventData_UserInfo{
							UserInfo: &api.UserInfo{
								UserInn: data["user_inn"].(string),
								Passport: &platform.Passport{
									Number: data["passport_number"].(string),
									Series: data["passport_series"].(string),
									Fcs: &platform.FCs{
										Name:       data["name"].(string),
										Surname:    data["surname"].(string),
										Patronymic: data["patronimic"].(string),
									},
									BirthDate:          birth_date_proto,
									BirthLocation:      data["birth_location"].(string),
									PickUpPoint:        data["pick_up_point"].(string),
									Authority:          data["authority"].(string),
									AuthorityDate:      authority_date_proto,
									RegistrationAdress: data["registration_adress"].(string),
								},
							},
						},
					}

					msg, err := proto.Marshal(kafka_data)
					if err != nil {
						h.regLog.Error(err)
						return err
					}
					err = h.kProducer.ProduceRecord(ServerTopicUsersConsumer, sarama.ByteEncoder(msg))
					if err != nil {
						h.regLog.Error(err)
						return err
					}
					break
				}
			default:
				{
					return ErrorInvalidServersOperation
				}
			}
			break
		}
	case ServerTypeAccounts:
		{
			//  TODO: complete
			break
		}
	default:
		{
			return ErrorInvalidServer
		}
	}

	return nil
}

func NewRegistrationGRPCHandlers(cfg *config.Config, kProducer *kafka.ProducerProvider, registrationUC registration.UseCase, regLog logger.Logger) registration.RegistrationGRPCHandlers {
	return &GRPCRegistrationHandlers{cfg: cfg, kProducer: kProducer, registrationUC: registrationUC, regLog: regLog}
}
