package usecase

import "errors"

var (
	ErrorInvalidEventType               = errors.New("Invalid event type")
	ErrorSagaWasNotFound                = errors.New("Saga was not found")
	ErrorEventWasNotFound               = errors.New("Event was not found")
	ErrorWrongEventStatusForRollBack    = errors.New("Wrong event status for rollback")
	ErrorInvalidEventUuid               = errors.New("Invalid event uuid")
	ErrorInvalidSagaType                = errors.New("Wrong saga type!")
	ErrorInvalidSagGroup                = errors.New("Invalid saga group!")
	ErrorWrongEventStatusForDelete      = errors.New("Wrong event status for delete")
	ErrorNoDependsData                  = errors.New("No depends data")
	ErrorWrongSagaStatusForRollBack     = errors.New("Wrong saga status for rollback")
	ErrorWrongOperation                 = errors.New("Wrong operation")
	ErrorSagaWasNotFoundInPreviousLayer = errors.New("Saga was not found in previous layer of tree")
	ErrorNoRootsSagas                   = errors.New("No roots sagas")
	ErrorWrongSagaType                  = errors.New("Wrong saga type")
	ErrorInvalidEventStatus             = errors.New("Invalid event status")
	ErrorInvalidSagaStatus              = errors.New("Invalid saga status")
	ErrorNoReventEvent                  = errors.New("No revent event")
	ErrorEventDataNotExist              = errors.New("Event data not exist")
)
