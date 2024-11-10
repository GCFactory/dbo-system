package registration

import (
	"context"
	"github.com/GCFactory/dbo-system/service/registration/internal/models"
	"github.com/google/uuid"
)

type UseCase interface {
	StartOperation(ctx context.Context, operation_type uint8, start_data map[string]interface{}) ([]*models.Event, error)
	ProcessingSagaAndEvents(ctx context.Context, saga_uuid uuid.UUID, event_uuid uuid.UUID, success bool, data map[string]interface{}) ([]*models.Event, error)
	GetEventData(ctx context.Context, event_uuid uuid.UUID) (map[string]interface{}, error)
	SetOrUpdateSagaData(ctx context.Context, saga_uuid uuid.UUID, new_saga_data map[string]interface{}) error
}
