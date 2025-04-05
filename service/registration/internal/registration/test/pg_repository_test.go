package test

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/GCFactory/dbo-system/platform/config"
	"github.com/GCFactory/dbo-system/service/registration/internal/models"
	"github.com/GCFactory/dbo-system/service/registration/internal/registration/repository"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"testing"
)

var (
	testCfg = &config.Config{
		Env: "Development",
		Logger: config.Logger{
			Development: true,
			Level:       "Debug",
		},
	}
)

func TestRepository_CreateSaga(t *testing.T) {

	t.Parallel()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "slqmock")
	defer sqlxDB.Close()

	regRepo := repository.NewRegistrationRepository(sqlxDB)

	saga := &models.Saga{
		Saga_uuid:   uuid.New(),
		Saga_name:   "test_name",
		Saga_status: 0,
		Saga_type:   1,
	}

	t.Run("Success", func(t *testing.T) {

		mock.ExpectBegin()
		mock.ExpectExec(repository.CreateSaga).
			WithArgs(
				&saga.Saga_uuid,
				&saga.Saga_status,
				&saga.Saga_type,
				&saga.Saga_name,
			).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err = regRepo.CreateSaga(context.Background(), saga)
		require.Nil(t, err)

		err = mock.ExpectationsWereMet()
		require.Nil(t, err)

	})

	t.Run("Error", func(t *testing.T) {

		mock.ExpectBegin()

		mock.ExpectExec(repository.CreateSaga).WithArgs(
			&saga.Saga_uuid,
			&saga.Saga_status,
			&saga.Saga_type,
			&saga.Saga_name,
		).WillReturnError(errors.New("test error"))

		mock.ExpectRollback()

		err = regRepo.CreateSaga(context.Background(), saga)
		require.Equal(t, err, repository.ErrorCreateSaga)

		err = mock.ExpectationsWereMet()
		require.Nil(t, err)

	})
}

func TestRepository_DeleteSaga(t *testing.T) {

	t.Parallel()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "slqmock")
	defer sqlxDB.Close()

	regRepo := repository.NewRegistrationRepository(sqlxDB)

	saga := &models.Saga{
		Saga_uuid:   uuid.New(),
		Saga_name:   "test_name",
		Saga_status: 0,
		Saga_type:   1,
	}

	t.Run("Success", func(t *testing.T) {

		mock.ExpectBegin()
		mock.ExpectExec(repository.DeleteSaga).
			WithArgs(
				&saga.Saga_uuid,
			).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err = regRepo.DeleteSaga(context.Background(), saga.Saga_uuid)
		require.Nil(t, err)

		err = mock.ExpectationsWereMet()
		require.Nil(t, err)

	})

	t.Run("Error", func(t *testing.T) {

		mock.ExpectBegin()

		mock.ExpectExec(repository.DeleteSaga).WithArgs(
			&saga.Saga_uuid,
		).WillReturnError(errors.New("test error"))

		mock.ExpectRollback()

		err = regRepo.DeleteSaga(context.Background(), saga.Saga_uuid)
		require.Equal(t, err, repository.ErrorDeleteSaga)

		err = mock.ExpectationsWereMet()
		require.Nil(t, err)

	})

	t.Run("Error no affected rows", func(t *testing.T) {

		mock.ExpectBegin()

		mock.ExpectExec(repository.DeleteSaga).WithArgs(
			&saga.Saga_uuid,
		).WillReturnResult(sqlmock.NewResult(0, 0))

		mock.ExpectRollback()

		err = regRepo.DeleteSaga(context.Background(), saga.Saga_uuid)
		require.Equal(t, err, repository.ErrorDeleteSaga)

		err = mock.ExpectationsWereMet()
		require.Nil(t, err)

	})
}

