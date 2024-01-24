package repository

const (
	createConfig = `INSERT INTO totp_service.totp_codes (totp_id, user_id, is_active, issuer, account_name, secret, url, created_at)
						VALUES ($1, $2, $3, $4, $5, $6, $7, now())
						RETURNING *`
	getURL       = `SELECT (url) FROM totp_service.totp_codes WHERE user_id = $1`
	updareActive = `UPDATE totp_service.totp_codes SET is_active = $2 WHERE user_id = $1`
)
