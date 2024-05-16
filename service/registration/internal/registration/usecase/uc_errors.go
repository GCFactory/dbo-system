package usecase

import "errors"

var (
	ErrorWrongEventName            = errors.New("Invalid event name, no such name")
	ErrorWrongSagaType             = errors.New("Invalid saga type, no such type")
	ErrorWrongEventNameForSagaType = errors.New("Invalid event name for this saga type")
	ErrorEmptySagaName             = errors.New("Empty saga name")
	ErrorNoSaga                    = errors.New("No such saga with this id")
	ErrorUnmarshal                 = errors.New("Unmarshal error")
	ErrorEventExist                = errors.New("Event already exist")
	ErrorInverseEventExist         = errors.New("Inverse event already exist")
	ErrorMarshal                   = errors.New("Marshal error")
	ErrorUpdateSagaEvents          = errors.New("Update saga error")
	ErrorNoEvent                   = errors.New("No such event in saga")
)
