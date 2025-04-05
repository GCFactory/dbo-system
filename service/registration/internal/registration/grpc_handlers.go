package registration

import (
	"context"
	"github.com/GCFactory/dbo-system/service/registration/internal/models"
	"github.com/google/uuid"
)

type RegistrationGRPCHandlers interface {
	StartOperation(ctx context.Context, operation_type uint8, operation_data map[string]interface{}) (uuid.UUID, error)
	Process(ctx context.Context, saga_uuid uuid.UUID, saga *models.Saga, event_uuid uuid.UUID, event *models.Event, data map[string]interface{}, is_success bool) error
	SendRequest(ctx context.Context, server uint8, operation_name string, saga_uuid uuid.UUID, event_uuid uuid.UUID, data map[string]interface{}) error
}
