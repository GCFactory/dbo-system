package usecase

const (
	AntiFrod      string = "AntiFrod"
	ReserveNumber string = "ReserveNumber"
)

var PossibleEventsList = [...]string{
	AntiFrod,
	ReserveNumber,
}

var PossibleEventsListForSagaType = map[uint][...]string{
	SagaRegistration: {
		AntiFrod,
		ReserveNumber,
	},
}