func TestRepository_GetSaga(t *testing.T) {

	t.Parallel()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "slqmock")
	defer sqlxDB.Close()

	regRepo := repository.NewRegistrationRepository(sqlxDB)

	saga := &models.Saga{
		Saga_uuid:   uuid.New(),
		Saga_name:   "test_name",
		Saga_status: 0,
		Saga_type:   1,
	}

	t.Run("Success", func(t *testing.T) {

		mock.ExpectBegin()

		rows := mock.NewRows([]string{"saga_uuid", "saga_status", "saga_type", "saga_name"}).AddRow(
			&saga.Saga_uuid,
			&saga.Saga_status,
			&saga.Saga_type,
			saga.Saga_name,
		)

		mock.ExpectQuery(repository.GetSaga).
			WithArgs(
				&saga.Saga_uuid,
			).WillReturnRows(rows)
		mock.ExpectCommit()

		result, err := regRepo.GetSaga(context.Background(), saga.Saga_uuid)
		require.Nil(t, err)
		require.Equal(t, saga, result)

		err = mock.ExpectationsWereMet()
		require.Nil(t, err)

	})

	t.Run("Error", func(t *testing.T) {

		mock.ExpectBegin()

		mock.ExpectQuery(repository.GetSaga).WithArgs(
			&saga.Saga_uuid,
		).WillReturnError(errors.New("test error"))

		mock.ExpectRollback()

		result, err := regRepo.GetSaga(context.Background(), saga.Saga_uuid)
		require.Equal(t, err, repository.ErrorGetSaga)
		require.Nil(t, result)

		err = mock.ExpectationsWereMet()
		require.Nil(t, err)

	})
}

func TestRepository_UpdateSaga(t *testing.T) {

	t.Parallel()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "slqmock")
	defer sqlxDB.Close()

	regRepo := repository.NewRegistrationRepository(sqlxDB)

	saga := &models.Saga{
		Saga_uuid:   uuid.New(),
		Saga_name:   "test_name",
		Saga_status: 0,
		Saga_type:   1,
	}

	t.Run("Success", func(t *testing.T) {

		mock.ExpectBegin()

		mock.ExpectExec(repository.UpdateSaga).
			WithArgs(
				&saga.Saga_uuid,
				&saga.Saga_status,
				&saga.Saga_type,
				&saga.Saga_name,
			).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err = regRepo.UpdateSaga(context.Background(), saga)
		require.Nil(t, err)

		err = mock.ExpectationsWereMet()
		require.Nil(t, err)

	})

	t.Run("Error", func(t *testing.T) {

		mock.ExpectBegin()

		mock.ExpectExec(repository.UpdateSaga).WillReturnError(errors.New("test error"))

		mock.ExpectRollback()

		err = regRepo.UpdateSaga(context.Background(), saga)
		require.Equal(t, err, repository.ErrorUpdateSaga)

		err = mock.ExpectationsWereMet()
		require.Nil(t, err)

	})

	t.Run("No affected rows", func(t *testing.T) {

		mock.ExpectBegin()

		mock.ExpectExec(repository.UpdateSaga).WillReturnResult(sqlmock.NewResult(0, 0))

		mock.ExpectRollback()

		err = regRepo.UpdateSaga(context.Background(), saga)
		require.Equal(t, err, repository.ErrorUpdateSaga)

		err = mock.ExpectationsWereMet()
		require.Nil(t, err)

	})
}

