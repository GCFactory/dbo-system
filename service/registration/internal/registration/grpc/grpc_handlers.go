package grpc

import (
	"context"
	"github.com/GCFactory/dbo-system/platform/pkg/logger"
	"github.com/GCFactory/dbo-system/service/registration/config"
	"github.com/GCFactory/dbo-system/service/registration/internal/models"
	"github.com/GCFactory/dbo-system/service/registration/internal/registration"
	"github.com/GCFactory/dbo-system/service/registration/internal/registration/usecase"
	"github.com/GCFactory/dbo-system/service/registration/pkg/kafka"
)

type GRPCRegistrationHandlers struct {
	cfg            *config.Config
	kProducer      *kafka.ProducerProvider
	registrationUC registration.UseCase
	accLog         logger.Logger
}

func (h *GRPCRegistrationHandlers) StartOperation(ctx context.Context, operation_type uint8, operation_data map[string]interface{}) (err error) {
	// TODO: complete

	var list_of_events []*models.Event
	switch operation_type {
	case usecase.OperationCreateUser:
		{
			list_of_events, err = h.registrationUC.StartOperation(ctx, operation_type, operation_data)
			if err != nil {
				return err
			}
			if len(list_of_events) == 0 {
				return ErrorEmptyStartEventList
			}

		}
	default:
		{
			h.accLog.Debug("Unknown type of operation!")
		}
	}

	return nil
}

func (h *GRPCRegistrationHandlers) Process(ctx context.Context, event_type string, event_data []byte) (err error) {
	// TODO: complete
	return nil
}

func NewRegistrationGRPCHandlers(cfg *config.Config, kProducer *kafka.ProducerProvider, registrationUC registration.UseCase, accLog logger.Logger) registration.RegistrationGRPCHandlers {
	return &GRPCRegistrationHandlers{cfg: cfg, kProducer: kProducer, registrationUC: registrationUC, accLog: accLog}
}
