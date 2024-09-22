package test

import (
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

		rows := mock.NewRows([]string{"event_uuid"}).
			AddRow(event1.Event_uuid).
			AddRow(event2.Event_uuid)

		mock.ExpectQuery(repository.GetListOfSagaEvents).
			WithArgs(saga.Saga_uuid).
			WillReturnRows(rows)

		result, err := regRepo.GetListOfSagaEvents(context.Background(), saga.Saga_uuid)
		require.Nil(t, err)
		require.Equal(t, result, events)

	})
}
