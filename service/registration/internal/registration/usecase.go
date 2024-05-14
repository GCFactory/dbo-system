package registration

import (
	"github.com/GCFactory/dbo-system/service/registration/internal/models"
	"github.com/google/uuid"
	"golang.org/x/net/context"
)

type UseCase interface {
	CreateSaga(ctx context.Context, saga models.Saga) error
	AddSagaEvent(ctx context.Context, saga_uuid uuid.UUID, event_name string) error
	RemoveSagaEvent(ctx context.Context, saga_uuid uuid.UUID, event_name string) error
	FallBackSaga(ctx context.Context, saga_uuid uuid.UUID) error
}
