package repository

import (
	"encoding/json"
	"github.com/GCFactory/dbo-system/service/registration/internal/models"
	"github.com/GCFactory/dbo-system/service/registration/internal/registration"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
	"golang.org/x/net/context"
	"time"
)

type registrationRepo struct {
	db *sqlx.DB
}

func (repo registrationRepo) CreateOperation(ctx context.Context, operation *models.Operation) error {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationRepo.CreateOperation")
	defer span.Finish()

	tx, err := repo.db.BeginTx(ctxWithTrace, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	_, err = repo.db.ExecContext(ctxWithTrace,
		CreateOperation,
		operation.Operation_uuid,
		operation.Operation_name,
		operation.Create_time,
		operation.Last_time_update,
	)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (repo registrationRepo) GetOperation(ctx context.Context, id uuid.UUID) (*models.Operation, error) {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationRepo.GetOperation")
	defer span.Finish()

	tx, err := repo.db.BeginTx(ctxWithTrace, nil)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	result := &models.Operation{}

	if err := repo.db.QueryRowxContext(ctxWithTrace,
		GetOperation,
		id).StructScan(result); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return result, nil
}

func (repo registrationRepo) UpdateOperation(ctx context.Context, operation *models.Operation) error {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationRepo.UpdateOperation")
	defer span.Finish()

	tx, err := repo.db.BeginTx(ctxWithTrace, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	result, err := repo.db.ExecContext(ctxWithTrace,
		UpdateOperation,
		operation.Operation_uuid,
		operation.Last_time_update,
	)

	if err != nil {
		return err
	} else if count, err := result.RowsAffected(); err != nil || count == 0 {
		return ErrorUpdateOperation
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (repo registrationRepo) DeleteOperation(ctx context.Context, id uuid.UUID) error {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationRepo.DeleteOperation")
	defer span.Finish()

	tx, err := repo.db.BeginTx(ctxWithTrace, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	result, err := repo.db.ExecContext(ctxWithTrace,
		DeleteOperation,
		id)

	if err != nil {
		return err
	} else if count, err := result.RowsAffected(); err != nil || count == 0 {
		return ErrorDeleteOperation
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil

}

func (repo registrationRepo) CreateSaga(ctx context.Context, saga *models.Saga) error {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationRepo.CreateSaga")
	defer span.Finish()

	tx, err := repo.db.BeginTx(ctxWithTrace, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	saga_data, err := json.Marshal(saga.Saga_data)
	if err != nil {
		return ErrorCreateSaga
	}

	_, err = repo.db.ExecContext(ctxWithTrace,
		CreateSaga,
		saga.Saga_uuid,
		saga.Saga_status,
		saga.Saga_type,
		saga.Saga_name,
		saga_data,
		saga.Operation_uuid,
	)
	if err != nil {
		return ErrorCreateSaga
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (repo registrationRepo) DeleteSaga(ctx context.Context, saga_uuid uuid.UUID) error {
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationRepo.DeleteSaga")
	defer span.Finish()

	tx, err := repo.db.BeginTx(ctxWithTrace, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	result, err := repo.db.ExecContext(ctxWithTrace,
		DeleteSaga,
		saga_uuid,
	)

	if err != nil {
		return ErrorDeleteSaga
	} else if count, err := result.RowsAffected(); err != nil || count == 0 {
		return ErrorDeleteSaga
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (repo registrationRepo) GetSaga(ctx context.Context, id uuid.UUID) (*models.Saga, error) {
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationRepo.GetSaga")
	defer span.Finish()

	tx, err := repo.db.BeginTx(ctxWithTrace, nil)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	data := &models.SagaFromDB{}

	if err := repo.db.QueryRowxContext(ctxWithTrace,
		GetSaga,
		id).StructScan(data); err != nil {
		return nil, ErrorGetSaga
	}

	saga_data := make(map[string]interface{})
	err = json.Unmarshal([]byte(data.Saga_data), &saga_data)
	if err != nil {
		return nil, ErrorGetSaga
	}

	result := &models.Saga{
		Saga_uuid:      data.Saga_uuid,
		Saga_name:      data.Saga_name,
		Saga_type:      data.Saga_type,
		Saga_status:    data.Saga_status,
		Saga_data:      saga_data,
		Operation_uuid: data.Operation_uuid,
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return result, nil
}

func (repo registrationRepo) UpdateSaga(ctx context.Context, saga *models.Saga) error {
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationRepo.UpdateSaga")
	defer span.Finish()

	tx, err := repo.db.BeginTx(ctxWithTrace, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	saga_data, err := json.Marshal(saga.Saga_data)
	if err != nil {
		return ErrorUpdateSaga
	}

	result, err := repo.db.ExecContext(ctxWithTrace,
		UpdateSaga,
		&saga.Saga_uuid,
		&saga.Saga_status,
		&saga.Saga_type,
		&saga.Saga_name,
		saga_data,
	)

	if err != nil {
		return ErrorUpdateSaga
	} else if count, err := result.RowsAffected(); err != nil || count == 0 {
		return ErrorUpdateSaga
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (repo registrationRepo) CreateSagaConnection(ctx context.Context, sagaConnection *models.SagaConnection) error {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationRepo.CreateSagaConnection")
	defer span.Finish()

	tx, err := repo.db.BeginTx(ctxWithTrace, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	_, err = repo.db.ExecContext(ctxWithTrace,
		CreateSagaConnection,
		&sagaConnection.Current_saga_uuid,
		&sagaConnection.Next_saga_uuid,
		&sagaConnection.Acc_connection_status,
	)

	if err != nil {
		return ErrorCreateSagaConnection
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (repo registrationRepo) GetSagaConnectionsCurrentSaga(ctx context.Context, current_saga_uuid uuid.UUID) (*models.ListOfSagaConnections, error) {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationRepo.GetSagaConnectionsCurrentSaga")
	defer span.Finish()

	tx, err := repo.db.BeginTx(ctxWithTrace, nil)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	rows, err := repo.db.QueryxContext(
		ctxWithTrace,
		GetSagaConnectionsCurrentSaga,
		&current_saga_uuid,
	)

	if err != nil {
		return nil, ErrorGetSagaCurrentConnections
	}

	result := &models.ListOfSagaConnections{}

	for rows.Next() {

		saga_connection := &models.SagaConnection{}

		if err := rows.StructScan(saga_connection); err != nil {
			return nil, ErrorGetSagaCurrentConnections
		}

		result.List_of_connetcions = append(result.List_of_connetcions, saga_connection)

	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return result, nil
}

func (repo registrationRepo) GetSagaConnectionsNextSaga(ctx context.Context, next_saga_uuid uuid.UUID) (*models.ListOfSagaConnections, error) {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationRepo.GetSagaConnectionsNextSaga")
	defer span.Finish()

	tx, err := repo.db.BeginTx(ctxWithTrace, nil)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	rows, err := repo.db.QueryxContext(
		ctxWithTrace,
		GetSagaConnectionsNextSaga,
		&next_saga_uuid,
	)

	if err != nil {
		return nil, ErrorGetSagaNextConnections
	}

	result := &models.ListOfSagaConnections{}

	for rows.Next() {

		saga_connection := &models.SagaConnection{}

		if err := rows.StructScan(saga_connection); err != nil {
			return nil, ErrorGetSagaNextConnections
		}

		result.List_of_connetcions = append(result.List_of_connetcions, saga_connection)

	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return result, nil
}

func (repo registrationRepo) DeleteSagaConnection(ctx context.Context, sagaConnection *models.SagaConnection) error {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationRepo.DeleteSagaConnection")
	defer span.Finish()

	tx, err := repo.db.BeginTx(ctxWithTrace, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	result, err := repo.db.ExecContext(ctxWithTrace,
		DeleteSagaConnection,
		&sagaConnection.Current_saga_uuid,
		&sagaConnection.Next_saga_uuid)

	if err != nil {
		return ErrorDeleteSagaConnection
	} else if count, err := result.RowsAffected(); err != nil || count == 0 {
		return ErrorDeleteSagaConnection
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (repo registrationRepo) UpdateSagaConnection(ctx context.Context, sagaConnection *models.SagaConnection) error {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationRepo.UpdateSagaConnection")
	defer span.Finish()

	tx, err := repo.db.BeginTx(ctxWithTrace, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	result, err := repo.db.ExecContext(ctxWithTrace,
		UpdateSagaConnection,
		&sagaConnection.Current_saga_uuid,
		&sagaConnection.Next_saga_uuid,
		&sagaConnection.Acc_connection_status,
	)

	if err != nil {
		return ErrorUpdateSagaConnection
	} else if count, err := result.RowsAffected(); err != nil || count == 0 {
		return ErrorUpdateSagaConnection
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil

}

func (repo registrationRepo) CreateEvent(ctx context.Context, event *models.Event) error {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationRepo.CreateEvent")
	defer span.Finish()

	tx, err := repo.db.BeginTx(ctxWithTrace, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	event_data, err := json.Marshal(event.Event_required_data)
	if err != nil {
		return ErrorCreateEvent
	}

	if _, err := repo.db.ExecContext(ctxWithTrace,
		CreateEvent,
		&event.Event_uuid,
		&event.Saga_uuid,
		&event.Event_status,
		&event.Event_name,
		&event.Event_is_roll_back,
		event_data,
		&event.Event_result,
		&event.Event_rollback_uuid,
	); err != nil {
		return ErrorCreateEvent
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (repo registrationRepo) GetEvent(ctx context.Context, id uuid.UUID) (*models.Event, error) {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationRepo.GetEvent")
	defer span.Finish()

	tx, err := repo.db.BeginTx(ctxWithTrace, nil)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	event_data := &models.EventFromDB{}

	if err := repo.db.QueryRowxContext(ctxWithTrace,
		GetEvent,
		id).StructScan(event_data); err != nil {
		return nil, ErrorGetEvent
	}

	var event_data_arr []string
	err = json.Unmarshal([]byte(event_data.Event_required_data), &event_data_arr)
	if err != nil {
		return nil, ErrorGetEvent
	}
	result := &models.Event{
		Event_uuid:          event_data.Event_uuid,
		Event_name:          event_data.Event_name,
		Event_is_roll_back:  event_data.Event_is_roll_back,
		Event_status:        event_data.Event_status,
		Event_result:        event_data.Event_result,
		Saga_uuid:           event_data.Saga_uuid,
		Event_rollback_uuid: event_data.Event_rollback_uuid,
		Event_required_data: event_data_arr,
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return result, nil
}

func (repo registrationRepo) DeleteEvent(ctx context.Context, event_uuid uuid.UUID) error {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationRepo.DeleteEvent")
	defer span.Finish()

	tx, err := repo.db.BeginTx(ctxWithTrace, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	result, err := repo.db.ExecContext(ctxWithTrace,
		DeleteEvent,
		event_uuid)

	if err != nil {
		return ErrorDeleteEvent
	} else if count, err := result.RowsAffected(); err != nil || count == 0 {
		return ErrorDeleteEvent
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (repo registrationRepo) UpdateEvent(ctx context.Context, event *models.Event) error {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationRepo.UpdateEvent")
	defer span.Finish()

	tx, err := repo.db.BeginTx(ctxWithTrace, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	event_data, err := json.Marshal(event.Event_required_data)
	if err != nil {
		return ErrorUpdateEvent
	}

	result, err := repo.db.ExecContext(ctxWithTrace,
		UpdateEvent,
		&event.Event_uuid,
		&event.Event_status,
		event_data,
		&event.Event_result,
		&event.Event_rollback_uuid,
	)

	if err != nil {
		return ErrorUpdateEvent
	} else if count, err := result.RowsAffected(); err != nil || count == 0 {
		return ErrorUpdateEvent
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (repo registrationRepo) GetListOfSagaEvents(ctx context.Context, saga_uuid uuid.UUID) (*models.SagaListEvents, error) {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationRepo.GetListOfSagaEvents")
	defer span.Finish()

	tx, err := repo.db.BeginTx(ctxWithTrace, nil)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	rows, err := repo.db.QueryxContext(
		ctxWithTrace,
		GetListOfSagaEvents,
		saga_uuid,
	)

	if err != nil {
		return nil, ErrorGetListOfSagaEvents
	}

	result := &models.SagaListEvents{}

	for rows.Next() {
		var event_uuid uuid.UUID

		if err := rows.Scan(&event_uuid); err != nil {
			return nil, ErrorGetListOfSagaEvents
		}

		result.EventList = append(result.EventList, event_uuid)
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return result, nil
}

func (repo registrationRepo) GetRevertEvent(ctx context.Context, event_uuid uuid.UUID) (*models.Event, error) {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationRepo.GetRevertEvent")
	defer span.Finish()

	tx, err := repo.db.BeginTx(ctxWithTrace, nil)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	event_data := &models.EventFromDB{}

	if err := repo.db.QueryRowxContext(ctxWithTrace,
		GetRevertEvent,
		event_uuid).StructScan(event_data); err != nil {
		return nil, ErrorGetEvent
	}

	var event_data_arr []string
	err = json.Unmarshal([]byte(event_data.Event_required_data), &event_data_arr)
	if err != nil {
		return nil, ErrorGetEvent
	}
	result := &models.Event{
		Event_uuid:          event_data.Event_uuid,
		Event_name:          event_data.Event_name,
		Event_is_roll_back:  event_data.Event_is_roll_back,
		Event_status:        event_data.Event_status,
		Event_result:        event_data.Event_result,
		Saga_uuid:           event_data.Saga_uuid,
		Event_rollback_uuid: event_data.Event_rollback_uuid,
		Event_required_data: event_data_arr,
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return result, nil
}

func (repo registrationRepo) GetOperationSaga(ctx context.Context, operation_uuid uuid.UUID) (*models.ListOfSaga, error) {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationRepo.GetOperationSaga")
	defer span.Finish()

	tx, err := repo.db.BeginTx(ctxWithTrace, nil)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	rows, err := repo.db.QueryxContext(
		ctxWithTrace,
		GetListOfOperationSagas,
		operation_uuid,
	)

	if err != nil {
		return nil, err
	}

	result := &models.ListOfSaga{}

	for rows.Next() {
		var saga_uuid uuid.UUID

		if err := rows.Scan(&saga_uuid); err != nil {
			return nil, err
		}

		result.ListId = append(result.ListId, saga_uuid)
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return result, nil
}

func (repo registrationRepo) GetOperationBetweenInterval(ctx context.Context, begin time.Time, end time.Time) ([]uuid.UUID, error) {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationRepo.GetOperationSaga")
	defer span.Finish()

	tx, err := repo.db.BeginTx(ctxWithTrace, nil)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()
	
	rows, err := repo.db.QueryxContext(
		ctxWithTrace,
		GetListOperationsBetweenInterval,
		begin,
		end,
	)

	if err != nil {
		return nil, err
	}

	result := make([]uuid.UUID, 0)

	for rows.Next() {
		var operation_uuid uuid.UUID

		if err := rows.Scan(&operation_uuid); err != nil {
			return nil, err
		}

		result = append(result, operation_uuid)
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return result, nil
}

func NewRegistrationRepository(db *sqlx.DB) registration.Repository {
	return &registrationRepo{db: db}
}
