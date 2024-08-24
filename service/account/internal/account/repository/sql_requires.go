package repository

const (
	CreateAccount = `INSERT INTO accounts (acc_uuid, acc_status, acc_culc_number, acc_corr_number, acc_bic, acc_cio, acc_money_value)
			VALUES($1, $2, $3, $4, $5, $6, $7);`
	InsertReserveReason = `INSERT INTO accounts_reserved (acc_uuid, reserve_reason)
			VALUES($1, $2);`
	DeleteAccount       = `DELETE FROM accounts WHERE acc_uuid = $1`
	GetAccountData      = `SELECT * FROM ONLY accounts WHERE acc_uuid = $1`
	UpdateAccountStatus = `UPDATE ONLY accounts SET acc_status = $2 WHERE acc_uuid = $1`
	GetAccountStatus    = `SELECT acc_status FROM ONLY accounts WHERE acc_uuid = $1`
	UpdateAccountAmount = `UPDATE ONLY accounts SET acc_money_amount = $2 WHERE acc_uuid = $1`
	GetAccountAmount    = `SELECT acc_money_amount FROM ONLY accounts WHERE acc_uuid = $1`
	GetReserveReason    = `SELECT * FROM accounts_reserved WHERE acc_uuid = $1`
	DeleteReserveReason = `DELETE FROM accounts_reserved WHERE acc_uuid = $1`
)
