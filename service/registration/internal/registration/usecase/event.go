package usecase

const (
	EventTypeCreateUser       string = "add_user"
	EventTypeGetUserData      string = "get_user_data"
	EventTypeReserveAccount   string = "reserve_acc"
	EventTypeCreateAccount    string = "create_account"
	EventTypeOpenAccount      string = "open_account"
	EventTypeAddAccountToUser string = "add_user_account"
)

var PossibleEventsList = [...]string{
	EventTypeCreateUser,
	EventTypeGetUserData,
	EventTypeReserveAccount,
	EventTypeCreateAccount,
	EventTypeOpenAccount,
	EventTypeAddAccountToUser,
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

var EventListOfData = map[string][]string{
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
}

const (
// Here will be reverts events
)

var RevertEvent = map[string]string{}

//type EventRealisation interface {
//	CreateEvent(ctx context.Context) (event *models.Event, err error)
//	StartEvent(ctx context.Context, event *models.Event) error
//	RevertEvent(ctx context.Context, event *models.Event) error
//}
