package repository

import (
	"context"
	registration "github.com/GCFactory/dbo-system/service/account/internal/account"
	"github.com/GCFactory/dbo-system/service/account/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
)

type accountRepo struct {
	db *sqlx.DB
}

func (repo accountRepo) CreateAccount(ctx context.Context, account *models.Account, reason *models.ReserverReason) error {
	span, local_ctx := opentracing.StartSpanFromContext(ctx, "accountRepo.CreateAccount")
	defer span.Finish()

	tx, err := repo.db.BeginTx(local_ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if _, err = repo.db.ExecContext(local_ctx,
		CreateAccount,
		account.Acc_uuid,
		account.Acc_culc_number,
		account.Acc_corr_number,
		account.Acc_bic,
		account.Acc_cio,
		account.Acc_money_value,
	); err != nil {
		return ErrorCreateAccount
	}

	if _, err = repo.db.ExecContext(local_ctx,
		InsertReserveReason,
		reason.Acc_uuid,
		reason.Reason,
	); err != nil {
		return ErrorAddReserveReason
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (repo accountRepo) DeleteAccount(ctx context.Context, acc_uuid uuid.UUID) error {
	span, local_ctx := opentracing.StartSpanFromContext(ctx, "accountRepo.DeleteAccount")
	defer span.Finish()

	res, err := repo.db.ExecContext(local_ctx,
		DeleteAccount,
		acc_uuid,
	)

	if err == nil {
		count, err := res.RowsAffected()

		if err != nil || count == 0 {
			return ErrorDeleteAccount
		}

		return nil
	}

	return ErrorDeleteAccount
}

func (repo accountRepo) GetAccountData(ctx context.Context, acc_uuid uuid.UUID) (*models.Account, error) {
	span, local_ctx := opentracing.StartSpanFromContext(ctx, "accountRepo.GetAccountData")
	defer span.Finish()

	var result models.Account

	if err := repo.db.QueryRowxContext(local_ctx,
		GetAccountData,
		&acc_uuid,
	).StructScan(&result); err != nil {
		return nil, ErrorGetAccountData
	}

	return &result, nil
}

func (repo accountRepo) UpdateAccountStatus(ctx context.Context, acc_uuid uuid.UUID, new_status uint8) error {
	span, local_ctx := opentracing.StartSpanFromContext(ctx, "accountRepo.UpdateAccountStatus")
	defer span.Finish()

	res, err := repo.db.ExecContext(local_ctx,
		UpdateAccountStatus,
		acc_uuid,
		new_status,
	)

	if err != nil {
		return ErrorUpdateAccountStatus
	} else {
		count, err := res.RowsAffected()
		if err != nil || count == 0 {
			return ErrorUpdateAccountStatus
		}
	}

	return nil
}

func (repo accountRepo) GetAccountStatus(ctx context.Context, acc_uuid uuid.UUID) (uint8, error) {
	span, local_ctx := opentracing.StartSpanFromContext(ctx, "accountRepo.GetAccountStatus")
	defer span.Finish()

	var result models.Account

	if err := repo.db.QueryRowxContext(local_ctx,
		GetAccountStatus,
		&acc_uuid,
	).StructScan(&result); err != nil {
		return 0, ErrorGetAccountStatus
	}

	return result.Acc_status, nil
}

func (repo accountRepo) UpdateAccountAmount(ctx context.Context, acc_uuid uuid.UUID, acc_new_amount float64) error {
	span, local_ctx := opentracing.StartSpanFromContext(ctx, "accountRepo.UpdateAccountAmount")
	defer span.Finish()

	res, err := repo.db.ExecContext(local_ctx,
		UpdateAccountAmount,
		acc_uuid,
		acc_new_amount,
	)
	if err != nil {
		return ErrorUpdateAccountAmount
	} else {
		count, err := res.RowsAffected()
		if err != nil || count == 0 {
			return ErrorUpdateAccountAmount
		}
	}

	return nil
}

func (repo accountRepo) GetAccountAmount(ctx context.Context, acc_uuid uuid.UUID) (float64, error) {
	span, local_ctx := opentracing.StartSpanFromContext(ctx, "accountRepo.GetAccountAmount")
	defer span.Finish()

	var result models.Account

	if err := repo.db.QueryRowxContext(local_ctx,
		GetAccountAmount,
		&acc_uuid,
	).StructScan(&result); err != nil {
		return 0, ErrorGetAccountAmount
	}

	return result.Acc_money_amount, nil
}

func (repo accountRepo) DeleteReserveReason(ctx context.Context, acc_uuid uuid.UUID) error {
	span, local_ctx := opentracing.StartSpanFromContext(ctx, "accountRepo.DeleteReserveReason")
	defer span.Finish()

	res, err := repo.db.ExecContext(local_ctx,
		DeleteReserveReason,
		acc_uuid,
	)

	if err == nil {
		count, err := res.RowsAffected()

		if err != nil || count == 0 {
			return ErrorDeleteReserveReason
		}

		return nil
	}

	return ErrorDeleteReserveReason
}

func (repo accountRepo) GetReserveReason(ctx context.Context, acc_uuid uuid.UUID) (*models.ReserverReason, error) {
	span, local_ctx := opentracing.StartSpanFromContext(ctx, "accountRepo.GetReserveReason")
	defer span.Finish()

	var result models.ReserverReason

	if err := repo.db.QueryRowxContext(local_ctx,
		GetReserveReason,
		&acc_uuid,
	).StructScan(&result); err != nil {
		return nil, ErrorGetReserveReason
	}

	return &result, nil
}

func NewAccountRepository(db *sqlx.DB) registration.Repository {
	return &accountRepo{db: db}
}
