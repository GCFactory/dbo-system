package usecase

import "errors"

var (
	ErrorTokenAlreadyExist     = errors.New("Token already exist")
	ErrorOperationProcessedYet = errors.New("Operation processed yet")
)
