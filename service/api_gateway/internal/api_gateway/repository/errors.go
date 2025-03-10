package repository

import "errors"

var (
	ErrorAddToken          = errors.New("Error adding repository")
	ErrorGetTokenValue     = errors.New("Error get token value")
	ErrorGetTokenExpire    = errors.New("Error get token expire")
	ErrorUpdateTokenExpire = errors.New("Error update token expire")
	ErrorDeleteToken       = errors.New("Error delete token")
)
