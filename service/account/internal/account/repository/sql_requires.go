package repository

const (
	CreateAccount = `INSERT INTO accounts (acc_uuid, acc_culc_number, acc_corr_number, acc_bic, acc_cio, acc_money_value)
		VALUES($1, $2, $3, $4, $5, $6)
		RETURNING *`
	DeleteAccount       = `DELETE FROM accounts WHERE acc_uuid = $1`
	GetAccountData      = `SELECT * FROM ONLY accounts WHERE acc_uuid = $1`
	UpdateAccountStatus = `UPDATE ONLY accounts SET acc_status = $2 WHERE acc_uuid = $1`
	GetAccountStatus    = `SELECT acc_status FROM ONLY accounts WHERE acc_uuid = $1`
	UpdateAccountAmount = `UPDATE ONLY accounts SET acc_money_amount = $2 WHERE acc_uuid = $1`
	GetAccountAmount    = `SELECT acc_money_amount FROM ONLY accounts WHERE acc_uuid = $1`
	AddReserveReason    = `INSERT INTO accounts_reserved (acc_uuid, reserve_reason)
		VALUES($1, $2)
		RETURNING *`
	GetReserveReason = `SELECT reserve_reason FROM accounts_reserved WHERE acc_uuid = $1`
)
