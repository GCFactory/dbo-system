package repository

const (
	AddPassport = `insert into pasport 
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
			user_accounts
		)
		values ($1, $2, $3, $4);`

	GetPassportData = `select * 
		from only passport 
		where passport_uuid = $1;`

	GetUserData = `select * 
		from users 
		where user_uuid = $1;`

	UpdatePassport = `update pasport 
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
		set users_accounts = $2 
		where user_uuid = $1;`

	GetUsersAccounts = `select users_accounts
		from only users
		where user_uuid = $1;`
)
