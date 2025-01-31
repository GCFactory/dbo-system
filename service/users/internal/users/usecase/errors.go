package usecase

import "errors"

var (
	ErrorUserExists           = errors.New("User already exists")
	ErrorUserNotFound         = errors.New("User not found")
	ErrorUnMarshal            = errors.New("Error un-marshalling")
	ErrorAccountAlreadyExists = errors.New("Account already exists")
	ErrorMarshal              = errors.New("Error marshalling")
)
