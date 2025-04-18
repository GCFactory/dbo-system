//go:generate mockgen -source repository.go -destination mock/repository_mock.go -package mock

package users

import (
	"context"
	"github.com/GCFactory/dbo-system/service/users/internal/models"
	"github.com/google/uuid"
)

type Repository interface {
	AddUser(ctx context.Context, user_data *models.User_full_data) error
	GetUserData(ctx context.Context, user_uuid uuid.UUID) (*models.User_full_data, error)
	UpdatePassport(ctx context.Context, passport *models.Passport) error
	DeleteUser(ctx context.Context, userId uuid.UUID) error
	UpdateUserAccount(ctx context.Context, user_uuid uuid.UUID, accounts string) error
	GetUsersAccounts(ctx context.Context, user_uuid uuid.UUID) (string, error)
	UpdateUserPassw(ctx context.Context, user_uuid uuid.UUID, new_passw string) error
	GetUserByLogin(ctx context.Context, login string) (*models.User, error)
	UpdateUserTotp(ctx context.Context, userUuid uuid.UUID, totpId uuid.UUID, totpUsage bool) error
}
