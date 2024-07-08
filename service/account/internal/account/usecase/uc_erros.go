package usecase

import "errors"

var (
	ErrorWrongAccStatus         = errors.New("Wrong account status!")
	ErrorWrongMoneyValue        = errors.New("Wrong money value!")
	ErrorAccIsExisting          = errors.New("Account for this acc uuid is existing!")
	ErrorReasonIsExisting       = errors.New("Reason for this acc uuid is existing!")
	ErrorCreateAcc              = errors.New("accountRepo.CreateAccount")
	ErrorAddReason              = errors.New("accountRepo.AddReserveReason")
	ErrorFatal                  = errors.New("Fatal error")
	ErrorNoFoundAcc             = errors.New("No account found!")
	ErrorWrongAccReservedStatus = errors.New("Acc has wrong status for creating it!")
	ErrorWrongAccOpenStatus     = errors.New("Acc has wrong status for opening it!")
	ErrorWrongAccCloseStatus    = errors.New("Acc has wrong status for closing it!")
	ErrorWrongAccBlockStatus    = errors.New("Acc has wrong status for blocking it!")
	ErrorUpdateAccStatus        = errors.New("accountRepo.UpdateAccountStatus")
	ErrorGetAccReservedReason   = errors.New("accountRepo.GetReserveReason")
	ErrorOverflowAmount         = errors.New("Amount overflow!")
	ErrorUpdateAmountValue      = errors.New("accountRepo.UpdateAccountAmount")
	ErrorNotEnoughMoneyAmount   = errors.New("Not enough money amount!")
)
