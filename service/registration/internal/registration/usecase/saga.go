package usecase

import "github.com/GCFactory/dbo-system/service/registration/internal/models"

const (
	SagaStatusUndefined         uint = 0
	SagaStatusCreated           uint = 10
	SagaStatusInProcess         uint = 20
	SagaStatusCompleted         uint = 30
	SagaStatusFallBackInProcess uint = 40
	SagaStatusFallBackSuccess   uint = 50
	SagaStatusFallBackError     uint = 250
	SagaStatusError             uint = 255
)

const (
	SagaTypeCreateUser              string = "create_user"
	SagaTypeCheckUser               string = "check_user"
	SagaTypeReserveAccount          string = "reserve_account"
	SagaTypeCreateAccount           string = "create_account"
	SagaTypeOpenAccountAndAddToUser string = "open_account_and_add_to_user"
	SagaTypeGetAccountData          string = "get_account_data"
	SagaTypeAddAccountCache         string = "add_account_cache"
	SagaTypeWidthAccountCache       string = "width_account_cache"
	SagaTypeCloseAccount            string = "close_account"
	SagaTypeUpdateUserPassword      string = "update_user_password"
	SagaTypeGetUserDataByLogin      string = "get_user_data_by_login"
	SagaTypeCheckUserPassword       string = "check_user_password"
)

var PossibleSagaTypes = []string{
	SagaTypeCreateUser,
	SagaTypeCheckUser,
	SagaTypeReserveAccount,
	SagaTypeCreateAccount,
	SagaTypeOpenAccountAndAddToUser,
	SagaTypeAddAccountCache,
	SagaTypeGetAccountData,
	SagaTypeWidthAccountCache,
	SagaTypeCloseAccount,
	SagaTypeUpdateUserPassword,
	SagaTypeGetUserDataByLogin,
	SagaTypeCheckUserPassword,
}

func ValidateSagaType(saga_type string) bool {

	for i := 0; i < len(PossibleSagaTypes); i++ {
		if PossibleSagaTypes[i] == saga_type {
			return true
		}
	}

	return false

}

// Список возможных операций в SAG-е
var PossibleEventsListForSagaType = map[string][]string{
	SagaTypeCreateUser: {
		EventTypeCreateUser,
	},
	SagaTypeCheckUser: {
		EventTypeGetUserData,
	},
	SagaTypeReserveAccount: {
		EventTypeReserveAccount,
		EventTypeRemoveAccount,
	},
	SagaTypeCreateAccount: {
		EventTypeCreateAccount,
		EventTypeRemoveAccount,
	},
	SagaTypeOpenAccountAndAddToUser: {
		EventTypeOpenAccount,
		EventTypeAddAccountToUser,
		EventTypeRemoveAccount,
		EventTypeRemoveUserAccount,
	},
	SagaTypeAddAccountCache: {
		EventTypeAddAccountCache,
	},
	SagaTypeGetAccountData: {
		EventTypeGetAccountData,
	},
	SagaTypeWidthAccountCache: {
		EventTypeWidthAccountCache,
	},
	SagaTypeCloseAccount: {
		EventTypeCloseAccount,
	},
	SagaTypeUpdateUserPassword: {
		EventTypeUpdateUserPassword,
	},
	SagaTypeGetUserDataByLogin: {
		EventTYpeGetUserDataByLogin,
	},
	SagaTypeCheckUserPassword: {
		EventTypeCheckUserPassword,
	},
}

// Список операций, входящих в SAG-у
var StartEventsListForSagaType = map[string][]string{
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
	SagaTypeAddAccountCache: {
		EventTypeAddAccountCache,
	},
	SagaTypeGetAccountData: {
		EventTypeGetAccountData,
	},
	SagaTypeWidthAccountCache: {
		EventTypeWidthAccountCache,
	},
	SagaTypeCloseAccount: {
		EventTypeCloseAccount,
	},
	SagaTypeUpdateUserPassword: {
		EventTypeUpdateUserPassword,
	},
	SagaTypeGetUserDataByLogin: {
		EventTYpeGetUserDataByLogin,
	},
	SagaTypeCheckUserPassword: {
		EventTypeCheckUserPassword,
	},
}

func CheckExistingEventTypeIntoSagaType(saga_type string, event_type string) (result bool) {

	result = false

	if ValidateSagaType(saga_type) &&
		ValidateEventType(event_type) {

		saga_events, ok := PossibleEventsListForSagaType[saga_type]
		if ok {

			for _, existing_event_type := range saga_events {

				if existing_event_type == event_type {
					result = true
					break
				}

			}

		}

	}

	return result

}