func TestRepository_CreateSagaConnection(t *testing.T) {

	t.Parallel()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "slqmock")
	defer sqlxDB.Close()

	regRepo := repository.NewRegistrationRepository(sqlxDB)

	saga_connection := &models.SagaConnection{
		Current_saga_uuid:     uuid.New(),
		Next_saga_uuid:        uuid.New(),
		Acc_connection_status: 0,
	}

	t.Run("Success", func(t *testing.T) {

		mock.ExpectBegin()

		mock.ExpectExec(repository.CreateSagaConnection).WithArgs(
			&saga_connection.Current_saga_uuid,
			&saga_connection.Next_saga_uuid,
			&saga_connection.Acc_connection_status,
		).WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		err = regRepo.CreateSagaConnection(context.Background(), saga_connection)
		require.Nil(t, err)

		err = mock.ExpectationsWereMet()
		require.Nil(t, err)

	})

	t.Run("Error", func(t *testing.T) {

		mock.ExpectBegin()

		mock.ExpectExec(repository.CreateSagaConnection).WithArgs(
			&saga_connection.Current_saga_uuid,
			&saga_connection.Next_saga_uuid,
			&saga_connection.Acc_connection_status,
		).WillReturnError(errors.New("test error"))

		mock.ExpectRollback()

		err = regRepo.CreateSagaConnection(context.Background(), saga_connection)
		require.Equal(t, err, repository.ErrorCreateSagaConnection)

		err = mock.ExpectationsWereMet()
		require.Nil(t, err)

	})
}

func TestRepository_GetSagaConnectionsCurrentSaga(t *testing.T) {

	t.Parallel()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "slqmock")
	defer sqlxDB.Close()

	regRepo := repository.NewRegistrationRepository(sqlxDB)

	current_saga_uuid := uuid.New()

	saga_connection1 := &models.SagaConnection{
		Current_saga_uuid:     current_saga_uuid,
		Next_saga_uuid:        uuid.New(),
		Acc_connection_status: 0,
	}

	saga_connection2 := &models.SagaConnection{
		Current_saga_uuid:     current_saga_uuid,
		Next_saga_uuid:        uuid.New(),
		Acc_connection_status: 0,
	}

	list_of_connections := &models.ListOfSagaConnections{}

	list_of_connections.List_of_connetcions = append(list_of_connections.List_of_connetcions, saga_connection1, saga_connection2)

	t.Run("Success", func(t *testing.T) {

		mock.ExpectBegin()

		rows := mock.NewRows([]string{"current_saga_uuid", "next_saga_uuid", "acc_connection_status"}).
			AddRow(
				&saga_connection1.Current_saga_uuid,
				&saga_connection1.Next_saga_uuid,
				&saga_connection1.Acc_connection_status,
			).AddRow(
			&saga_connection2.Current_saga_uuid,
			&saga_connection2.Next_saga_uuid,
			&saga_connection2.Acc_connection_status,
		)

		mock.ExpectQuery(repository.GetSagaConnectionsCurrentSaga).WithArgs(
			&current_saga_uuid,
		).WillReturnRows(rows)

		mock.ExpectCommit()

		result, err := regRepo.GetSagaConnectionsCurrentSaga(context.Background(), current_saga_uuid)
		require.Nil(t, err)
		require.Equal(t, result, list_of_connections)

		err = mock.ExpectationsWereMet()
		require.Nil(t, err)

	})

	t.Run("Error", func(t *testing.T) {

		mock.ExpectBegin()

		mock.ExpectQuery(repository.GetSagaConnectionsCurrentSaga).WithArgs(
			&current_saga_uuid,
		).WillReturnError(errors.New("test error"))

		mock.ExpectRollback()

		result, err := regRepo.GetSagaConnectionsCurrentSaga(context.Background(), current_saga_uuid)
		require.Nil(t, result)
		require.Equal(t, err, repository.ErrorGetSagaCurrentConnections)

		err = mock.ExpectationsWereMet()
		require.Nil(t, err)

	})
}

