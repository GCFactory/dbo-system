package test

import (
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/GCFactory/dbo-system/platform/config"
	"github.com/GCFactory/dbo-system/service/registration/internal/models"
	"github.com/GCFactory/dbo-system/service/registration/internal/registration/repository"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"strconv"
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

func TestRegistrationRepo_CreateSaga(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "slqmock")
	defer sqlxDB.Close()

	regRepo := repository.NewRegistrationRepository(sqlxDB)

	t.Run("Success", func(t *testing.T) {
		saga_uuid := uuid.New()
		var saga_status uint = 255
		var saga_type uint = 0
		saga_name := "test_saga"

		events := models.SagaListEvents{}
		events.EventList = append(events.EventList, "Event1")
		events.EventList = append(events.EventList, "Event2")

		list_of_events, err := json.Marshal(events)
		require.NoError(t, err)
		require.NotEmpty(t, list_of_events)

		row := sqlmock.NewRows([]string{"saga_uuid", "saga_status", "saga_type", "saga_name", "list_of_events"}).AddRow(
			saga_uuid.String(),
			strconv.FormatInt(int64(saga_status), 10),
			strconv.FormatInt(int64(saga_type), 10),
			&saga_name,
			&list_of_events,
		)

		saga := models.Saga{
			Saga_uuid:           saga_uuid,
			Saga_status:         saga_status,
			Saga_type:           saga_type,
			Saga_name:           saga_name,
			Saga_list_of_events: string(list_of_events),
		}

		mock.ExpectQuery(repository.CreateSaga).WithArgs(
			&saga.Saga_uuid,
			&saga.Saga_status,
			&saga.Saga_type,
			&saga.Saga_name,
			&saga.Saga_list_of_events,
		).WillReturnRows(row)

		err = regRepo.CreateSaga(context.Background(), saga)

		require.NoError(t, err)
		require.Nil(t, err)
	})
	t.Run("WrongData", func(t *testing.T) {
		saga_uuid := uuid.New()
		var saga_status uint = 255
		var saga_type uint = 0

		events := models.SagaListEvents{}
		events.EventList = append(events.EventList, "Event1")
		events.EventList = append(events.EventList, "Event2")

		list_of_events, err := json.Marshal(events)
		require.NoError(t, err)
		require.NotEmpty(t, list_of_events)

		saga := models.Saga{
			Saga_uuid:           saga_uuid,
			Saga_status:         saga_status,
			Saga_type:           saga_type,
			Saga_name:           "",
			Saga_list_of_events: string(list_of_events),
		}

		mock.ExpectQuery(repository.CreateSaga).WithArgs(
			&saga.Saga_uuid,
			&saga.Saga_status,
			&saga.Saga_type,
			&saga.Saga_name,
			&saga.Saga_list_of_events,
		).WillReturnError(err)

		err = regRepo.CreateSaga(context.Background(), saga)

		require.NotNil(t, err)
		require.Error(t, err)
		require.Equal(t, err, repository.ErrorCreateSaga)
	})
}

func TestRegistrationRepo_DeleteSaga(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "slqmock")
	defer sqlxDB.Close()

	regRepo := repository.NewRegistrationRepository(sqlxDB)

	t.Run("Success", func(t *testing.T) {
		saga_uuid := uuid.New()

		mock.ExpectExec(repository.DeleteSaga).WithArgs(saga_uuid).WillReturnResult(sqlmock.NewResult(1, 1))

		err = regRepo.DeleteSaga(context.Background(), saga_uuid)

		require.Nil(t, err)
	})
	t.Run("WrongData", func(t *testing.T) {
		saga_uuid := uuid.New()

		mock.ExpectExec(repository.DeleteSaga).WithArgs(saga_uuid).WillReturnError(err)

		err = regRepo.DeleteSaga(context.Background(), saga_uuid)

		require.NotNil(t, err)
		require.Equal(t, err, repository.ErrorDeleteSaga)
	})
}

