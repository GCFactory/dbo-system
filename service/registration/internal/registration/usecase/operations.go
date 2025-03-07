package usecase

import (
	"errors"
	"slices"
)

const (
	OperationUnknown                 uint8 = 0
	OperationCreateUser              uint8 = 1
	OperationAddAccount              uint8 = 2
	OperationAddAccountCache         uint8 = 3
	OperationWidthAccountCache       uint8 = 4
	OperationCloseAccount            uint8 = 5
	OperationGetUserData             uint8 = 6
	OperationGetAccountData          uint8 = 7
	OperationGroupUpdateUserPassword uint8 = 8
	OperationGetUserDataByLogin      uint8 = 9
	OperationCheckUserPassword       uint8 = 10
	OperationError                   uint8 = 255
)

var PossibleOperations = []uint8{
	OperationUnknown,
	OperationCreateUser,
	OperationAddAccount,
	OperationAddAccountCache,
	OperationWidthAccountCache,
	OperationCloseAccount,
	OperationGetUserData,
	OperationGetAccountData,
	OperationGroupUpdateUserPassword,
	OperationGetUserDataByLogin,
	OperationCheckUserPassword,
	OperationError,
}

func ValidateOperation(operation uint8) bool {
	for _, op := range PossibleOperations {
		if op == operation {
			return true
		}
	}
	return false
}

// Список корневых SAG
var OperationsRootsSagas = map[uint8][]string{
	OperationCreateUser: []string{
		SagaTypeCreateUser,
	},
	OperationAddAccount: []string{
		SagaTypeCheckUser,
	},
	OperationAddAccountCache: []string{
		SagaTypeCheckUser,
	},
	OperationWidthAccountCache: {
		SagaTypeCheckUser,
	},
	OperationCloseAccount: {
		SagaTypeCheckUser,
	},
	OperationGetUserData: {
		SagaTypeCheckUser,
	},
	OperationGetAccountData: {
		SagaTypeCheckUser,
	},
	OperationGroupUpdateUserPassword: {
		SagaTypeUpdateUserPassword,
	},
	OperationGetUserDataByLogin: {
		SagaTypeGetUserDataByLogin,
	},
	OperationCheckUserPassword: {
		SagaTypeCheckUserPassword,
	},
}

// Имена операций
var OperationName = map[uint8]string{
	OperationCreateUser:              "create_user",
	OperationAddAccount:              "add_account",
	OperationAddAccountCache:         "add_account_cache",
	OperationWidthAccountCache:       "width_account_cache",
	OperationCloseAccount:            "close_account",
	OperationGetUserData:             "get_user_data",
	OperationGetAccountData:          "get_account_data",
	OperationGroupUpdateUserPassword: "update_user_password",
	OperationGetUserDataByLogin:      "get_user_data_by_login",
	OperationCheckUserPassword:       "check_user_password",
}

func OperationNameFromCode(operation_code uint8) (string, error) {

	result := ""
	var err error

	if slices.Contains(PossibleOperations, operation_code) {
		if name, ok := OperationName[operation_code]; ok {
			result = name
		}
	}

	if result == "" {
		err = errors.New("No operation name")
	}

	return result, err
}