func TestRepository_GetSagaConnectionsNextSaga(t *testing.T) {

	t.Parallel()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "slqmock")
	defer sqlxDB.Close()

	regRepo := repository.NewRegistrationRepository(sqlxDB)

	next_saga_uuid := uuid.New()

	saga_connection1 := &models.SagaConnection{
		Current_saga_uuid:     uuid.New(),
		Next_saga_uuid:        next_saga_uuid,
		Acc_connection_status: 0,
	}

	saga_connection2 := &models.SagaConnection{
		Current_saga_uuid:     uuid.New(),
		Next_saga_uuid:        next_saga_uuid,
		Acc_connection_status: 0,
	}

	list_of_connections := &models.ListOfSagaConnections{}

	list_of_connections.List_of_connetcions = append(list_of_connections.List_of_connetcions, saga_connection1, saga_connection2)

	t.Run("Success", func(t *testing.T) {

		mock.ExpectBegin()

		rows := mock.NewRows([]string{"current_saga_uuid", "next_saga_uuid", "acc_connection_status"}).
			AddRow(
				&saga_connection1.Current_saga_uuid,
				&saga_connection1.Next_saga_uuid,
				&saga_connection1.Acc_connection_status,
			).AddRow(
			&saga_connection2.Current_saga_uuid,
			&saga_connection2.Next_saga_uuid,
			&saga_connection2.Acc_connection_status,
		)

		mock.ExpectQuery(repository.GetSagaConnectionsNextSaga).WithArgs(
			&next_saga_uuid,
		).WillReturnRows(rows)

		mock.ExpectCommit()

		result, err := regRepo.GetSagaConnectionsNextSaga(context.Background(), next_saga_uuid)
		require.Nil(t, err)
		require.Equal(t, result, list_of_connections)

		err = mock.ExpectationsWereMet()
		require.Nil(t, err)

	})

	t.Run("Error", func(t *testing.T) {

		mock.ExpectBegin()

		mock.ExpectQuery(repository.GetSagaConnectionsNextSaga).WithArgs(
			&next_saga_uuid,
		).WillReturnError(errors.New("test error"))

		mock.ExpectRollback()

		result, err := regRepo.GetSagaConnectionsNextSaga(context.Background(), next_saga_uuid)
		require.Nil(t, result)
		require.Equal(t, err, repository.ErrorGetSagaNextConnections)

		err = mock.ExpectationsWereMet()
		require.Nil(t, err)

	})
}

func TestRepository_DeleteSagaConnection(t *testing.T) {

	t.Parallel()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "slqmock")
	defer sqlxDB.Close()

	regRepo := repository.NewRegistrationRepository(sqlxDB)

	saga_connection := &models.SagaConnection{
		Current_saga_uuid:     uuid.New(),
		Next_saga_uuid:        uuid.New(),
		Acc_connection_status: 0,
	}

	t.Run("Success", func(t *testing.T) {

		mock.ExpectBegin()

		mock.ExpectExec(repository.DeleteSagaConnection).
			WithArgs(
				&saga_connection.Current_saga_uuid,
				&saga_connection.Next_saga_uuid,
			).WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		err = regRepo.DeleteSagaConnection(context.Background(), saga_connection)
		require.Nil(t, err)

		err = mock.ExpectationsWereMet()
		require.Nil(t, err)

	})

	t.Run("Error", func(t *testing.T) {

		mock.ExpectBegin()

		mock.ExpectExec(repository.DeleteSagaConnection).
			WithArgs(
				&saga_connection.Current_saga_uuid,
				&saga_connection.Next_saga_uuid,
			).WillReturnError(errors.New("test error"))

		mock.ExpectRollback()

		err = regRepo.DeleteSagaConnection(context.Background(), saga_connection)
		require.Equal(t, err, repository.ErrorDeleteSagaConnection)

		err = mock.ExpectationsWereMet()
		require.Nil(t, err)

	})
}

