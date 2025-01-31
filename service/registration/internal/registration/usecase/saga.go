package usecase

import "github.com/GCFactory/dbo-system/service/registration/internal/models"

const (
	SagaStatusUndefined         uint = 0
	SagaStatusCreated           uint = 10
	SagaStatusInProcess         uint = 20
	SagaStatusCompleted         uint = 30
	SagaStatusFallBackInProcess uint = 40
	SagaStatusFallBackSuccess   uint = 50
	SagaStatusFallBackError     uint = 60
	SagaStatusError             uint = 255
)

const (
	SagaTypeCreateUser              string = "create_user"
	SagaTypeCheckUser               string = "check_user"
	SagaTypeReserveAccount          string = "reserve_account"
	SagaTypeCreateAccount           string = "create_account"
	SagaTypeOpenAccountAndAddToUser string = "open_account_and_add_to_user"
)

var PossibleSagaTypes = []string{
	SagaTypeCreateUser,
	SagaTypeCheckUser,
	SagaTypeReserveAccount,
	SagaTypeCreateAccount,
	SagaTypeOpenAccountAndAddToUser,
}

func ValidateSagaType(saga_type string) bool {

	for i := 0; i < len(PossibleSagaTypes); i++ {
		if PossibleSagaTypes[i] == saga_type {
			return true
		}
	}

	return false

}

const (
	SagaGroupCreateUser    uint8 = 1
	SagaGroupCreateAccount uint8 = 2
)

var PossibleSagaGroups = []uint8{
	SagaGroupCreateUser,
	SagaGroupCreateAccount,
}

func ValidateSagaGroup(saga_group uint8) bool {
	for i := 0; i < len(PossibleSagaGroups); i++ {
		if saga_group == PossibleSagaGroups[i] {
			return true
		}
	}

	return false
}

var PossibleEventsListForSagaType = map[string][]string{
	SagaTypeCreateUser: {
		EventTypeCreateUser,
	},
	SagaTypeCheckUser: {
		EventTypeGetUserData,
	},
	SagaTypeReserveAccount: {
		EventTypeReserveAccount,
	},
	SagaTypeCreateAccount: {
		EventTypeCreateAccount,
	},
	SagaTypeOpenAccountAndAddToUser: {
		EventTypeOpenAccount,
		EventTypeAddAccountToUser,
	},
}

var ListOfSagaDepend = map[string]models.SagaDepend{
	SagaTypeCreateUser: {
		Parents:  nil,
		Children: nil,
	},
	SagaTypeCheckUser: models.SagaDepend{
		Parents: nil,
		Children: []string{
			SagaTypeReserveAccount,
		},
	},
	SagaTypeReserveAccount: models.SagaDepend{
		Parents: []string{
			SagaTypeCheckUser,
		},
		Children: []string{
			SagaTypeCreateAccount,
		},
	},
	SagaTypeCreateAccount: models.SagaDepend{
		Parents: []string{
			SagaTypeReserveAccount,
		},
		Children: []string{
			SagaTypeOpenAccountAndAddToUser,
		},
	},
	SagaTypeOpenAccountAndAddToUser: models.SagaDepend{
		Parents: []string{
			SagaTypeCreateAccount,
		},
		Children: nil,
	},
}

const (
	SagaConnectionStatusUnknown  uint8 = 0
	SagaConnectionStatusWaiting  uint8 = 10
	SagaConnectionStatusSuccess  uint8 = 20
	SagaConnectionStatusFallBack uint8 = 30
	SagaConnectionStatusFailed   uint8 = 255
)

var PossibleSagaConnectionStatus = []uint8{
	SagaConnectionStatusUnknown,
	SagaConnectionStatusWaiting,
	SagaConnectionStatusSuccess,
	SagaConnectionStatusFallBack,
	SagaConnectionStatusFailed,
}

func ValidateSagaConnectionStatus(saga_connection_status uint8) bool {
	for i := 0; i < len(PossibleSagaConnectionStatus); i++ {
		if PossibleSagaConnectionStatus[i] == saga_connection_status {
			return true
		}
	}

	return false
}
