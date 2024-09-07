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

		sqlxDB.Close()
		mock.ExpectBegin().WillReturnError(err)

		err = userRepo.AddUser(context.Background(), user_full)
		require.Error(t, err)

		sqlxDB.Conn(context.Background())

	})
}

func TestRepository_GetUSerData(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "slqmock")
	defer sqlxDB.Close()

	userRepo := repository.NewUserRepository(sqlxDB)

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

		passport_row := mock.NewRows([]string{
			"passport_uuid",
			"passport_series",
			"passport_number",
			"name",
			"surname",
			"patronimic",
			"birth_date",
			"birth_location",
			"pick_up_point",
			"authority",
			"authority_date",
			"registration_adress"},
		).AddRow(
			passport.Passport_uuid,
			passport.Passport_series,
			passport.Passport_number,
			passport.Name,
			passport.Surname,
			passport.Patronimic,
			passport.Birth_date,
			passport.Birth_location,
			passport.Pick_up_point,
			passport.Authority,
			passport.Authority_date,
			passport.Registration_adress,
		)

		user_row := mock.NewRows([]string{
			"user_uuid",
			"passport_uuid",
			"user_inn",
			"user_accounts"},
		).AddRow(
			user.User_uuid,
			user.Passport_uuid,
			user.User_inn,
			user.User_accounts,
		)

		mock.ExpectBegin()

		mock.ExpectQuery(repository.GetUserData).WithArgs(
			user.User_uuid,
		).WillReturnRows(user_row)
		mock.ExpectQuery(repository.GetPassportData).WithArgs(
			passport.Passport_uuid,
		).WillReturnRows(passport_row)

		mock.ExpectCommit()

		result, err := userRepo.GetUserData(context.Background(), user.User_uuid)
		require.Nil(t, err)
		require.Equal(t, user_full, result)

		err = mock.ExpectationsWereMet()
		require.Nil(t, err)

	})

	t.Run("Error user uuid", func(t *testing.T) {

		mock.ExpectBegin()

		mock.ExpectQuery(repository.GetUserData).WithArgs(user.User_uuid).WillReturnError(err)

		mock.ExpectRollback()

		result, err := userRepo.GetUserData(context.Background(), user.User_uuid)
		require.Equal(t, err, repository.ErrorGetUser)
		require.Nil(t, result)

		err = mock.ExpectationsWereMet()
		require.Nil(t, err)

	})

	t.Run("Error getting passport data", func(t *testing.T) {

		user_row := mock.NewRows([]string{
			"user_uuid",
			"passport_uuid",
			"user_inn",
			"user_accounts"},
		).AddRow(
			user.User_uuid,
			user.Passport_uuid,
			user.User_inn,
			user.User_accounts,
		)

		mock.ExpectBegin()

		mock.ExpectQuery(repository.GetUserData).WithArgs(
			user.User_uuid,
		).WillReturnRows(user_row)
		mock.ExpectQuery(repository.GetPassportData).WithArgs(
			passport.Passport_uuid,
		).WillReturnError(err)

		mock.ExpectRollback()

		result, err := userRepo.GetUserData(context.Background(), user.User_uuid)
		require.Nil(t, result)
		require.Equal(t, err, repository.ErrorGetPassport)

		err = mock.ExpectationsWereMet()
		require.Nil(t, err)
	})

	t.Run("Error rollback", func(t *testing.T) {

		mock.ExpectBegin()

		mock.ExpectQuery(repository.GetUserData).WithArgs(
			user.User_uuid,
		).WillReturnError(err)

		mock.ExpectRollback().WillReturnError(err)

		result, err := userRepo.GetUserData(context.Background(), user.User_uuid)
		require.Nil(t, result)
		require.Error(t, err)

		err = mock.ExpectationsWereMet()
		require.Nil(t, err)
	})

	t.Run("Error begin transaction", func(t *testing.T) {

		sqlxDB.Close()
		mock.ExpectBegin().WillReturnError(err)

		result, err := userRepo.GetUserData(context.Background(), user.User_uuid)
		require.Nil(t, result)
		require.Error(t, err)

		sqlxDB.Conn(context.Background())

	})

}

