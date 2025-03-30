package repo

const (
	AddUserSettings = `INSERT INTO notification
						(
						 user_uuid,
						 email_usage,
						 email
						)
						VALUES ($1, $2, $3);`
	GetUserSettings = `SELECT *
						FROM ONLY notification
						WHERE user_uuid = $1;`
	UpdateUserSettings = `UPDATE ONLY notification
						SET email_usage = $2,
						    email = $3
						WHERE user_uuid = $1;`
	DeleteUserSettings = `DELETE FROM ONLY notification
							WHERE user_uuid = $1;`
)
