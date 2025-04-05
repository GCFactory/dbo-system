package usecase

import "errors"

var (
	ErrorUserExists           = errors.New("User already exists")
	ErrorUserNotFound         = errors.New("User not found")
	ErrorUnMarshal            = errors.New("Error un-marshalling")
	ErrorAccountAlreadyExists = errors.New("Account already exists")
	ErrorAccountWasNotFound   = errors.New("Account wasn't found")
	ErrorMarshal              = errors.New("Error marshalling")
	ErrorWrongPassword        = errors.New("Error wrong password")
)
