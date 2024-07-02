package repository

import (
	"context"
	registration "github.com/GCFactory/dbo-system/service/account/internal/account"
	"github.com/GCFactory/dbo-system/service/account/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type accountRepo struct {
	db *sqlx.DB
}

func (repo accountRepo) CreateAccount(ctx context.Context, account models.Account) error {
	//	todo: реализовать
	return nil
}

func (repo accountRepo) DeleteAccount(ctx context.Context, acc_uuid uuid.UUID) error {
	//	todo: реализовать
	return nil
}

func (repo accountRepo) GetAccountData(ctx context.Context, acc_uuid uuid.UUID) (*models.Account, error) {
	//	todo: реализовать
	return &models.Account{}, nil
}

func (repo accountRepo) UpdateAccountStatus(ctx context.Context, acc_uuid uuid.UUID) error {
	//	todo: реализовать
	return nil
}

func (repo accountRepo) GetAccountStatus(ctx context.Context, acc_uuid uuid.UUID) (uint8, error) {
	//	todo: реализовать
	return 0, nil
}

func (repo accountRepo) UpdateAccountAmount(ctx context.Context, acc_uuid uuid.UUID, acc_new_amount float32) error {
	//	todo: реализовать
	return nil
}

func (repo accountRepo) GetAccountAmount(ctx context.Context, acc_uuid uuid.UUID) (float64, error) {
	//	todo: реализовать
	return 0.0, nil
}

func (repo accountRepo) AddReserveReason(ctx context.Context, acc_uuid uuid.UUID, reason string) error {
	//	todo: реализовать
	return nil
}
func (repo accountRepo) GetReserveReason(ctx context.Context, acc_uuid uuid.UUID) (string, error) {
	//	todo: реализовать
	return "", nil
}

func NewAccountRepository(db *sqlx.DB) registration.Repository {
	return &accountRepo{db: db}
}