func TestRepository_UpdatePassport(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "slqmock")
	defer sqlxDB.Close()

	userRepo := repository.NewUserRepository(sqlxDB)

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

	t.Run("Success", func(t *testing.T) {

		mock.ExpectBegin()

		mock.ExpectExec(repository.UpdatePassport).WithArgs(
			passport.Passport_uuid,
			passport.Passport_series,
			passport.Passport_number,
			passport.Name,
			passport.Surname,
			passport.Patronimic,
			passport.Birth_date,
			passport.Birth_location,
			passport.Pick_up_point,
			passport.Authority,
			passport.Authority_date,
			passport.Registration_adress,
		).WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		err = userRepo.UpdatePassport(context.Background(), passport)
		require.Nil(t, err)

		err = mock.ExpectationsWereMet()
		require.Nil(t, err)

	})

	t.Run("Error update passport", func(t *testing.T) {

		mock.ExpectBegin()

		mock.ExpectExec(repository.UpdatePassport).WithArgs(
			passport.Passport_uuid,
			passport.Passport_series,
			passport.Passport_number,
			passport.Name,
			passport.Surname,
			passport.Patronimic,
			passport.Birth_date,
			passport.Birth_location,
			passport.Pick_up_point,
			passport.Authority,
			passport.Authority_date,
			passport.Registration_adress,
		).WillReturnError(err)

		mock.ExpectRollback()

		err = userRepo.UpdatePassport(context.Background(), passport)
		require.Equal(t, err, repository.ErrorUpdatePassport)

		err = mock.ExpectationsWereMet()
		require.Nil(t, err)

	})

	t.Run("Error update passport, no such row", func(t *testing.T) {

		mock.ExpectBegin()

		mock.ExpectExec(repository.UpdatePassport).WithArgs(
			passport.Passport_uuid,
			passport.Passport_series,
			passport.Passport_number,
			passport.Name,
			passport.Surname,
			passport.Patronimic,
			passport.Birth_date,
			passport.Birth_location,
			passport.Pick_up_point,
			passport.Authority,
			passport.Authority_date,
			passport.Registration_adress,
		).WillReturnResult(sqlmock.NewResult(0, 0))

		mock.ExpectRollback()

		err = userRepo.UpdatePassport(context.Background(), passport)
		require.Equal(t, err, repository.ErrorUpdatePassport)

		err = mock.ExpectationsWereMet()
		require.Nil(t, err)

	})

	t.Run("Error rollback", func(t *testing.T) {

		mock.ExpectBegin()

		mock.ExpectExec(repository.UpdatePassport).WithArgs(
			passport.Passport_uuid,
			passport.Passport_series,
			passport.Passport_number,
			passport.Name,
			passport.Surname,
			passport.Patronimic,
			passport.Birth_date,
			passport.Birth_location,
			passport.Pick_up_point,
			passport.Authority,
			passport.Authority_date,
			passport.Registration_adress,
		).WillReturnError(err)

		mock.ExpectRollback().WillReturnError(err)

		err = userRepo.UpdatePassport(context.Background(), passport)
		require.Error(t, err)

		err = mock.ExpectationsWereMet()
		require.Nil(t, err)

	})

	t.Run("Error begin transaction", func(t *testing.T) {

		sqlxDB.Close()
		mock.ExpectBegin().WillReturnError(err)

		mock.ExpectRollback()

		err = userRepo.UpdatePassport(context.Background(), passport)
		require.Error(t, err)
		sqlxDB.Conn(context.Background())

	})

}

