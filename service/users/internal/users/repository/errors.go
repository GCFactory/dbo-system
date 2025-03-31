package repository

import "errors"

var (
	ErrorAddPassport      = errors.New("Error adding passport data")
	ErrorAddUser          = errors.New("Error adding user data")
	ErrorGetPassport      = errors.New("Error getting passport")
	ErrorGetUser          = errors.New("Error getting user")
	ErrorUpdatePassport   = errors.New("Error updating passport")
	ErrorUpdateAccounts   = errors.New("Error updating accounts")
	ErrorGetUsersAccounts = errors.New("Error getting users accounts")
	ErrorUpdatePassword   = errors.New("Error update password")
	ErrorNoUserFound      = errors.New("No user found")
)
