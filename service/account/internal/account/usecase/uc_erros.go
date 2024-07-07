package usecase

import "errors"

var (
	ErrorWrongAccStatus   = errors.New("Wrong account status!")
	ErrorWrongMoneyValue  = errors.New("Wrong money value!")
	ErrorAccIsExisting    = errors.New("Account for this acc uuid is existing!")
	ErrorReasonIsExisting = errors.New("Reason for this acc uuid is existing!")
	ErrorCreateAcc        = errors.New("accountUC.ReservAcc.CreateAccount")
	ErrorAddReason        = errors.New("accountUC.ReservAcc.AddReserveReason")
	ErrorFatal            = errors.New("Fatal error")
)
