package test

import (
	"context"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/GCFactory/dbo-system/platform/config"
	"github.com/GCFactory/dbo-system/service/account/internal/account/repository"
	"github.com/GCFactory/dbo-system/service/account/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
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

func TestAccountRepo_CreateAccount(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "slqmock")
	defer sqlxDB.Close()

	accRepo := repository.NewAccountRepository(sqlxDB)

	accRepo = accRepo

	acc_uuid := uuid.New()
	acc_culc_number := "01234567890123456789"
	acc_corr_number := "01234567890123456789"
	acc_bic := "123456789"
	acc_cio := "123456789"
	var acc_money_value uint = 0

	t.Run("Success", func(t *testing.T) {
		account := models.Account{
			Acc_uuid:        acc_uuid,
			Acc_culc_number: acc_culc_number,
			Acc_corr_number: acc_corr_number,
			Acc_bic:         acc_bic,
			Acc_cio:         acc_cio,
			Acc_money_value: acc_money_value,
		}

		row := sqlmock.NewRows([]string{"acc_uuid", "acc_culc_number", "acc_corr_number", "acc_bic", "acc_cio", "acc_money_value"}).AddRow(
			&acc_uuid,
			&acc_culc_number,
			&acc_corr_number,
			&acc_bic,
			&acc_cio,
			&acc_money_value,
		)

		mock.ExpectQuery(repository.CreateAccount).WithArgs(
			&account.Acc_uuid,
			&account.Acc_culc_number,
			&account.Acc_corr_number,
			&account.Acc_bic,
			&account.Acc_cio,
			&account.Acc_money_value,
		).WillReturnRows(row)

		err = accRepo.CreateAccount(context.Background(), &account)
		require.Nil(t, err)
	})
	t.Run("Error", func(t *testing.T) {
		account := models.Account{
			Acc_uuid:        acc_uuid,
			Acc_culc_number: acc_culc_number,
			Acc_corr_number: acc_corr_number,
			Acc_bic:         acc_bic,
			Acc_cio:         acc_cio,
			Acc_money_value: acc_money_value,
		}

		mock.ExpectQuery(repository.CreateAccount).WithArgs(
			&account.Acc_uuid,
			&account.Acc_culc_number,
			&account.Acc_corr_number,
			&account.Acc_bic,
			&account.Acc_cio,
			&account.Acc_money_value,
		).WillReturnError(repository.ErrorCreateAccount)

		err = accRepo.CreateAccount(context.Background(), &account)
		require.Equal(t, err, repository.ErrorCreateAccount)
	})
}

func TestAccountRepo_DeleteAccount(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "slqmock")
	defer sqlxDB.Close()

	accRepo := repository.NewAccountRepository(sqlxDB)

	acc_uuid := uuid.New()

	t.Run("Success", func(t *testing.T) {

		mock.ExpectExec(repository.DeleteAccount).WithArgs(acc_uuid).WillReturnResult(sqlmock.NewResult(1, 1))

		err = accRepo.DeleteAccount(context.Background(), acc_uuid)
		require.Nil(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		mock.ExpectExec(repository.DeleteAccount).WithArgs(acc_uuid).WillReturnError(fmt.Errorf("error"))

		err = accRepo.DeleteAccount(context.Background(), acc_uuid)
		require.Equal(t, err, repository.ErrorDeleteAccount)
	})

	t.Run("Error no data", func(t *testing.T) {
		mock.ExpectExec(repository.DeleteAccount).WithArgs(acc_uuid).WillReturnResult(sqlmock.NewResult(0, 0))

		err = accRepo.DeleteAccount(context.Background(), acc_uuid)
		require.Equal(t, err, repository.ErrorDeleteAccount)
	})
}

