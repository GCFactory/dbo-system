package repository

import (
	"context"
	"encoding/json"
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

	tx, err := repo.db.Begin()
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

	accounts, err := json.Marshal(user.User_accounts)

	if _, err = repo.db.ExecContext(local_ctx,
		AddUser,
		user.User_uuid,
		user.Passport_uuid,
		user.User_inn,
		accounts,
		user.User_login,
		user.User_passw,
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

	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var user models.User

	var passport models.Passport

	var tmp []byte

	if err = repo.db.QueryRowxContext(local_ctx,
		GetUserData,
		&user_uuid,
	).Scan(&user.User_uuid, &user.Passport_uuid, &user.User_inn, &tmp, &user.User_login, &user.User_passw, &user.UsingTotp, &user.TotpId); err != nil {
		return nil, ErrorGetUser
	}

	if err = json.Unmarshal(tmp, &user.User_accounts); err != nil {
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

	tx, err := repo.db.Begin()
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

	tx, err := repo.db.Begin()
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

	tx, err := repo.db.Begin()
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	var result models.Accounts

	if err = repo.db.QueryRowxContext(local_ctx,
		GetUsersAccounts,
		&user_uuid,
	).StructScan(&result); err != nil {
		return "", ErrorGetUsersAccounts
	}

	if err = tx.Commit(); err != nil {
		return "", err
	}

	return result.User_accounts, nil
}

func (repo UserRepository) UpdateUserPassw(ctx context.Context, user_uuid uuid.UUID, new_passw string) error {

	span, local_ctx := opentracing.StartSpanFromContext(ctx, "UserRepository.UpdateUserPassw")
	defer span.Finish()

	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	res, err := repo.db.ExecContext(local_ctx,
		UpdateUserPassw,
		user_uuid,
		new_passw,
	)

	if err != nil {
		return ErrorUpdatePassword
	} else if count, err := res.RowsAffected(); err != nil || count == 0 {
		return ErrorUpdatePassword
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (repo UserRepository) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {

	span, local_ctx := opentracing.StartSpanFromContext(ctx, "UserRepository.GetUserByLogin")
	defer span.Finish()

	var result = &models.User{}

	tx, err := repo.db.Begin()
	if err != nil {
		return result, err
	}
	defer tx.Rollback()

	var tmp []byte

	if err = repo.db.QueryRowxContext(local_ctx,
		GetUserByLogin,
		&login,
	).Scan(&result.User_uuid, &result.Passport_uuid, &result.User_inn, &tmp, &result.User_login, &result.User_passw, &result.UsingTotp, &result.TotpId); err != nil {
		return nil, ErrorGetUser
	}

	if err = json.Unmarshal(tmp, &result.User_accounts); err != nil {
		return nil, ErrorGetUser
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return result, nil
}

func (repo UserRepository) DeleteUser(ctx context.Context, userId uuid.UUID) error {

	span, local_ctx := opentracing.StartSpanFromContext(ctx, "UserRepository.DeleteUser")
	defer span.Finish()

	res, err := repo.db.ExecContext(local_ctx,
		DeleteUser,
		userId,
	)

	if err != nil {
		return err
	} else if count, err := res.RowsAffected(); err != nil || count == 0 {
		return ErrorNoUserFound
	}

	return nil
}

func (repo UserRepository) UpdateUserTotp(ctx context.Context, userUuid uuid.UUID, totpId uuid.UUID, totpUsage bool) error {

	span, local_ctx := opentracing.StartSpanFromContext(ctx, "UserRepository.UpdateUserTotp")
	defer span.Finish()

	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	res, err := repo.db.ExecContext(local_ctx,
		UpdateUserTotp,
		userUuid,
		totpId,
		totpUsage,
	)

	if err != nil {
		return ErrorUpdateTotpInfo
	} else if count, err := res.RowsAffected(); err != nil || count == 0 {
		return ErrorUpdateTotpInfo
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil

}

func NewUserRepository(db *sqlx.DB) users.Repository {
	return &UserRepository{db: db}
}
