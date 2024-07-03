package test

import (
	"context"
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

//func TestAccountRepo_DeleteAccount(t *testing.T) {
//	t.Parallel()
//	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
//	require.NoError(t, err)
//
//	sqlxDB := sqlx.NewDb(db, "slqmock")
//	defer sqlxDB.Close()
//
//	accRepo := repository.NewAccountRepository(sqlxDB)
//
//	accRepo = accRepo
//
//	acc_uuid := uuid.New()
//	acc_culc_number := "01234567890123456789"
//	acc_corr_number := "01234567890123456789"
//	acc_bic := "123456789"
//	acc_cio := "123456789"
//	var acc_money_value uint = 0
//
//	t.Run("Success", func(t *testing.T) {
//
//		row := sqlmock.NewRows([]string{"acc_uuid", "acc_culc_number", "acc_corr_number", "acc_bic", "acc_cio", "acc_money_value"}).AddRow(
//			&acc_uuid,
//			&acc_culc_number,
//			&acc_corr_number,
//			&acc_bic,
//			&acc_cio,
//			&acc_money_value,
//		)
//
//		mock.ExpectBegin()
//		mock.ExpectExec(repository.CreateAccount).WithArgs(acc_uuid, acc_culc_number, acc_corr_number, acc_bic, acc_cio, acc_money_value).WillReturnResult(sqlmock.NewResult(1, 1))
//		mock.ExpectQuery(repository.DeleteAccount).WithArgs(
//			&acc_uuid,
//		).WillReturnRows(row)
//		mock.ExpectCommit()
//
//		err = accRepo.DeleteAccount(context.Background(), acc_uuid)
//		require.Nil(t, err)
//	})
//}
