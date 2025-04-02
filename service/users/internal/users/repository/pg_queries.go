package repository

const (
	AddPassport = `insert into passport 
    	(
			passport_uuid, 
			passport_series, 
			passport_number, 
			name, 
			surname, 
			patronimic, 
			birth_date, 
			birth_location, 
			pick_up_point, 
			authority, 
			authority_date, 
			registration_adress
		) 
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);`

	AddUser = `insert into users 
    	(
			user_uuid,
			passport_uuid, 
			user_inn, 
			user_accounts,
    	 	user_login,
    	 	user_password
		)
		values ($1, $2, $3, $4, $5, $6);`

	DeleteUser = `DELETE FROM ONLY passport
					WHERE passport_uuid = 
					      (SELECT passport_uuid
					       FROM users
					       WHERE user_uuid = $1
					       LIMIT 1);`

	GetPassportData = `select * 
		from only passport 
		where passport_uuid = $1;`

	GetUserData = `select * 
		from users 
		where user_uuid = $1;`

	UpdatePassport = `update passport 
		set 
			passport_series = $2, 
			passport_number = $3, 
			name = $4, 
			surname = $5, 
			patronimic = $6, 
			birth_date = $7, 
			birth_location = $8, 
			pick_up_point = $9, 
			authority = $10, 
			authority_date = $11, 
			registration_adress = $12
		where passport_uuid = $1;`

	UpdateUsersAccounts = `update users 
		set user_accounts = $2 
		where user_uuid = $1;`

	GetUsersAccounts = `select user_accounts
		from only users
		where user_uuid = $1;`

	UpdateUserPassw = `UPDATE users
						SET user_password = $2
						WHERE user_uuid = $1;`

	GetUserByLogin = `SELECT *
						FROM users
						WHERE user_login = $1;`

	UpdateUserTotp = `UPDATE users
						SET totp_id = $2,
						    using_totp = $3
						WHERE user_uuid = $1;`
)
