package usecase

const (
	OperationUnknown    uint8 = 0
	OperationCreateUser uint8 = 1
	OperationAddAccount uint8 = 2
	OperationError      uint8 = 255
)

var PossibleOperations = []uint8{
	OperationUnknown,
	OperationCreateUser,
	OperationAddAccount,
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

var OperationsRootsSagas = map[uint8][]string{
	OperationCreateUser: []string{
		SagaTypeCreateUser,
	},
	OperationAddAccount: []string{
		SagaTypeReserveAccount,
	},
}