func TestRepository_UpdateSagaConnection(t *testing.T) {

	t.Parallel()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "slqmock")
	defer sqlxDB.Close()

	regRepo := repository.NewRegistrationRepository(sqlxDB)

	saga_connection := &models.SagaConnection{
		Current_saga_uuid:     uuid.New(),
		Next_saga_uuid:        uuid.New(),
		Acc_connection_status: 0,
	}

	t.Run("Success", func(t *testing.T) {

		mock.ExpectBegin()

		mock.ExpectExec(repository.UpdateSagaConnection).
			WithArgs(
				&saga_connection.Current_saga_uuid,
				&saga_connection.Next_saga_uuid,
				&saga_connection.Acc_connection_status,
			).WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		err = regRepo.UpdateSagaConnection(context.Background(), saga_connection)
		require.Nil(t, err)

		err = mock.ExpectationsWereMet()
		require.Nil(t, err)

	})

	t.Run("Error", func(t *testing.T) {

		mock.ExpectBegin()

		mock.ExpectExec(repository.UpdateSagaConnection).
			WithArgs(
				&saga_connection.Current_saga_uuid,
				&saga_connection.Next_saga_uuid,
				&saga_connection.Acc_connection_status,
			).WillReturnError(errors.New("test error"))

		mock.ExpectRollback()

		err = regRepo.UpdateSagaConnection(context.Background(), saga_connection)
		require.Equal(t, err, repository.ErrorUpdateSagaConnection)

		err = mock.ExpectationsWereMet()
		require.Nil(t, err)

	})

	t.Run("No affected rows", func(t *testing.T) {

		mock.ExpectBegin()

		mock.ExpectExec(repository.UpdateSagaConnection).
			WithArgs(
				&saga_connection.Current_saga_uuid,
				&saga_connection.Next_saga_uuid,
				&saga_connection.Acc_connection_status,
			).WillReturnResult(sqlmock.NewResult(0, 0))

		mock.ExpectRollback()

		err = regRepo.UpdateSagaConnection(context.Background(), saga_connection)
		require.Equal(t, err, repository.ErrorUpdateSagaConnection)

		err = mock.ExpectationsWereMet()
		require.Nil(t, err)

	})
}

func TestRepository_CreateEvent(t *testing.T) {

	t.Parallel()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "slqmock")
	defer sqlxDB.Close()

	regRepo := repository.NewRegistrationRepository(sqlxDB)

	event := &models.Event{
		Event_uuid:          uuid.New(),
		Saga_uuid:           uuid.New(),
		Event_status:        0,
		Event_name:          "test_name",
		Event_is_roll_back:  false,
		Event_result:        "{}",
		Event_rollback_uuid: uuid.Nil,
	}

	t.Run("Success", func(t *testing.T) {

		mock.ExpectBegin()

		mock.ExpectExec(repository.CreateEvent).
			WithArgs(
				&event.Event_uuid,
				&event.Saga_uuid,
				&event.Event_status,
				&event.Event_name,
				&event.Event_is_roll_back,
				&event.Event_result,
				&event.Event_rollback_uuid,
			).WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		err = regRepo.CreateEvent(context.Background(), event)
		require.Nil(t, err)

		err = mock.ExpectationsWereMet()
		require.Nil(t, err)

	})

	t.Run("Error", func(t *testing.T) {

		mock.ExpectBegin()

		mock.ExpectExec(repository.CreateEvent).
			WithArgs(
				&event.Event_uuid,
				&event.Saga_uuid,
				&event.Event_status,
				&event.Event_name,
				&event.Event_is_roll_back,
				&event.Event_result,
				&event.Event_rollback_uuid,
			).WillReturnError(errors.New("test error"))

		mock.ExpectRollback()

		err = regRepo.CreateEvent(context.Background(), event)
		require.Equal(t, err, repository.ErrorCreateEvent)

		err = mock.ExpectationsWereMet()
		require.Nil(t, err)

	})
}

