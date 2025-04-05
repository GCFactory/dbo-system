//go:generate mockgen -source pg_repository.go -destination mock/pg_repository_mock.go -package mock
package registration

import (
	"context"
	"github.com/GCFactory/dbo-system/service/registration/internal/models"
	"github.com/google/uuid"
	"time"
)

type Repository interface {
	CreateOperation(ctx context.Context, operation *models.Operation) error
	GetOperation(ctx context.Context, id uuid.UUID) (*models.Operation, error)
	UpdateOperation(ctx context.Context, operation *models.Operation) error
	DeleteOperation(ctx context.Context, id uuid.UUID) error
	CreateSaga(ctx context.Context, saga *models.Saga) error
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
	GetRevertEvent(ctx context.Context, event_uuid uuid.UUID) (*models.Event, error)
	GetOperationSaga(ctx context.Context, operation_uuid uuid.UUID) (*models.ListOfSaga, error)
	GetOperationBetweenInterval(ctx context.Context, begin time.Time, end time.Time) ([]uuid.UUID, error)
}
