package tests

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/GCFactory/dbo-system/platform/config"
	"github.com/GCFactory/dbo-system/service/users/internal/models"
	"github.com/GCFactory/dbo-system/service/users/internal/users/repository"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
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

func TestRepository_AddUser(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "slqmock")
	defer sqlxDB.Close()

	userRepo := repository.NewUserRepository(sqlxDB)

	mock.ExpectationsWereMet()

	passport := &models.Passport{
		Passport_uuid:       uuid.New(),
		Registration_adress: "...",
		Authority_date:      time.DateTime,
		Authority:           "5",
		Pick_up_point:       "000-111",
		Birth_location:      "4",
		Birth_date:          time.DateTime,
		Patronimic:          "3",
		Name:                "1",
		Passport_number:     "123456",
		Passport_series:     "1234",
		Surname:             "2",
	}

	user := &models.User{
		User_uuid:     uuid.New(),
		Passport_uuid: passport.Passport_uuid,
		User_inn:      "01234567890123456789",
		User_accounts: "{\"accounts\": [" +
			"\"" + uuid.New().String() + "\"," +
			"\"" + uuid.New().String() + "\"" +
			"]" +
			"}",
	}

	user_full := &models.User_full_data{
		User:     user,
		Passport: passport,
	}

	t.Run("Success", func(t *testing.T) {

		mock.ExpectBegin()

		mock.ExpectExec(repository.AddPassport).WithArgs(
			&passport.Passport_uuid,
			&passport.Passport_series,
			&passport.Passport_number,
			&passport.Name,
			&passport.Surname,
			&passport.Patronimic,
			&passport.Birth_date,
			&passport.Birth_location,
			&passport.Pick_up_point,
			&passport.Authority,
			&passport.Authority_date,
			&passport.Registration_adress,
		).WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec(repository.AddUser).WithArgs(
			&user.User_uuid,
			&user.Passport_uuid,
			&user.User_inn,
			&user.User_accounts,
		).WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		err = userRepo.AddUser(context.Background(), user_full)
		require.Nil(t, err)

		err = mock.ExpectationsWereMet()
		require.Nil(t, err)

	})

	t.Run("Error in passport data", func(t *testing.T) {

		tmp := passport.Passport_number
		passport.Passport_number = ""

		mock.ExpectBegin()

		mock.ExpectExec(repository.AddPassport).WithArgs(
			&passport.Passport_uuid,
			&passport.Passport_series,
			&passport.Passport_number,
			&passport.Name,
			&passport.Surname,
			&passport.Patronimic,
			&passport.Birth_date,
			&passport.Birth_location,
			&passport.Pick_up_point,
			&passport.Authority,
			&passport.Authority_date,
			&passport.Registration_adress,
		).WillReturnError(err)

		mock.ExpectRollback()

		err = userRepo.AddUser(context.Background(), user_full)
		require.Equal(t, err, repository.ErrorAddPassport)

		err = mock.ExpectationsWereMet()
		require.Nil(t, err)

		passport.Passport_series = tmp

	})

	t.Run("Error in user data", func(t *testing.T) {

		tmp := user.Passport_uuid
		user.Passport_uuid = uuid.New()

		mock.ExpectBegin()

		mock.ExpectExec(repository.AddPassport).WithArgs(
			&passport.Passport_uuid,
			&passport.Passport_series,
			&passport.Passport_number,
			&passport.Name,
			&passport.Surname,
			&passport.Patronimic,
			&passport.Birth_date,
			&passport.Birth_location,
			&passport.Pick_up_point,
			&passport.Authority,
			&passport.Authority_date,
			&passport.Registration_adress,
		).WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec(repository.AddUser).WithArgs(
			&user.User_uuid,
			&user.Passport_uuid,
			&user.User_inn,
			&user.User_accounts,
		).WillReturnError(err)

		mock.ExpectRollback()

		err = userRepo.AddUser(context.Background(), user_full)
		require.Equal(t, err, repository.ErrorAddUser)

		err = mock.ExpectationsWereMet()
		require.Nil(t, err)

		user.Passport_uuid = tmp

	})

	t.Run("Error rollback transaction", func(t *testing.T) {

		tmp := user.Passport_uuid
		user.Passport_uuid = uuid.New()

		mock.ExpectBegin()

		mock.ExpectExec(repository.AddPassport).WithArgs(
			&passport.Passport_uuid,
			&passport.Passport_series,
			&passport.Passport_number,
			&passport.Name,
			&passport.Surname,
			&passport.Patronimic,
			&passport.Birth_date,
			&passport.Birth_location,
			&passport.Pick_up_point,
			&passport.Authority,
			&passport.Authority_date,
			&passport.Registration_adress,
		).WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec(repository.AddUser).WithArgs(
			&user.User_uuid,
			&user.Passport_uuid,
			&user.User_inn,
			&user.User_accounts,
		).WillReturnError(err)

		mock.ExpectRollback().WillReturnError(err)

		err = userRepo.AddUser(context.Background(), user_full)
		require.Error(t, err)

		err = mock.ExpectationsWereMet()
		require.Nil(t, err)

		user.Passport_uuid = tmp
	})

	t.Run("Error begin transaction", func(t *testing.T) {

		mock.ExpectClose()
		mock.ExpectBegin().WillReturnError(err)

		err = userRepo.AddUser(context.Background(), user_full)
		require.Error(t, err)

		err = mock.ExpectationsWereMet()
		require.Error(t, err)

	})

	t.Run("Error commit transaction", func(t *testing.T) {

		mock.ExpectBegin()

		mock.ExpectExec(repository.AddPassport).WithArgs(
			&passport.Passport_uuid,
			&passport.Passport_series,
			&passport.Passport_number,
			&passport.Name,
			&passport.Surname,
			&passport.Patronimic,
			&passport.Birth_date,
			&passport.Birth_location,
			&passport.Pick_up_point,
			&passport.Authority,
			&passport.Authority_date,
			&passport.Registration_adress,
		).WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec(repository.AddUser).WithArgs(
			&user.User_uuid,
			&user.Passport_uuid,
			&user.User_inn,
			&user.User_accounts,
		).WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectClose()
		mock.ExpectCommit().WillReturnError(err)

		err = userRepo.AddUser(context.Background(), user_full)
		require.Error(t, err)

		err = mock.ExpectationsWereMet()
		require.Error(t, err)
	})
}