func TestRepository_GetEvent(t *testing.T) {

	t.Parallel()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "slqmock")
	defer sqlxDB.Close()

	regRepo := repository.NewRegistrationRepository(sqlxDB)

	event := &models.Event{
		Event_uuid:          uuid.New(),
		Saga_uuid:           uuid.New(),
		Event_status:        0,
		Event_name:          "test_name",
		Event_is_roll_back:  false,
		Event_result:        "{}",
		Event_rollback_uuid: uuid.Nil,
	}

	t.Run("Success", func(t *testing.T) {

		mock.ExpectBegin()

		rows := mock.NewRows([]string{"event_uuid", "saga_uuid", "event_status", "event_name", "event_is_roll_back", "event_result", "event_rollback_uuid"}).
			AddRow(
				&event.Event_uuid,
				&event.Saga_uuid,
				&event.Event_status,
				&event.Event_name,
				&event.Event_is_roll_back,
				&event.Event_result,
				&event.Event_rollback_uuid,
			)

		mock.ExpectQuery(repository.GetEvent).
			WithArgs(
				&event.Event_uuid,
			).WillReturnRows(rows)

		mock.ExpectCommit()

		result, err := regRepo.GetEvent(context.Background(), event.Event_uuid)
		require.Nil(t, err)
		require.Equal(t, event, result)

		err = mock.ExpectationsWereMet()
		require.Nil(t, err)

	})

	t.Run("Error", func(t *testing.T) {

		mock.ExpectBegin()

		mock.ExpectQuery(repository.GetEvent).
			WithArgs(
				&event.Event_uuid,
			).WillReturnError(errors.New("test error"))

		mock.ExpectRollback()

		result, err := regRepo.GetEvent(context.Background(), event.Event_uuid)
		require.Nil(t, result)
		require.Equal(t, err, repository.ErrorGetEvent)

		err = mock.ExpectationsWereMet()
		require.Nil(t, err)

	})
}

func TestRepository_DeleteEvent(t *testing.T) {

	t.Parallel()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "slqmock")
	defer sqlxDB.Close()

	regRepo := repository.NewRegistrationRepository(sqlxDB)

	event := &models.Event{
		Event_uuid:          uuid.New(),
		Saga_uuid:           uuid.New(),
		Event_status:        0,
		Event_name:          "test_name",
		Event_is_roll_back:  false,
		Event_result:        "{}",
		Event_rollback_uuid: uuid.Nil,
	}

	t.Run("Success", func(t *testing.T) {

		mock.ExpectBegin()

		mock.ExpectExec(repository.DeleteEvent).
			WithArgs(
				&event.Event_uuid,
			).WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		err = regRepo.DeleteEvent(context.Background(), event.Event_uuid)
		require.Nil(t, err)

		err = mock.ExpectationsWereMet()
		require.Nil(t, err)

	})

	t.Run("Error", func(t *testing.T) {

		mock.ExpectBegin()

		mock.ExpectExec(repository.DeleteEvent).
			WithArgs(
				&event.Event_uuid,
			).WillReturnError(errors.New("test error"))

		mock.ExpectRollback()

		err = regRepo.DeleteEvent(context.Background(), event.Event_uuid)
		require.Equal(t, err, repository.ErrorDeleteEvent)

		err = mock.ExpectationsWereMet()
		require.Nil(t, err)

	})

	t.Run("No affected rows", func(t *testing.T) {

		mock.ExpectBegin()

		mock.ExpectExec(repository.DeleteEvent).
			WithArgs(
				&event.Event_uuid,
			).WillReturnResult(sqlmock.NewResult(0, 0))

		mock.ExpectRollback()

		err = regRepo.DeleteEvent(context.Background(), event.Event_uuid)
		require.Equal(t, err, repository.ErrorDeleteEvent)

		err = mock.ExpectationsWereMet()
		require.Nil(t, err)

	})
}