func TestRepository_UpdateUserAccount(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "slqmock")
	defer sqlxDB.Close()

	userRepo := repository.NewUserRepository(sqlxDB)

	user_uuid := uuid.New()

	new_accounts := "{\"accounts\": [" + uuid.New().String() + "]}"

	t.Run("Success", func(t *testing.T) {

		mock.ExpectBegin()

		mock.ExpectExec(repository.UpdateUsersAccounts).WithArgs(
			user_uuid,
			new_accounts,
		).WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		err = userRepo.UpdateUserAccount(context.Background(), user_uuid, new_accounts)
		require.Nil(t, err)

		err = mock.ExpectationsWereMet()
		require.Nil(t, err)

	})

	t.Run("Error update account", func(t *testing.T) {

		mock.ExpectBegin()

		mock.ExpectExec(repository.UpdateUsersAccounts).WithArgs(
			user_uuid,
			new_accounts,
		).WillReturnError(err)

		mock.ExpectRollback()

		err = userRepo.UpdateUserAccount(context.Background(), user_uuid, new_accounts)
		require.Equal(t, err, repository.ErrorUpdateAccounts)

		err = mock.ExpectationsWereMet()
		require.Nil(t, err)

	})

	t.Run("Error update account, no such row", func(t *testing.T) {

		mock.ExpectBegin()

		mock.ExpectExec(repository.UpdateUsersAccounts).WithArgs(
			user_uuid,
			new_accounts,
		).WillReturnResult(sqlmock.NewResult(0, 0))

		mock.ExpectRollback()

		err = userRepo.UpdateUserAccount(context.Background(), user_uuid, new_accounts)
		require.Equal(t, err, repository.ErrorUpdateAccounts)

		err = mock.ExpectationsWereMet()
		require.Nil(t, err)

	})

	t.Run("Error rollbcak transaction", func(t *testing.T) {

		mock.ExpectBegin()

		mock.ExpectExec(repository.UpdateUsersAccounts).WithArgs(
			user_uuid,
			new_accounts,
		).WillReturnResult(sqlmock.NewResult(0, 0))

		mock.ExpectRollback().WillReturnError(err)

		err = userRepo.UpdateUserAccount(context.Background(), user_uuid, new_accounts)
		require.Error(t, err)

		err = mock.ExpectationsWereMet()
		require.Nil(t, err)

	})

	t.Run("Error begin transaction", func(t *testing.T) {

		sqlxDB.Close()
		mock.ExpectBegin().WillReturnError(err)

		err = userRepo.UpdateUserAccount(context.Background(), user_uuid, new_accounts)
		require.Error(t, err)

		sqlxDB.Conn(context.Background())

	})

}

func TestRepository_GetUsersAccounts(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "slqmock")
	defer sqlxDB.Close()

	userRepo := repository.NewUserRepository(sqlxDB)

	user_uuid := uuid.New()

	accounts := "{\"accounts\": [" + uuid.New().String() + "]}"

	t.Run("Success", func(t *testing.T) {

		rows := mock.NewRows([]string{"user_accounts"}).AddRow(accounts)

		mock.ExpectBegin()

		mock.ExpectQuery(repository.GetUsersAccounts).WithArgs(
			user_uuid,
		).WillReturnRows(rows)

		mock.ExpectCommit()

		result, err := userRepo.GetUsersAccounts(context.Background(), user_uuid)
		require.Nil(t, err)
		require.Equal(t, result, accounts)

		err = mock.ExpectationsWereMet()
		require.Nil(t, err)

	})

	t.Run("Error get accounts", func(t *testing.T) {

		mock.ExpectBegin()

		mock.ExpectQuery(repository.GetUsersAccounts).WithArgs(
			user_uuid,
		).WillReturnError(err)

		mock.ExpectRollback()

		result, err := userRepo.GetUsersAccounts(context.Background(), user_uuid)
		require.Equal(t, result, "")
		require.Equal(t, err, repository.ErrorGetUsersAccounts)

		err = mock.ExpectationsWereMet()
		require.Nil(t, err)

	})

	t.Run("Error rollback transaction", func(t *testing.T) {

		mock.ExpectBegin()

		mock.ExpectQuery(repository.GetUsersAccounts).WithArgs(
			user_uuid,
		).WillReturnError(err)

		mock.ExpectRollback().WillReturnError(err)

		result, err := userRepo.GetUsersAccounts(context.Background(), user_uuid)
		require.Equal(t, result, "")
		require.Error(t, err)

		err = mock.ExpectationsWereMet()
		require.Nil(t, err)

	})

	t.Run("Error begin transaction", func(t *testing.T) {

		sqlxDB.Close()

		mock.ExpectBegin().WillReturnError(err)

		result, err := userRepo.GetUsersAccounts(context.Background(), user_uuid)
		require.Equal(t, result, "")
		require.Error(t, err)

		sqlxDB.Conn(context.Background())
	})

}
