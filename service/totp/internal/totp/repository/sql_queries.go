package repository

const (
	createConfig = `INSERT INTO totp_service.totp_codes (totp_id, user_id, is_active, issuer, account_name, secret, url, created_at)
						VALUES ($1, $2, $3, $4, $5, $6, $7, now())
						RETURNING *`
)
