package repository

const (
	createConfig = `INSERT INTO totp_service.totp_codes (totp_id, user_id, is_active, issuer, account_name, secret, url, created_at)
						VALUES ($1, $2, $3, $4, $5, $6, $7, now())
						RETURNING *`
	getActiveConfig   = `SELECT * from totp_service.totp_codes WHERE user_id = $1 and is_active = true`
	getConfigByTotpId = `SELECT * from totp_service.totp_codes WHERE totp_id = $1`
	getConfigByUserId = `SELECT * from totp_service.totp_codes WHERE user_id = $1`

	updateTotpActivityByTotpId = `UPDATE totp_service.totp_codes SET is_active = $2, updated_at = now() WHERE totp_id = $1`
)
