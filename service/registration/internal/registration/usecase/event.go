package usecase

const (
	AntiFrod      string = "AntiFrod"
	ReserveNumber string = "ReserveNumber"
	TestEvent     string = "TestEvent"
)

var PossibleEventsList = [...]string{
	AntiFrod,
	ReserveNumber,
	TestEvent,
}

var PossibleEventsListForSagaType = map[uint][]string{
	SagaRegistration: {
		AntiFrod,
		ReserveNumber,
	},
}