func TestAccountRepo_GetAccountData(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "slqmock")
	defer sqlxDB.Close()

	accRepo := repository.NewAccountRepository(sqlxDB)

	acc_uuid := uuid.New()
	acc_culc_number := "01234567890123456789"
	acc_corr_number := "01234567890123456789"
	acc_bic := "123456789"
	acc_cio := "123456789"
	var acc_money_value uint = 0
	var acc_money_amount float64 = 0.0
	var acc_status uint8 = 0

	t.Run("Success", func(t *testing.T) {
		acc := models.Account{
			Acc_uuid:         acc_uuid,
			Acc_culc_number:  acc_culc_number,
			Acc_corr_number:  acc_corr_number,
			Acc_bic:          acc_bic,
			Acc_cio:          acc_cio,
			Acc_money_value:  acc_money_value,
			Acc_money_amount: acc_money_amount,
			Acc_status:       acc_status,
		}

		row := sqlmock.NewRows([]string{"acc_uuid", "acc_culc_number", "acc_corr_number", "acc_bic", "acc_cio", "acc_money_value", "acc_money_amount", "acc_status"}).AddRow(
			&acc_uuid,
			&acc_culc_number,
			&acc_corr_number,
			&acc_bic,
			&acc_cio,
			&acc_money_value,
			&acc_money_amount,
			&acc_status,
		)

		mock.ExpectQuery(repository.GetAccountData).WithArgs(acc_uuid).WillReturnRows(row)

		result, err := accRepo.GetAccountData(context.Background(), acc_uuid)
		require.Nil(t, err)
		require.Equal(t, result, &acc)
	})

	t.Run("No data", func(t *testing.T) {
		mock.ExpectQuery(repository.GetAccountData).WithArgs(acc_uuid).WillReturnError(fmt.Errorf("error"))

		result, err := accRepo.GetAccountData(context.Background(), acc_uuid)
		require.Equal(t, err, repository.ErrorGetAccountData)
		require.Nil(t, result)
	})
}

func TestAccountRepo_UpdateAccountStatus(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "slqmock")
	defer sqlxDB.Close()

	accRepo := repository.NewAccountRepository(sqlxDB)

	acc_uuid := uuid.New()
	var acc_status uint8 = 0

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec(repository.UpdateAccountStatus).WithArgs(acc_uuid, acc_status).WillReturnResult(sqlmock.NewResult(1, 1))

		err = accRepo.UpdateAccountStatus(context.Background(), acc_uuid, acc_status)
		require.Nil(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		mock.ExpectExec(repository.UpdateAccountStatus).WithArgs(acc_uuid, acc_status).WillReturnError(fmt.Errorf("error"))

		err = accRepo.UpdateAccountStatus(context.Background(), acc_uuid, acc_status)
		require.Equal(t, err, repository.ErrorUpdateAccountStatus)
	})

	t.Run("Error no data", func(t *testing.T) {
		mock.ExpectExec(repository.UpdateAccountStatus).WithArgs(acc_uuid, acc_status).WillReturnResult(sqlmock.NewResult(0, 0))

		err = accRepo.UpdateAccountStatus(context.Background(), acc_uuid, acc_status)
		require.Equal(t, err, repository.ErrorUpdateAccountStatus)
	})
}

func TestAccountRepo_GetAccountStatus(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "slqmock")
	defer sqlxDB.Close()

	accRepo := repository.NewAccountRepository(sqlxDB)

	acc_uuid := uuid.New()
	var acc_status uint8 = 0

	t.Run("Success", func(t *testing.T) {
		row := mock.NewRows([]string{"acc_status"}).AddRow(
			&acc_status,
		)

		mock.ExpectQuery(repository.GetAccountStatus).WithArgs(acc_uuid).WillReturnRows(row)

		result, err := accRepo.GetAccountStatus(context.Background(), acc_uuid)
		require.Nil(t, err)
		require.Equal(t, result, acc_status)
	})
	t.Run("Error", func(t *testing.T) {
		mock.ExpectQuery(repository.GetAccountStatus).WithArgs(acc_uuid).WillReturnError(fmt.Errorf("error"))

		result, err := accRepo.GetAccountStatus(context.Background(), acc_uuid)
		require.Equal(t, err, repository.ErrorGetAccountStatus)
		require.Equal(t, result, uint8(0))
	})
}

func TestAccountRepo_UpdateAccountAmount(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "slqmock")
	defer sqlxDB.Close()

	accRepo := repository.NewAccountRepository(sqlxDB)

	acc_uuid := uuid.New()
	var acc_amount float64 = 1.0

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec(repository.UpdateAccountAmount).WithArgs(acc_uuid, acc_amount).WillReturnResult(sqlmock.NewResult(1, 1))

		err = accRepo.UpdateAccountAmount(context.Background(), acc_uuid, acc_amount)
		require.Nil(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		mock.ExpectExec(repository.UpdateAccountAmount).WithArgs(acc_uuid, acc_amount).WillReturnError(fmt.Errorf("error"))

		err = accRepo.UpdateAccountAmount(context.Background(), acc_uuid, acc_amount)
		require.Equal(t, err, repository.ErrorUpdateAccountAmount)
	})

	t.Run("Error no data", func(t *testing.T) {
		mock.ExpectExec(repository.UpdateAccountAmount).WithArgs(acc_uuid, acc_amount).WillReturnResult(sqlmock.NewResult(0, 0))

		err = accRepo.UpdateAccountAmount(context.Background(), acc_uuid, acc_amount)
		require.Equal(t, err, repository.ErrorUpdateAccountAmount)
	})
}

