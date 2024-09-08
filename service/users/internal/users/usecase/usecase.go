package usecase

import (
	"context"
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

	var result models.ListOfAccounts

	err := json.Unmarshal([]byte(user.User.User_accounts), &result)
	if err != nil {
		return ErrorUnMarshal
	}

	list_of_accounts := result.Data

	for i := 0; i < len(list_of_accounts); i++ {
		if list_of_accounts[i] == account {
			return ErrorAccountAlreadyExists
		}
	}

	list_of_accounts = append(list_of_accounts, account)

	result.Data = list_of_accounts

	bytes, err := json.Marshal(result)
	if err != nil {
		return ErrorMarshal
	}

	err = uc.usersRepo.UpdateUserAccount(ctxWithTrace, user_uuid, string(bytes))
	if err != nil {
		return err
	}

	return nil
}

func (uc userUsecase) GetUserAccounts(ctx context.Context, user_uuid uuid.UUID) (string, error) {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "userUsecase.GetUserAccounts")
	defer span.Finish()

	accounts, _ := uc.usersRepo.GetUsersAccounts(ctxWithTrace, user_uuid)
	if accounts == "" {
		return "", ErrorUserNotFound
	}

	return accounts, nil
}

func NewUsersUseCase(cfg *config.Config, users_repo users.Repository, log logger.Logger) users.UseCase {
	return &userUsecase{cfg: cfg, usersRepo: users_repo, logger: log}
}
