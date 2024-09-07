package repository

import (
	"context"
	"github.com/GCFactory/dbo-system/service/users/internal/models"
	"github.com/GCFactory/dbo-system/service/users/internal/users"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
)

type UserRepository struct {
	db *sqlx.DB
}

func (repo UserRepository) AddUser(ctx context.Context, user_data *models.User_full_data) error {
	span, local_ctx := opentracing.StartSpanFromContext(ctx, "UserRepository.AddUser")
	defer span.Finish()

	tx, err := repo.db.BeginTx(local_ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	passport := user_data.Passport

	if _, err = repo.db.ExecContext(local_ctx,
		AddPassport,
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
	); err != nil {
		return ErrorAddPassport
	}

	user := user_data.User
	if _, err = repo.db.ExecContext(local_ctx,
		AddUser,
		user.User_uuid,
		user.Passport_uuid,
		user.User_inn,
		user.User_accounts,
	); err != nil {
		return ErrorAddUser
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (repo UserRepository) GetUserData(ctx context.Context, user_uuid uuid.UUID) (*models.User_full_data, error) {

	span, local_ctx := opentracing.StartSpanFromContext(ctx, "UserRepository.GetUserData")
	defer span.Finish()

	tx, err := repo.db.BeginTx(local_ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var user models.User

	var passport models.Passport

	if err = repo.db.QueryRowxContext(local_ctx,
		GetUserData,
		&user_uuid,
	).StructScan(&user); err != nil {
		return nil, ErrorGetUser
	}

	if err = repo.db.QueryRowxContext(local_ctx,
		GetPassportData,
		&user.Passport_uuid,
	).StructScan(&passport); err != nil {
		return nil, ErrorGetPassport
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &models.User_full_data{
		User:     &user,
		Passport: &passport,
	}, nil
}

func (repo UserRepository) UpdatePassport(ctx context.Context, passport *models.Passport) error {

	span, local_ctx := opentracing.StartSpanFromContext(ctx, "UserRepository.UpdatePassport")
	defer span.Finish()

	tx, err := repo.db.BeginTx(local_ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	res, err := repo.db.ExecContext(local_ctx,
		UpdatePassport,
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

	if err != nil {
		return ErrorUpdatePassport
	} else if count, err := res.RowsAffected(); err != nil || count == 0 {
		return ErrorUpdatePassport
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (repo UserRepository) UpdateUserAccount(ctx context.Context, user_uuid uuid.UUID, accounts string) error {

	span, local_ctx := opentracing.StartSpanFromContext(ctx, "UserRepository.UpdateUserAccount")
	defer span.Finish()

	tx, err := repo.db.BeginTx(local_ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	res, err := repo.db.ExecContext(local_ctx,
		UpdateUsersAccounts,
		user_uuid,
		accounts,
	)

	if err != nil {
		return ErrorUpdateAccounts
	} else if count, err := res.RowsAffected(); err != nil || count == 0 {
		return ErrorUpdateAccounts
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (repo UserRepository) GetUsersAccounts(ctx context.Context, user_uuid uuid.UUID) (string, error) {

	span, local_ctx := opentracing.StartSpanFromContext(ctx, "UserRepository.GetUsersAccounts")
	defer span.Finish()

	tx, err := repo.db.BeginTx(local_ctx, nil)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	var result string

	if err = repo.db.QueryRowxContext(local_ctx,
		GetUsersAccounts,
		&user_uuid,
	).StructScan(&result); err != nil {
		return "", ErrorGetUsersAccounts
	}

	if err = tx.Commit(); err != nil {
		return "", err
	}

	return result, nil
}

func NewUserRepository(db *sqlx.DB) users.Repository {
	return &UserRepository{db: db}
}
