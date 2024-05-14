//go:generate mockgen -source pg_repository.go -destination mock/pg_repository_mock.go -package mock
package registration

import (
	"context"
	"github.com/GCFactory/dbo-system/service/registration/internal/models"
	"github.com/google/uuid"
)

type Repository interface {
	CreateSaga(ctx context.Context, saga models.Saga) error
	GetSagaById(ctx context.Context, saga_uuid uuid.UUID) (*models.Saga, error)
	UpdateSagaStatus(ctx context.Context, saga_uuid uuid.UUID, saga_status uint) error
	UpdateSagaEvents(ctx context.Context, saga_uuid uuid.UUID, events string) error
}
