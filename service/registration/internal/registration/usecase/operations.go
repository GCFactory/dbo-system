package usecase

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
}
