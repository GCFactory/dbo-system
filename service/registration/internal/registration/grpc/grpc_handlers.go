package grpc

import (
	"context"
	"github.com/GCFactory/dbo-system/platform/pkg/logger"
	"github.com/GCFactory/dbo-system/service/registration/config"
	"github.com/GCFactory/dbo-system/service/registration/internal/registration"
	"github.com/GCFactory/dbo-system/service/registration/internal/registration/usecase"
	"github.com/GCFactory/dbo-system/service/registration/pkg/kafka"
)

type RegistrationHandlers struct {
	cfg            *config.Config
	kProducer      *kafka.ProducerProvider
	registrationUC registration.UseCase
	accLog         logger.Logger
}

func (h *RegistrationHandlers) StartOperation(ctx context.Context, operation_type uint8, operation_data []byte) (err error) {
	// TODO: complete

	switch operation_type {
	case usecase.OperationCreateUser:
		{
			list_of_events, err := h.registrationUC.StartOperation(ctx, operation_type)
			if err != nil {
				return err
			}
			for _, event := range list_of_events {
				
			}

		}
	default:
		{
			h.accLog.Debug("Unknown type of operation!")
		}
	}

	return nil
}

func (h *RegistrationHandlers) Process(ctx context.Context, event_type string, event_data []byte) (err error) {
	// TODO: complete
	return nil
}

func NewRegistrationGRPCHandlers(cfg *config.Config, kProducer *kafka.ProducerProvider, registrationUC registration.UseCase, accLog logger.Logger) registration.RegistrationGRPCHandlers {
	return &RegistrationHandlers{cfg: cfg, kProducer: kProducer, registrationUC: registrationUC, accLog: accLog}
}
