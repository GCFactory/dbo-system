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
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationRepo.CreateSaga")
	defer span.Finish()

	var res models.Saga

	if err := repo.db.QueryRowxContext(ctxWithTrace,
		CreateSaga,
		&saga.Saga_uuid,
		&saga.Saga_status,
		&saga.Saga_type,
		&saga.Saga_name,
		&saga.Saga_list_of_events,
	).StructScan(&res); err != nil {
		return ErrorCreateSaga
	}

	return nil
}

func (repo registrationRepo) GetSagaById(ctx context.Context, saga_uuid uuid.UUID) (*models.Saga, error) {
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationRepo.GetSagaById")
	defer span.Finish()

	var saga models.Saga

	if err := repo.db.QueryRowxContext(ctxWithTrace,
		GetSaga,
		saga_uuid,
	).StructScan(&saga); err != nil {
		return nil, ErrorGetSaga
	}
	return &saga, nil
}

func (repo registrationRepo) UpdateSagaStatus(ctx context.Context, saga_uuid uuid.UUID, saga_status uint) error {
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationRepo.UpdateSagaStatus")
	defer span.Finish()

	_, err := repo.db.ExecContext(ctxWithTrace,
		UpdateSagaStatus,
		saga_uuid,
		saga_status,
	)

	if err != nil {
		return ErrorUpdateSagaStatus
	}

	return nil
}

func (repo registrationRepo) UpdateSagaEvents(ctx context.Context, saga_uuid uuid.UUID, events string) error {
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationRepo.UpdateSagaEvents")
	defer span.Finish()

	_, err := repo.db.ExecContext(ctxWithTrace,
		UpdateSagaEvents,
		saga_uuid,
		events,
	)

	if err != nil {
		return ErrorUpdateSagaEvents
	}

	return nil
}

func NewRegistrationRepository(db *sqlx.DB) registration.Repository {
	return &registrationRepo{db: db}
}
