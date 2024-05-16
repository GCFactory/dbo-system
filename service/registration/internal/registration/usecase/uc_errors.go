package usecase

import "errors"

var (
	ErrorWrongEventName            = errors.New("registrationUC.ValidateEventName.WrongEventName")
	ErrorWrongSagaType             = errors.New("registrationUC.ValidateEventNameForSagaType.WrongSagaType")
	ErrorWrongEventNameForSagaType = errors.New("registrationUC.ValidateEventNameForSagaType.WrongEventNameForSagaType")
	ErrorEmptySagaName             = errors.New("registrationUC.AddSagaEvent.EmptySagaName")
	ErrorNoSaga                    = errors.New("registrationUC.AddSagaEvent.GetSagaById")
	ErrorUnmarshal                 = errors.New("registrationUC.AddSagaEvent.Unmarshal")
	ErrorEventExist                = errors.New("registrationUC.AddSagaEvent.EventExist")
	ErrorInverseEventExist         = errors.New("registrationUC.AddSagaEvent.InverseEventExist")
	ErrorMarshal                   = errors.New("registrationUC.AddSagaEvent.Marshal")
	ErrorUpdateSagaEvents          = errors.New("registrationUC.AddSagaEvent.UpdateSagaEvents")
)