func TestAccountRepo_GetAccountAmount(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "slqmock")
	defer sqlxDB.Close()

	accRepo := repository.NewAccountRepository(sqlxDB)

	acc_uuid := uuid.New()
	var acc_amount float64 = 1.0

	t.Run("Success", func(t *testing.T) {
		row := mock.NewRows([]string{"acc_money_amount"}).AddRow(
			&acc_amount,
		)

		mock.ExpectQuery(repository.GetAccountAmount).WithArgs(acc_uuid).WillReturnRows(row)

		result, err := accRepo.GetAccountAmount(context.Background(), acc_uuid)
		require.Nil(t, err)
		require.Equal(t, result, acc_amount)
	})

	t.Run("Error", func(t *testing.T) {
		mock.ExpectQuery(repository.GetAccountAmount).WithArgs(acc_uuid).WillReturnError(fmt.Errorf("error"))

		result, err := accRepo.GetAccountAmount(context.Background(), acc_uuid)
		require.Equal(t, err, repository.ErrorGetAccountAmount)
		require.Equal(t, result, float64(0))
	})
}

func TestAccountRepo_AddReserveReason(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "slqmock")
	defer sqlxDB.Close()

	accRepo := repository.NewAccountRepository(sqlxDB)

	acc_uuid := uuid.New()
	acc_reason := "reason"

	t.Run("Success", func(t *testing.T) {
		row := mock.NewRows([]string{"acc_uuid", "reserve_reason"}).AddRow(
			&acc_uuid,
			&acc_reason,
		)

		reason := models.ReserverReason{
			Acc_uuid: acc_uuid,
			Reason:   acc_reason,
		}

		mock.ExpectQuery(repository.AddReserveReason).WithArgs(acc_uuid, acc_reason).WillReturnRows(row)

		err = accRepo.AddReserveReason(context.Background(), &reason)
		require.Nil(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		reason := models.ReserverReason{
			Acc_uuid: acc_uuid,
			Reason:   acc_reason,
		}

		mock.ExpectQuery(repository.AddReserveReason).WithArgs(acc_uuid, acc_reason).WillReturnError(fmt.Errorf("error"))

		err = accRepo.AddReserveReason(context.Background(), &reason)
		require.Equal(t, err, repository.ErrorAddReserveReason)
	})
}

func TestAccountRepo_DeleteReserveReason(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "slqmock")
	defer sqlxDB.Close()

	accRepo := repository.NewAccountRepository(sqlxDB)

	acc_uuid := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec(repository.DeleteReserveReason).WithArgs(acc_uuid).WillReturnResult(sqlmock.NewResult(1, 1))

		err = accRepo.DeleteReserveReason(context.Background(), acc_uuid)
		require.Nil(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		mock.ExpectExec(repository.DeleteReserveReason).WithArgs(acc_uuid).WillReturnError(fmt.Errorf("error"))

		err = accRepo.DeleteReserveReason(context.Background(), acc_uuid)
		require.Equal(t, err, repository.ErrorDeleteReserveReason)
	})

	t.Run("Error no data", func(t *testing.T) {
		mock.ExpectExec(repository.DeleteReserveReason).WithArgs(acc_uuid).WillReturnResult(sqlmock.NewResult(0, 0))

		err = accRepo.DeleteReserveReason(context.Background(), acc_uuid)
		require.Equal(t, err, repository.ErrorDeleteReserveReason)
	})
}

func TestAccountRepo_GetReserveReason(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(db, "slqmock")
	defer sqlxDB.Close()

	accRepo := repository.NewAccountRepository(sqlxDB)

	acc_uuid := uuid.New()
	acc_reason := "reason"
	reason := models.ReserverReason{
		Acc_uuid: acc_uuid,
		Reason:   acc_reason,
	}

	reason = reason

	t.Run("Success", func(t *testing.T) {
		row := mock.NewRows([]string{"acc_uuid", "reserve_reason"}).AddRow(
			&acc_uuid,
			&acc_reason,
		)

		mock.ExpectQuery(repository.GetReserveReason).WithArgs(acc_uuid).WillReturnRows(row)

		result, err := accRepo.GetReserveReason(context.Background(), acc_uuid)
		require.Nil(t, err)
		require.Equal(t, result, &reason)
	})

	t.Run("Error", func(t *testing.T) {
		mock.ExpectQuery(repository.GetReserveReason).WithArgs(acc_uuid).WillReturnError(fmt.Errorf("error"))

		result, err := accRepo.GetReserveReason(context.Background(), acc_uuid)
		require.Equal(t, err, repository.ErrorGetReserveReason)
		require.Nil(t, result)
	})
}
