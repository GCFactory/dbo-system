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

		event1 := models.SagaEvent{Name: "Event1"}
		event2 := models.SagaEvent{Name: "Event2"}
		events := models.SagaListEvents{}
		events.EventList = append(events.EventList, &event1)
		events.EventList = append(events.EventList, &event2)

		list_of_events, err := json.Marshal(events)
		require.NoError(t, err)
		require.NotEmpty(t, list_of_events)

		row := sqlmock.NewRows([]string{"saga_uuid", "sage_status", "saga_type", "saga_name", "list_of_events"}).AddRow(
			saga_uuid.String(),
			strconv.FormatInt(int64(saga_status), 10),
			strconv.FormatInt(int64(saga_type), 10),
			saga_name,
			string(list_of_events),
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
		require.Equal(t, err, nil)
	})
}
