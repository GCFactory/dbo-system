//go:generate mockgen -source pg_repository.go -destination mock/pg_repository_mock.go -package mock
package registration

import (
	"context"
	"github.com/GCFactory/dbo-system/service/registration/internal/models"
	"github.com/google/uuid"
)

type Repository interface {
	CreateSaga(ctx context.Context, saga models.Saga) error
	GetSaga(ctx context.Context, id uuid.UUID) (*models.Saga, error)
	DeleteSaga(ctx context.Context, saga_uuid uuid.UUID) error
	UpdateSaga(ctx context.Context, saga *models.Saga) error
	CreateSagaConnection(ctx context.Context, sagaConnection *models.SagaConnection) error
	GetSagaConnectionsCurrentSaga(ctx context.Context, current_saga_uuid uuid.UUID) (*models.ListOfSagaConnections, error)
	GetSagaConnectionsNextSaga(ctx context.Context, next_saga_uuid uuid.UUID) (*models.ListOfSagaConnections, error)
	UpdateSagaConnection(ctx context.Context, sagaConnection *models.SagaConnection) error
	DeleteSagaConnection(ctx context.Context, sagaConnection *models.SagaConnection) error
	CreateEvent(ctx context.Context, event *models.Event) error
	GetEvent(ctx context.Context, id uuid.UUID) (*models.Event, error)
	DeleteEvent(ctx context.Context, event_uuid uuid.UUID) error
	UpdateEvent(ctx context.Context, event *models.Event) error
	GetListOfSagaEvents(ctx context.Context, saga_uuid uuid.UUID) (*models.SagaListEvents, error)
}