const (
	SagaGroupCreateUser         uint8 = 1
	SagaGroupCreateAccount      uint8 = 2
	SagaGroupAddAccountCache    uint8 = 3
	SagaGroupWidthAccountCache  uint8 = 4
	SagaGroupCloseAccount       uint8 = 5
	SagaGroupGetUserData        uint8 = 6
	SagaGroupGetAccountData     uint8 = 7
	SagaGroupUpdateUserPassword uint8 = 8
	SagaGroupGetUserDataByLogin uint8 = 9
	SagaGroupCheckUserPassword  uint8 = 10
)

var PossibleSagaGroups = []uint8{
	SagaGroupCreateUser,
	SagaGroupCreateAccount,
	SagaGroupAddAccountCache,
	SagaGroupWidthAccountCache,
	SagaGroupCloseAccount,
	SagaGroupGetUserData,
	SagaGroupGetAccountData,
	SagaGroupUpdateUserPassword,
	SagaGroupGetUserDataByLogin,
	SagaGroupCheckUserPassword,
}

func ValidateSagaGroup(saga_group uint8) bool {
	for i := 0; i < len(PossibleSagaGroups); i++ {
		if saga_group == PossibleSagaGroups[i] {
			return true
		}
	}

	return false
}

// Список зависимостей SAG в транзакции
var ListOfSagaDepend = map[uint8]map[string]models.SagaDepend{
	SagaGroupCreateUser: map[string]models.SagaDepend{
		SagaTypeCreateUser: {
			Parents:  nil,
			Children: nil,
		},
	},
	SagaGroupCreateAccount: map[string]models.SagaDepend{
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
	},
	SagaGroupAddAccountCache: map[string]models.SagaDepend{
		SagaTypeCheckUser: models.SagaDepend{
			Parents: nil,
			Children: []string{
				SagaTypeGetAccountData,
			},
		},
		SagaTypeGetAccountData: models.SagaDepend{
			Parents: []string{
				SagaTypeCheckUser,
			},
			Children: []string{
				SagaTypeAddAccountCache,
			},
		},
		SagaTypeAddAccountCache: models.SagaDepend{
			Parents: []string{
				SagaTypeGetAccountData,
			},
			Children: nil,
		},
	},
	SagaGroupWidthAccountCache: map[string]models.SagaDepend{
		SagaTypeCheckUser: models.SagaDepend{
			Parents: nil,
			Children: []string{
				SagaTypeGetAccountData,
			},
		},
		SagaTypeGetAccountData: models.SagaDepend{
			Parents: []string{
				SagaTypeCheckUser,
			},
			Children: []string{
				SagaTypeWidthAccountCache,
			},
		},
		SagaTypeWidthAccountCache: models.SagaDepend{
			Parents: []string{
				SagaTypeGetAccountData,
			},
			Children: nil,
		},
	},
	SagaGroupCloseAccount: map[string]models.SagaDepend{
		SagaTypeCheckUser: models.SagaDepend{
			Parents: nil,
			Children: []string{
				SagaTypeGetAccountData,
			},
		},
		SagaTypeGetAccountData: models.SagaDepend{
			Parents: []string{
				SagaTypeCheckUser,
			},
			Children: []string{
				SagaTypeCloseAccount,
			},
		},
		SagaTypeCloseAccount: models.SagaDepend{
			Parents: []string{
				SagaTypeGetAccountData,
			},
			Children: nil,
		},
	},
	SagaGroupGetUserData: map[string]models.SagaDepend{
		SagaTypeCheckUser: models.SagaDepend{
			Parents:  nil,
			Children: nil,
		},
	},
	SagaGroupGetAccountData: map[string]models.SagaDepend{
		SagaTypeCheckUser: models.SagaDepend{
			Parents: nil,
			Children: []string{
				SagaTypeGetAccountData,
			},
		},
		SagaTypeGetAccountData: models.SagaDepend{
			Parents: []string{
				SagaTypeCheckUser,
			},
			Children: nil,
		},
	},
	SagaGroupUpdateUserPassword: map[string]models.SagaDepend{
		SagaTypeUpdateUserPassword: {
			Parents:  nil,
			Children: nil,
		},
	},
	SagaGroupGetUserDataByLogin: map[string]models.SagaDepend{
		SagaTypeGetUserDataByLogin: models.SagaDepend{
			Parents:  nil,
			Children: nil,
		},
	},
	SagaGroupCheckUserPassword: map[string]models.SagaDepend{
		SagaTypeCheckUserPassword: models.SagaDepend{
			Parents:  nil,
			Children: nil,
		},
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

// Список данных, которые извлекаются, как результат операции
var SagaGroupResultDataUpdate = map[uint8]map[string][]string{
	SagaGroupCreateUser: map[string][]string{
		SagaTypeCreateUser: []string{
			"user_uuid",
		},
	},
	SagaGroupCreateAccount: map[string][]string{
		SagaTypeCheckUser: []string{
			"accounts",
		},
		SagaTypeReserveAccount: []string{
			"acc_id",
		},
	},
	SagaGroupAddAccountCache: map[string][]string{
		SagaTypeCheckUser: []string{
			"accounts",
		},
		SagaTypeGetAccountData: []string{
			"acc_status",
		},
	},
	SagaGroupWidthAccountCache: map[string][]string{
		SagaTypeCheckUser: []string{
			"accounts",
		},
		SagaTypeGetAccountData: []string{
			"acc_status",
		},
	},
	SagaGroupCloseAccount: map[string][]string{
		SagaTypeCheckUser: []string{
			"accounts",
		},
		SagaTypeGetAccountData: []string{
			"acc_status",
			"acc_cache",
		},
	},
	SagaGroupGetUserData: map[string][]string{
		SagaTypeCheckUser: []string{
			"inn",
			"accounts",
			"passport_series",
			"passport_number",
			"passport_first_name",
			"passport_first_surname",
			"passport_first_patronimic",
			"passport_birth_date",
			"passport_birth_location",
			"passport_pick_up_point",
			"passport_authority",
			"passport_authority_date",
			"passport_registration_address",
			"user_id",
			"user_login",
		},
	},
	SagaGroupGetAccountData: map[string][]string{
		SagaTypeCheckUser: []string{
			"accounts",
		},
		SagaTypeGetAccountData: []string{
			"acc_name",
			"acc_status",
			"acc_cache",
			"acc_cache_value",
			"acc_culc_number",
			"acc_corr_number",
			"acc_bic",
			"acc_cio",
			"acc_reserve_reason",
		},
	},
	SagaGroupGetUserDataByLogin: map[string][]string{
		SagaTypeGetUserDataByLogin: {
			"user_inn",
			"user_id",
			"user_login",
			"accounts",
		},
	},
}

// Возвращаемые данные при получении статуса операции
var SagaGroupDataIsResult = map[uint8]map[string]map[string][]string{
	SagaGroupCreateUser: map[string]map[string][]string{
		SagaTypeCreateUser: map[string][]string{
			EventTypeCreateUser: []string{
				"user_uuid",
			},
		},
	},
	SagaGroupCreateAccount: map[string]map[string][]string{
		SagaTypeReserveAccount: map[string][]string{
			EventTypeReserveAccount: []string{
				"acc_id",
			},
		},
	},
	SagaGroupAddAccountCache:   nil,
	SagaGroupWidthAccountCache: nil,
	SagaGroupCloseAccount:      nil,
	SagaGroupGetUserData: map[string]map[string][]string{
		SagaTypeCheckUser: map[string][]string{
			EventTypeGetUserData: []string{
				"inn",
				"accounts",
				"passport_series",
				"passport_number",
				"passport_first_name",
				"passport_first_surname",
				"passport_first_patronimic",
				"passport_birth_date",
				"passport_birth_location",
				"passport_pick_up_point",
				"passport_authority",
				"passport_authority_date",
				"passport_registration_address",
				"user_id",
				"user_login",
			},
		},
	},
	SagaGroupGetAccountData: map[string]map[string][]string{
		SagaTypeGetAccountData: map[string][]string{
			EventTypeGetAccountData: []string{
				"acc_name",
				"acc_status",
				"acc_cache",
				"acc_cache_value",
				"acc_culc_number",
				"acc_corr_number",
				"acc_bic",
				"acc_cio",
				"acc_reserve_reason",
			},
		},
	},
	SagaGroupGetUserDataByLogin: map[string]map[string][]string{
		SagaTypeGetUserDataByLogin: map[string][]string{
			EventTYpeGetUserDataByLogin: {
				"user_inn",
				"user_id",
				"user_login",
				"accounts",
			},
		},
	},
	SagaGroupCheckUserPassword: nil,
}