func TestRepository_UpdateEvent(t *testing.T) {

	t.Parallel()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "slqmock")
	defer sqlxDB.Close()

	regRepo := repository.NewRegistrationRepository(sqlxDB)

	event := &models.Event{
		Event_uuid:          uuid.New(),
		Saga_uuid:           uuid.New(),
		Event_status:        0,
		Event_name:          "test_name",
		Event_is_roll_back:  false,
		Event_result:        "{}",
		Event_rollback_uuid: uuid.Nil,
	}

	t.Run("Success", func(t *testing.T) {

		mock.ExpectBegin()

		mock.ExpectExec(repository.UpdateEvent).
			WithArgs(
				&event.Event_uuid,
				&event.Event_status,
				&event.Event_result,
				&event.Event_rollback_uuid,
			).WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		err = regRepo.UpdateEvent(context.Background(), event)
		require.Nil(t, err)

		err = mock.ExpectationsWereMet()
		require.Nil(t, err)

	})

	t.Run("Error", func(t *testing.T) {

		mock.ExpectBegin()

		mock.ExpectExec(repository.UpdateEvent).
			WithArgs(
				&event.Event_uuid,
				&event.Event_status,
				&event.Event_result,
				&event.Event_rollback_uuid,
			).WillReturnError(errors.New("test error"))

		mock.ExpectRollback()

		err = regRepo.UpdateEvent(context.Background(), event)
		require.Equal(t, err, repository.ErrorUpdateEvent)

		err = mock.ExpectationsWereMet()
		require.Nil(t, err)

	})

	t.Run("No affected rows", func(t *testing.T) {

		mock.ExpectBegin()

		mock.ExpectExec(repository.UpdateEvent).
			WithArgs(
				&event.Event_uuid,
				&event.Event_status,
				&event.Event_result,
				&event.Event_rollback_uuid,
			).WillReturnResult(sqlmock.NewResult(0, 0))

		mock.ExpectRollback()

		err = regRepo.UpdateEvent(context.Background(), event)
		require.Equal(t, err, repository.ErrorUpdateEvent)

		err = mock.ExpectationsWereMet()
		require.Nil(t, err)

	})
}

func TestRepository_GetListOfSagaEvents(t *testing.T) {

	t.Parallel()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "slqmock")
	defer sqlxDB.Close()

	regRepo := repository.NewRegistrationRepository(sqlxDB)

	saga := &models.Saga{
		Saga_uuid:   uuid.New(),
		Saga_name:   "test_name",
		Saga_status: 10,
		Saga_type:   1,
	}

	events := &models.SagaListEvents{}

	event1 := &models.Event{
		Event_uuid:          uuid.New(),
		Saga_uuid:           saga.Saga_uuid,
		Event_rollback_uuid: uuid.Nil,
		Event_result:        "{}",
		Event_name:          "test_name",
		Event_status:        10,
		Event_is_roll_back:  false,
	}

	event2 := &models.Event{
		Event_uuid:          uuid.New(),
		Saga_uuid:           saga.Saga_uuid,
		Event_rollback_uuid: uuid.Nil,
		Event_result:        "{}",
		Event_name:          "test_name",
		Event_status:        10,
		Event_is_roll_back:  false,
	}

	events.EventList = append(events.EventList, event1.Event_uuid, event2.Event_uuid)

	t.Run("Success", func(t *testing.T) {

		mock.ExpectBegin()

		rows := mock.NewRows([]string{"event_uuid"}).
			AddRow(event1.Event_uuid).
			AddRow(event2.Event_uuid)

		mock.ExpectQuery(repository.GetListOfSagaEvents).
			WithArgs(saga.Saga_uuid).
			WillReturnRows(rows)

		mock.ExpectCommit()

		result, err := regRepo.GetListOfSagaEvents(context.Background(), saga.Saga_uuid)
		require.Nil(t, err)
		require.Equal(t, result, events)

	})

	t.Run("Error", func(t *testing.T) {

		mock.ExpectBegin()

		mock.ExpectQuery(repository.GetListOfSagaEvents).
			WithArgs(saga.Saga_uuid).
			WillReturnError(errors.New("test error"))

		mock.ExpectRollback()

		result, err := regRepo.GetListOfSagaEvents(context.Background(), saga.Saga_uuid)
		require.Nil(t, result)
		require.Equal(t, err, repository.ErrorGetListOfSagaEvents)

	})
}
