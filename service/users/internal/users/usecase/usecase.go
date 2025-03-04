package usecase

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"github.com/GCFactory/dbo-system/platform/config"
	"github.com/GCFactory/dbo-system/platform/pkg/logger"
	"github.com/GCFactory/dbo-system/service/users/internal/models"
	"github.com/GCFactory/dbo-system/service/users/internal/users"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
)

type userUsecase struct {
	cfg       *config.Config
	usersRepo users.Repository
	logger    logger.Logger
}

func (uc userUsecase) AddUser(ctx context.Context, user *models.User_full_data) error {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "userUsecase.AddUser")
	defer span.Finish()

	sha := sha512.Sum512([]byte(user.User.User_passw))
	user.User.User_passw = hex.EncodeToString(sha[:])

	bd_user, _ := uc.usersRepo.GetUserData(ctxWithTrace, user.User.User_uuid)
	if bd_user != nil {
		return ErrorUserExists
	}

	if err := uc.usersRepo.AddUser(ctxWithTrace, user); err != nil {
		return err
	}

	return nil
}

func (uc userUsecase) GetUserData(ctx context.Context, user_uuid uuid.UUID) (*models.User_full_data, error) {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "userUsecase.GetUserData")
	defer span.Finish()

	result, err := uc.usersRepo.GetUserData(ctxWithTrace, user_uuid)
	if err != nil {
		return nil, ErrorUserNotFound
	}

	return result, nil
}

func (uc userUsecase) UpdateUsersPassport(ctx context.Context, user_uuid uuid.UUID, passport *models.Passport) error {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "userUsecase.UpdateUsersPassport")
	defer span.Finish()

	user, _ := uc.usersRepo.GetUserData(ctxWithTrace, user_uuid)
	if user == nil {
		return ErrorUserNotFound
	}

	passport.Passport_uuid = user.Passport.Passport_uuid

	err := uc.usersRepo.UpdatePassport(ctxWithTrace, passport)
	if err != nil {
		return err
	}

	return nil
}

func (uc userUsecase) AddUserAccount(ctx context.Context, user_uuid uuid.UUID, account uuid.UUID) error {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "userUsecase.AddUserAccount")
	defer span.Finish()

	user, _ := uc.usersRepo.GetUserData(ctxWithTrace, user_uuid)
	if user == nil {
		return ErrorUserNotFound
	}

	list_of_accounts := user.User.User_accounts.Data

	for i := 0; i < len(list_of_accounts); i++ {
		if list_of_accounts[i] == account {
			return ErrorAccountAlreadyExists
		}
	}

	list_of_accounts = append(list_of_accounts, account)
	user.User.User_accounts.Data = list_of_accounts

	bytes, err := json.Marshal(user.User.User_accounts)
	if err != nil {
		return ErrorMarshal
	}

	err = uc.usersRepo.UpdateUserAccount(ctxWithTrace, user_uuid, string(bytes))
	if err != nil {
		return err
	}

	return nil
}

func (uc userUsecase) RemoveUserAccount(ctx context.Context, user_uuid uuid.UUID, account uuid.UUID) error {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "userUsecase.RemoveUserAccount")
	defer span.Finish()

	user, _ := uc.usersRepo.GetUserData(ctxWithTrace, user_uuid)
	if user == nil {
		return ErrorUserNotFound
	}

	list_of_accounts := user.User.User_accounts.Data

	is_found := false
	account_number := -1
	for acc_number, account_uuid := range list_of_accounts {
		if account_uuid == account {
			account_number = acc_number
			is_found = true
			break
		}
	}

	if is_found {
		list_of_accounts = append(list_of_accounts[:account_number], list_of_accounts[account_number+1:]...)
		user.User.User_accounts.Data = list_of_accounts

		bytes, err := json.Marshal(user.User.User_accounts)
		if err != nil {
			return ErrorMarshal
		}

		err = uc.usersRepo.UpdateUserAccount(ctxWithTrace, user_uuid, string(bytes))
		if err != nil {
			return err
		}
	} else {
		return ErrorAccountWasNotFound
	}

	return nil
}

func (uc userUsecase) GetUserAccounts(ctx context.Context, user_uuid uuid.UUID) (*models.ListOfAccounts, error) {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "userUsecase.GetUserAccounts")
	defer span.Finish()

	accounts, _ := uc.usersRepo.GetUsersAccounts(ctxWithTrace, user_uuid)
	if accounts == "" {
		return nil, ErrorUserNotFound
	}

	var result *models.ListOfAccounts

	err := json.Unmarshal([]byte(accounts), &result)
	if err != nil {
		return nil, ErrorUnMarshal
	}

	return result, nil
}

func (uc userUsecase) UpdateUserPassword(ctx context.Context, user_uuid uuid.UUID, passw string) error {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "userUsecase.UpdateUserPassword")
	defer span.Finish()

	user, err := uc.usersRepo.GetUserData(ctxWithTrace, user_uuid)
	if err != nil {
		return ErrorUserNotFound
	}

	if user != nil {
		err = uc.usersRepo.UpdateUserPassw(ctxWithTrace, user_uuid, passw)
		if err != nil {
			return err
		}
	}

	return nil
}

func (uc userUsecase) GetUserDataByLogin(ctx context.Context, login string) (*models.User, error) {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "userUsecase.GetUserDataByLogin")
	defer span.Finish()

	user, err := uc.usersRepo.GetUserByLogin(ctxWithTrace, login)
	if err != nil {
		return nil, ErrorUserNotFound
	}

	return user, nil
}

func NewUsersUseCase(cfg *config.Config, users_repo users.Repository, log logger.Logger) users.UseCase {
	return &userUsecase{cfg: cfg, usersRepo: users_repo, logger: log}
}