func TestRegistrationRepo_GetSagaById(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "slqmock")
	defer sqlxDB.Close()

	regRepo := repository.NewRegistrationRepository(sqlxDB)

	t.Run("Success", func(t *testing.T) {
		saga_uuid := uuid.New()
		var saga_status uint = 255
		var saga_type uint = 0
		saga_name := "test_saga"

		events := models.SagaListEvents{}
		events.EventList = append(events.EventList, "Event1")
		events.EventList = append(events.EventList, "Event2")

		list_of_events, err := json.Marshal(events)
		require.NoError(t, err)
		require.NotEmpty(t, list_of_events)

		row := sqlmock.NewRows([]string{"saga_uuid", "saga_status", "saga_type", "saga_name", "list_of_events"}).AddRow(
			saga_uuid.String(),
			strconv.FormatInt(int64(saga_status), 10),
			strconv.FormatInt(int64(saga_type), 10),
			&saga_name,
			&list_of_events,
		)

		mock.ExpectQuery(repository.GetSaga).WithArgs(
			&saga_uuid,
		).WillReturnRows(row)

		saga := models.Saga{
			Saga_uuid:           saga_uuid,
			Saga_status:         saga_status,
			Saga_type:           saga_type,
			Saga_name:           saga_name,
			Saga_list_of_events: string(list_of_events),
		}

		result, err := regRepo.GetSagaById(context.Background(), saga_uuid)

		require.NoError(t, err)
		require.Nil(t, err)
		require.Equal(t, *result, saga)
	})
	t.Run("WrongData", func(t *testing.T) {
		saga_uuid := uuid.New()

		mock.ExpectQuery(repository.GetSaga).WithArgs(
			&saga_uuid,
		).WillReturnError(err)

		result, err := regRepo.GetSagaById(context.Background(), saga_uuid)

		require.NotNil(t, err)
		require.Equal(t, err, repository.ErrorGetSaga)
		require.Nil(t, result)
	})
}

func TestRegistrationRepo_UpdateSagaStatus(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "slqmock")
	defer sqlxDB.Close()

	regRepo := repository.NewRegistrationRepository(sqlxDB)

	t.Run("Success", func(t *testing.T) {
		saga_uuid := uuid.New()
		var saga_status uint = 255

		mock.ExpectExec(repository.UpdateSagaStatus).WithArgs(
			&saga_uuid,
			&saga_status,
		).WillReturnResult(sqlmock.NewResult(1, 1))

		err := regRepo.UpdateSagaStatus(context.Background(), saga_uuid, saga_status)

		require.NoError(t, err)
		require.Nil(t, err)
	})
	t.Run("WrongData", func(t *testing.T) {
		saga_uuid := uuid.New()
		var saga_status uint = 255

		mock.ExpectExec(repository.UpdateSagaStatus).WithArgs(
			&saga_uuid,
			&saga_status,
		).WillReturnError(err)

		err := regRepo.UpdateSagaStatus(context.Background(), saga_uuid, saga_status)

		require.NotNil(t, err)
		require.Equal(t, err, repository.ErrorUpdateSagaStatus)
	})
}

func TestRegistrationRepo_UpdateSagaEvents(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "slqmock")
	defer sqlxDB.Close()

	regRepo := repository.NewRegistrationRepository(sqlxDB)

	t.Run("Success", func(t *testing.T) {
		saga_uuid := uuid.New()

		event_list := models.SagaListEvents{}
		event_list.EventList = append(event_list.EventList, "Event1")
		event_list.EventList = append(event_list.EventList, "Event2")

		event_list_str, err := json.Marshal(event_list)
		require.Nil(t, err)
		require.NotEmpty(t, event_list_str)

		mock.ExpectExec(repository.UpdateSagaEvents).WithArgs(
			&saga_uuid,
			string(event_list_str),
		).WillReturnResult(sqlmock.NewResult(1, 1))

		err = regRepo.UpdateSagaEvents(context.Background(), saga_uuid, string(event_list_str))

		require.NoError(t, err)
		require.Nil(t, err)
	})
	t.Run("WrongData", func(t *testing.T) {
		saga_uuid := uuid.New()

		event_list := models.SagaListEvents{}
		event_list.EventList = append(event_list.EventList, "Event1")
		event_list.EventList = append(event_list.EventList, "Event2")

		event_list_str, err := json.Marshal(event_list)
		require.Nil(t, err)
		require.NotEmpty(t, event_list_str)

		mock.ExpectExec(repository.UpdateSagaEvents).WithArgs(
			&saga_uuid,
			string(event_list_str),
		).WillReturnError(err)

		err = regRepo.UpdateSagaEvents(context.Background(), saga_uuid, string(event_list_str))

		require.NotNil(t, err)
		require.Equal(t, err, repository.ErrorUpdateSagaEvents)
	})
}
