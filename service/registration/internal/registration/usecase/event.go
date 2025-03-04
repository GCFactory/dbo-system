package usecase

const (
	EventTypeCreateUser         string = "add_user"
	EventTypeGetUserData        string = "get_user_data"
	EventTypeReserveAccount     string = "reserve_acc"
	EventTypeCreateAccount      string = "create_acc"
	EventTypeOpenAccount        string = "open_acc"
	EventTypeAddAccountToUser   string = "add_user_account"
	EventTypeAddAccountCache    string = "adding_acc"
	EventTypeWidthAccountCache  string = "width_acc"
	EventTypeGetAccountData     string = "get_acc_data"
	EventTypeCloseAccount       string = "close_acc"
	EventTypeRemoveAccount      string = "remove_acc"
	EventTypeRemoveUserAccount  string = "remove_user_account"
	EventTypeUpdateUserPassword string = "update_user_password"
)

var PossibleEventsList = [...]string{
	EventTypeCreateUser,
	EventTypeGetUserData,
	EventTypeReserveAccount,
	EventTypeCreateAccount,
	EventTypeOpenAccount,
	EventTypeAddAccountToUser,
	EventTypeAddAccountCache,
	EventTypeWidthAccountCache,
	EventTypeGetAccountData,
	EventTypeCloseAccount,
	EventTypeRemoveAccount,
	EventTypeRemoveUserAccount,
	EventTypeUpdateUserPassword,
}

func ValidateEventType(eventType string) bool {

	for i := 0; i < len(PossibleEventsList); i++ {
		if eventType == PossibleEventsList[i] {
			return true
		}
	}

	return false

}

const (
	EventStatusUndefined         uint8 = 0
	EventStatusCreated           uint8 = 10
	EventStatusInProgress        uint8 = 20
	EventStatusCompleted         uint8 = 30
	EventStatusFallBackInProcess uint8 = 40
	EventStatusFallBackCompleted uint8 = 50
	EventStatusFallBackError     uint8 = 250
	EventStatusError             uint8 = 255
)

var PossibleEventStatus = [...]uint8{
	EventStatusUndefined,
	EventStatusCreated,
	EventStatusInProgress,
	EventStatusCompleted,
	EventStatusFallBackInProcess,
	EventStatusFallBackCompleted,
	EventStatusFallBackError,
	EventStatusError,
}

func ValidateEventStatus(status uint8) bool {
	for i := 0; i < len(PossibleEventStatus); i++ {
		if status == PossibleEventStatus[i] {
			return true
		}
	}

	return false
}

var RequiredEventListOfData = map[string][]string{
	EventTypeCreateUser: []string{
		"user_inn",
		"passport_number",
		"passport_series",
		"name",
		"surname",
		"patronimic",
		"birth_date",
		"birth_location",
		"pick_up_point",
		"authority",
		"authority_date",
		"registration_adress",
		"login",
		"password",
	},
	EventTypeGetUserData: []string{
		"user_id",
	},
	EventTypeReserveAccount: []string{
		"culc_number",
		"corr_number",
		"bic",
		"cio",
		"reserve_reason",
	},
	EventTypeCreateAccount: []string{
		"acc_id",
	},
	EventTypeOpenAccount: []string{
		"acc_id",
	},
	EventTypeAddAccountToUser: []string{
		"user_id",
		"acc_id",
	},
	EventTypeAddAccountCache: []string{
		"acc_id",
		"cache_diff",
	},
	EventTypeWidthAccountCache: []string{
		"acc_id",
		"cache_diff",
	},
	EventTypeGetAccountData: []string{
		"acc_id",
	},
	EventTypeCloseAccount: []string{
		"acc_id",
	},
	EventTypeRemoveUserAccount: []string{
		"user_id",
		"acc_id",
	},
	EventTypeRemoveAccount: []string{
		"acc_id",
	},
	EventTypeUpdateUserPassword: {
		"user_id",
		"new_password",
	},
}

// Обратные события
var RevertEvent = map[uint8]map[string]map[string]string{
	SagaGroupCreateAccount: map[string]map[string]string{
		SagaTypeOpenAccountAndAddToUser: map[string]string{
			EventTypeAddAccountToUser: EventTypeRemoveUserAccount,
		},
		SagaTypeReserveAccount: map[string]string{
			EventTypeReserveAccount: EventTypeRemoveAccount,
		},
	},
}
