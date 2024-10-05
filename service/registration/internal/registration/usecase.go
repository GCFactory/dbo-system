package registration

import (
	"context"
	"github.com/GCFactory/dbo-system/service/registration/internal/models"
	"github.com/google/uuid"
)

type UseCase interface {
	StartOperation(ctx context.Context, operation_type uint8) ([]*models.Event, error)
	ProcessingSagaAndEvents(ctx context.Context, saga_uuid uuid.UUID, event_uuid uuid.UUID, success bool) ([]*models.Event, error)
}
