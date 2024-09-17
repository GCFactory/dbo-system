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

	_, err := repo.db.ExecContext(ctxWithTrace,
		CreateSaga,
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

func (repo registrationRepo) DeleteSaga(ctx context.Context, saga_uuid uuid.UUID) error {
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationRepo.DeleteSaga")
	defer span.Finish()

	result, err := repo.db.ExecContext(ctxWithTrace,
		DeleteSaga,
		saga_uuid,
	)

	if err != nil {
		return ErrorDeleteSaga
	} else if count, err := result.RowsAffected(); err != nil || count == 0 {
		return ErrorDeleteSaga
	}

	return nil
}

func (repo registrationRepo) GetSaga(ctx context.Context, id uuid.UUID) (*models.Saga, error) {
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationRepo.GetSaga")
	defer span.Finish()

	var result *models.Saga

	if err := repo.db.QueryRowContext(ctxWithTrace,
		GetSaga,
		id).Scan(&result); err != nil {
		return nil, ErrorGetSaga
	}

	return result, nil
}

func (repo registrationRepo) UpdateSaga(ctx context.Context, saga *models.Saga) error {
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationRepo.UpdateSaga")
	defer span.Finish()

	result, err := repo.db.ExecContext(ctxWithTrace,
		UpdateSaga,
		&saga.Saga_uuid,
		&saga.Saga_status,
		&saga.Saga_type,
		&saga.Saga_name,
		&saga.Saga_list_of_events,
	)

	if err != nil {
		return ErrorUpdateSaga
	} else if count, err := result.RowsAffected(); err != nil || count == 0 {
		return ErrorUpdateSaga
	}

	return nil
}

func (repo registrationRepo) CreateSagaConnection(ctx context.Context, sagaConnection *models.SagaConnection) error {
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationRepo.CreateSagaConnection")
	defer span.Finish()

	_, err := repo.db.ExecContext(ctxWithTrace,
		CreateSagaConnection,
		&sagaConnection.Current_saga_uuid,
		&sagaConnection.Next_saga_uuid,
	)

	if err != nil {
		return ErrorCreateSagaConnection
	}

	return nil
}

func (repo registrationRepo) GetSagaConnections(ctx context.Context, saga_uuid uuid.UUID) (*models.ListOfSagaConnections, error) {
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationRepo.GetSagaConnections")
	defer span.Finish()

	var result *models.ListOfSagaConnections

	if err := repo.db.QueryRowxContext(ctxWithTrace,
		GetSagaConnections,
		&saga_uuid).StructScan(&result); err != nil {
		return nil, ErrorGetSagaConnections
	}

	return result, nil
}

func (repo registrationRepo) DeleteSagaConnection(ctx context.Context, sagaConnection *models.SagaConnection) error {
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationRepo.DeleteSagaConnection")
	defer span.Finish()

	result, err := repo.db.ExecContext(ctxWithTrace,
		DeleteSagaConnection,
		&sagaConnection.Current_saga_uuid,
		&sagaConnection.Next_saga_uuid)

	if err != nil {
		return ErrorDeleteSagaConnection
	} else if count, err := result.RowsAffected(); err != nil || count == 0 {
		return ErrorDeleteSagaConnection
	}

	return nil
}

func (repo registrationRepo) CreateEvent(ctx context.Context, event *models.Event) error {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationRepo.CreateEvent")
	defer span.Finish()

	if _, err := repo.db.ExecContext(ctxWithTrace,
		CreateEvent,
		&event.Event_uuid,
		&event.Event_status,
		&event.Event_name,
		&event.Event_result,
		&event.Event_rollback_uuid,
	); err != nil {
		return ErrorCreateEvent
	}

	return nil
}

func (repo registrationRepo) GetEvent(ctx context.Context, id uuid.UUID) (*models.Event, error) {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationRepo.GetEvent")
	defer span.Finish()

	var result *models.Event

	if err := repo.db.QueryRowxContext(ctxWithTrace,
		GetEvent,
		id).StructScan(&result); err != nil {
		return nil, ErrorGetEvent
	}

	return result, nil
}

func (repo registrationRepo) DeleteEvent(ctx context.Context, event_uuid uuid.UUID) error {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationRepo.DeleteEvent")
	defer span.Finish()

	result, err := repo.db.ExecContext(ctxWithTrace,
		DeleteEvent,
		event_uuid)

	if err != nil {
		return ErrorDeleteEvent
	} else if count, err := result.RowsAffected(); err != nil || count == 0 {
		return ErrorDeleteEvent
	}

	return nil
}

func (repo registrationRepo) UpdateEvent(ctx context.Context, event *models.Event) error {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationRepo.UpdateEvent")
	defer span.Finish()

	result, err := repo.db.ExecContext(ctxWithTrace,
		UpdateEvent,
		&event.Event_uuid,
		&event.Event_status,
		&event.Event_name,
		&event.Event_result,
		&event.Event_rollback_uuid,
	)

	if err != nil {
		return ErrorUpdateEvent
	} else if count, err := result.RowsAffected(); err != nil || count == 0 {
		return ErrorUpdateEvent
	}

	return nil
}

func (repo registrationRepo) GetListOfSagaEvents(ctx context.Context, saga_uuid uuid.UUID) (*models.SagaListEvents, error) {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationRepo.GetListOfSagaEvents")
	defer span.Finish()

	var result *models.SagaListEvents

	if err := repo.db.QueryRowxContext(ctxWithTrace,
		GetListOfSagaEvents,
		saga_uuid,
	).StructScan(&result); err != nil {
		return nil, ErrorGetListOfSagaEvents
	}

	return result, nil
}

func NewRegistrationRepository(db *sqlx.DB) registration.Repository {
	return &registrationRepo{db: db}
}
