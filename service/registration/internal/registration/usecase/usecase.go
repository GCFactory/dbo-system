package usecase

import (
	"encoding/json"
	"github.com/GCFactory/dbo-system/platform/config"
	"github.com/GCFactory/dbo-system/platform/pkg/logger"
	"github.com/GCFactory/dbo-system/service/registration/internal/models"
	"github.com/GCFactory/dbo-system/service/registration/internal/registration"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"golang.org/x/net/context"
)

type registrationUC struct {
	cfg              *config.Config
	registrationRepo registration.Repository
	logger           logger.Logger
}

func (regUC *registrationUC) CreateSaga(ctx context.Context, saga models.Saga) error {
	// TODO: реализовать
	return nil
}

func (regUC *registrationUC) DeleteSaga(ctx context.Context, saga_uuid uuid.UUID) error {
	//	TODO: реализовать
	return nil
}

func (regUC *registrationUC) AddSagaEvent(ctx context.Context, saga_uuid uuid.UUID, event_name string) error {
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationUC.AddSagaEvent")
	defer span.Finish()

	//	Проверка на пустоту строки
	if event_name == "" {
		return ErrorEmptySagaName
	}

	//	Проверка на то, что такое название event-а вообще возможно
	err := regUC.ValidateEventName(ctxWithTrace, event_name)
	if err != nil {
		return err
	}

	//	Проверка на существование saga с таким id
	saga, err := regUC.registrationRepo.GetSagaById(ctxWithTrace, saga_uuid)
	if err != nil {
		return ErrorNoSaga
	}

	err = regUC.ValidateEventNameForSagaType(ctxWithTrace, saga.Saga_type, event_name)
	if err != nil {
		return err
	}

	list_of_events_str := saga.Saga_list_of_events

	var list_of_events models.SagaListEvents

	err = json.Unmarshal([]byte(list_of_events_str), &list_of_events)
	if err != nil {
		return ErrorUnmarshal
	}

	//	Проверка на то, что в saga есть такой же event или противоположный ему
	for i := 0; i < len(list_of_events.EventList); i++ {
		inv_event := ""
		if event_name[0] == '-' {
			inv_event = inv_event[1:]
		} else {
			inv_event = "-" + event_name
		}
		//	Такой же
		if list_of_events.EventList[i] == event_name {
			return ErrorEventExist
		} else if list_of_events.EventList[i] == inv_event { //	Противоположный
			return ErrorInverseEventExist
		}
	}

	list_of_events.EventList = append(list_of_events.EventList, event_name)

	result_list_of_events_str, err := json.Marshal(list_of_events)
	if err != nil {
		return ErrorMarshal
	}

	err = regUC.registrationRepo.UpdateSagaEvents(ctxWithTrace, saga_uuid, string(result_list_of_events_str))
	if err != nil {
		return ErrorUpdateSagaEvents
	}

	return nil
}

func (regUC *registrationUC) RemoveSagaEvent(ctx context.Context, saga_uuid uuid.UUID, event_name string) error {
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationUC.RemoveSagaEvent")
	defer span.Finish()

	//	Проверка на пустоту строки
	if event_name == "" {
		return ErrorEmptySagaName
	}

	//	Проверка на то, что такое название event-а вообще возможно
	err := regUC.ValidateEventName(ctxWithTrace, event_name)
	if err != nil {
		return err
	}

	//	Проверка на существование saga с таким id
	saga, err := regUC.registrationRepo.GetSagaById(ctxWithTrace, saga_uuid)
	if err != nil {
		return ErrorNoSaga
	}

	err = regUC.ValidateEventNameForSagaType(ctxWithTrace, saga.Saga_type, event_name)
	if err != nil {
		return err
	}

	list_of_events_str := saga.Saga_list_of_events

	var list_of_events models.SagaListEvents

	err = json.Unmarshal([]byte(list_of_events_str), &list_of_events)
	if err != nil {
		return ErrorUnmarshal
	}

	var flag_exist = false
	var position int

	//	Проверка на то, что в saga есть такой же event
	for i := 0; i < len(list_of_events.EventList); i++ {
		if list_of_events.EventList[i] == event_name {
			flag_exist = true
			position = i
			break
		}
	}

	if flag_exist == false {
		return ErrorNoEvent
	}

	list_of_events.EventList[position] = list_of_events.EventList[len(list_of_events.EventList)-1]
	list_of_events.EventList[len(list_of_events.EventList)-1] = ""
	list_of_events.EventList = list_of_events.EventList[:len(list_of_events.EventList)-1]

	result_list_of_events_str, err := json.Marshal(list_of_events)
	if err != nil {
		return ErrorMarshal
	}

	err = regUC.registrationRepo.UpdateSagaEvents(ctxWithTrace, saga_uuid, string(result_list_of_events_str))
	if err != nil {
		return ErrorUpdateSagaEvents
	}

	return nil
}

func (regUC *registrationUC) FallBackSaga(ctx context.Context, saga_uuid uuid.UUID) error {
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationUC.FallBackSaga")
	defer span.Finish()

	saga, err := regUC.registrationRepo.GetSagaById(ctxWithTrace, saga_uuid)
	if err != nil {
		return ErrorNoSaga
	}

	saga = saga
	//	TODO: доделать

	return nil
}

func (regUC *registrationUC) ValidateEventName(ctx context.Context, event_name string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "registrationUC.ValidateEventName")
	defer span.Finish()

	var flag_exist = false

	for i := 0; i < len(PossibleEventsList); i++ {
		inv_possiple_event := "-" + PossibleEventsList[i]
		if PossibleEventsList[i] == event_name || inv_possiple_event == event_name {
			flag_exist = true
			break
		}
	}

	if flag_exist == false {
		return ErrorWrongEventName
	}

	return nil
}

func (regUC *registrationUC) ValidateEventNameForSagaType(ctx context.Context, saga_type uint, event_name string) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "registrationUC.ValidateEventNameForSagaType")
	defer span.Finish()

	if _, ok := PossibleEventsListForSagaType[saga_type]; !ok {
		return ErrorWrongSagaType
	}

	//	Проверка на то, возможно ли такое название event-а для данного типа sag-и
	var flag_exist = false

	for i := 0; i < len(PossibleEventsListForSagaType[saga_type]); i++ {
		tmp_str := "-" + PossibleEventsListForSagaType[saga_type][i]
		if event_name == PossibleEventsListForSagaType[saga_type][i] || event_name == tmp_str {
			flag_exist = true
			break
		}
	}

	if flag_exist == false {
		return ErrorWrongEventNameForSagaType
	}

	return nil
}

func NewRegistrationUseCase(cfg *config.Config, registration_repo registration.Repository, log logger.Logger) registration.UseCase {
	return &registrationUC{cfg: cfg, registrationRepo: registration_repo, logger: log}
}
