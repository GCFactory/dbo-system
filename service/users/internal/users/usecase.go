package users

import (
	"context"
	"github.com/GCFactory/dbo-system/service/users/internal/models"
	"github.com/google/uuid"
)

type UseCase interface {
	AddUser(ctx context.Context, user *models.User_full_data) error
	GetUserData(ctx context.Context, user_uuid uuid.UUID) (*models.User_full_data, error)
	UpdateUsersPassport(ctx context.Context, user_uuid uuid.UUID, passport *models.Passport) error
	RemoveUser(ctx context.Context, userId uuid.UUID) error
	AddUserAccount(ctx context.Context, user_uuid uuid.UUID, account uuid.UUID) error
	RemoveUserAccount(ctx context.Context, user_id uuid.UUID, account uuid.UUID) error
	GetUserAccounts(ctx context.Context, user_uuid uuid.UUID) (*models.ListOfAccounts, error)
	UpdateUserPassword(ctx context.Context, user_uuid uuid.UUID, passw string) error
	GetUserDataByLogin(ctx context.Context, login string) (*models.User, error)
	CheckUserPassword(ctx context.Context, user_uuid uuid.UUID, passw string) error
}
