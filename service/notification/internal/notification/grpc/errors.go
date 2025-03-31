package grpc

import (
	"errors"
)

var (
	ErrorNoUserEmail          = errors.New("No user email")
	ErrorUnknownOperationType = errors.New("unknown operation type")
	ErrorInvalidInputData     = errors.New("invalid input data")
)

var GRPCErrors = map[error]uint32{

	//repository.ErrorGetUsersAccounts: 10,
	//repository.ErrorUpdateAccounts:   20,
	//repository.ErrorUpdatePassport:   30,
	//repository.ErrorGetUser:          40,
	//repository.ErrorGetPassport:      50,
	//repository.ErrorAddUser:          60,
	//repository.ErrorAddPassport:      70,
	//repository.ErrorUpdatePassword:   80,
	//
	//usecase.ErrorAccountAlreadyExists: 200,
	//usecase.ErrorUserExists:           210,
	//usecase.ErrorUserNotFound:         220,
	//usecase.ErrorMarshal:              300,
	//usecase.ErrorUnMarshal:            310,
	//
	//ErrorInvalidInputData:     400,
	//ErrorUnkonwnOperationType: 410,
}

func GetErrorCode(err error) uint32 {
	var result uint32
	result, ok := GRPCErrors[err]

	if !ok {
		result = 4294967295
	}

	return result
}
