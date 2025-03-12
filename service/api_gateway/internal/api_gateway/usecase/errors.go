package usecase

import "errors"

var (
	ErrorTokenAlreadyExist     = errors.New("Token already exist")
	ErrorOperationProcessedYet = errors.New("Operation processed yet")
	ErrorWrongPassword         = errors.New("Wrong password")
	ErrorNoUserData            = errors.New("No user data")
	ErrorNoAccountData         = errors.New("No account data")
)
