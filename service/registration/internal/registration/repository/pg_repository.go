package repository

import (
	"github.com/GCFactory/dbo-system/service/registration/internal/models"
	"github.com/GCFactory/dbo-system/service/registration/internal/registration"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
	"golang.org/x/net/context"
)

type registrationRepo struct {
	db *sqlx.DB
}

func (repo registrationRepo) CreateSaga(ctx context.Context, saga models.Saga) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "registrationRepo.CreateSaga")
	defer span.Finish()

	_, err := repo.db.ExecContext(ctx,
		createSaga,
		&saga.Saga_uuid,
		&saga.Saga_status,
		&saga.Saga_type,
		&saga.Saga_name,
		&saga.Saga_list_of_events,
	)

	if err != nil {
		return ErrorCreateSaga
	}

	return nil
}

func (repo registrationRepo) ChangeSagaStatus(ctx context.Context, saga_uuid uuid.UUID, saga_status uint) error {
	return nil
}

func (repo registrationRepo) RemoveEventFromSaga(ctx context.Context, saga_uuid uuid.UUID, event_name string) error {
	return nil
}

func (repo registrationRepo) AddEventInToSaga(ctx context.Context, saga_uuid uuid.UUID, event_name string) error {
	return nil
}

func (repo registrationRepo) GetSagaById(ctx context.Context, saga_uuid uuid.UUID) (models.Saga, error) {
	return models.Saga{}, nil
}

func NewRegistrationRepository(db *sqlx.DB) registration.Repository {
	return &registrationRepo{db: db}
}
