package usecase

const (
	EventTypeCreateUser       string = "create_user"
	EventTypeReserveAccount   string = "reserve_account"
	EventTypeCreateAccount    string = "create_account"
	EventTypeOpenAccount      string = "open_account"
	EventTypeAddAccountToUser string = "add_account_to_user"
)

var PossibleEventsList = [...]string{
	EventTypeCreateUser,
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

const (
// Here will be reverts events
)

var RevertEvent = map[string]string{}

//type EventRealisation interface {
//	CreateEvent(ctx context.Context) (event *models.Event, err error)
//	StartEvent(ctx context.Context, event *models.Event) error
//	RevertEvent(ctx context.Context, event *models.Event) error
//}
