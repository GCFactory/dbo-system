package usecase

import "errors"

var (
	ErrorWrongEventName                      = errors.New("Invalid event name, no such name")
	ErrorWrongSagaType                       = errors.New("Invalid saga type, no such type")
	ErrorWrongEventNameForSagaType           = errors.New("Invalid event name for this saga type")
	ErrorEmptySagaName                       = errors.New("Empty saga name")
	ErrorNoSaga                              = errors.New("No such saga with this id")
	ErrorUnmarshal                           = errors.New("Unmarshal error")
	ErrorMarshal                             = errors.New("Marshal error")
	ErrorUpdateSagaEvents                    = errors.New("Update saga error")
	ErrorNoEvent                             = errors.New("No such event in saga")
	ErrorNotEmptyEventList                   = errors.New("Saga hasn't empty events list")
	ErrorDeleteSaga                          = errors.New("Delete saga error")
	ErrorSagaExist                           = errors.New("Saga with this id already exist")
	ErrorCreateSaga                          = errors.New("Create saga error")
	ErrorSagaAlreadyFallBack                 = errors.New("Saga already fall back")
	ErrorSagaWasNotCompeted                  = errors.New("Saga was not completed")
	ErrorUpdateSagaStatus                    = errors.New("Update saga error")
	ErrorWrongFallBackStatus                 = errors.New("Wrong saga status for fall back")
	ErrorWrongFallBackStatusForFallBackEvent = errors.New("Wrong saga status for fall back its event")
	ErrorEventExist                          = errors.New("Event already exist")
	ErrorInversedEvent                       = errors.New("Event is inversed")
	ErrorFallBackEvent                       = errors.New("Error fall back event")
)
