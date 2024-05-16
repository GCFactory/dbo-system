package usecase

const (
	AntiFrod      = "AntiFrod"
	ReserveNumber = "ReserveNumber"
)

var PossibleEventsList = [...]string{
	AntiFrod,
	ReserveNumber,
}

var PossibleEventsListForSagaType = map[uint][...]string{
	0: {
		AntiFrod,
		ReserveNumber,
	},
}
