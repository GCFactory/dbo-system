package repository

import "errors"

var (
	ErrorCreateAccount       = errors.New("accountRepo.CreateAccount.QueryRowxContext")
	ErrorGetAccountData      = errors.New("accountRepo.GetAccountData.QueryRowxContext")
	ErrorGetAccountStatus    = errors.New("accountRepo.GetAccountStatus.QueryRowxContext")
	ErrorGetAccountAmount    = errors.New("accountRepo.GetAccountAmount.QueryRowxContext")
	ErrorGetReserveReason    = errors.New("accountRepo.GetReserveReason.QueryRowxContext")
	ErrorUpdateAccountStatus = errors.New("accountRepo.UpdateAccountStatus.ExecContext")
	ErrorAddReserveReason    = errors.New("accountRepo.AddReserveReason.QueryRowxContext")
	ErrorUpdateAccountAmount = errors.New("accountRepo.UpdateAccountAmount.ExecContext")
	ErrorDeleteAccount       = errors.New("accountRepo.DeleteAccount.ExecContext")
	ErrorDeleteReserveReason = errors.New("accountRepo.DeleteReserveReason.ExecContext")
)
