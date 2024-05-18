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
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationUC.CreateSaga")
	defer span.Finish()

	_, err := regUC.registrationRepo.GetSagaById(ctxWithTrace, saga.Saga_uuid)
	if err == nil {
		return ErrorSagaExist
	}

	if _, ok := PossibleEventsListForSagaType[saga.Saga_type]; !ok {
		return ErrorWrongSagaType
	}

	saga.Saga_status = StatusCreated

	list_of_events := models.SagaListEvents{}

	for i := 0; i < len(PossibleEventsListForSagaType[saga.Saga_type]); i++ {
		list_of_events.EventList = append(list_of_events.EventList, PossibleEventsListForSagaType[saga.Saga_type][i])
	}

	list_of_events_byte, err := json.Marshal(list_of_events)
	if err != nil {
		return ErrorMarshal
	}

	saga.Saga_list_of_events = string(list_of_events_byte)
	saga.Saga_status = StatusInProcess

	err = regUC.registrationRepo.CreateSaga(ctxWithTrace, saga)
	if err != nil {
		return ErrorCreateSaga
	}

	return nil
}

func (regUC *registrationUC) DeleteSaga(ctx context.Context, saga_uuid uuid.UUID) error {
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationUC.DeleteSaga")
	defer span.Finish()

	saga, err := regUC.registrationRepo.GetSagaById(ctxWithTrace, saga_uuid)
	if err != nil {
		return ErrorNoSaga
	}

	if saga.Saga_status != StatusCompleted {
		return ErrorSagaWasNotCompeted
	}

	list_of_events := models.SagaListEvents{}
	err = json.Unmarshal([]byte(saga.Saga_list_of_events), list_of_events)
	if err != nil {
		return ErrorUnmarshal
	}

	if len(list_of_events.EventList) != 0 {
		return ErrorNotEmptyEventList
	}

	err = regUC.registrationRepo.DeleteSaga(ctxWithTrace, saga_uuid)
	if err != nil {
		return ErrorDeleteSaga
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

	//	Если saga откатывается, то удаление простого event-а порождает обратный ему
	if saga.Saga_status == StatusFallBack && event_name[0] != '-' {
		inverse_event_name := "-" + event_name
		list_of_events.EventList = append(list_of_events.EventList, inverse_event_name)
	}

	result_list_of_events_str, err := json.Marshal(list_of_events)
	if err != nil {
		return ErrorMarshal
	}

	err = regUC.registrationRepo.UpdateSagaEvents(ctxWithTrace, saga_uuid, string(result_list_of_events_str))
	if err != nil {
		return ErrorUpdateSagaEvents
	}

	if len(list_of_events.EventList) == 0 {
		if saga.Saga_type == StatusInProcess { //	Если saga была в процессе
			if err = regUC.registrationRepo.UpdateSagaStatus(ctxWithTrace, saga_uuid, StatusCompleted); err != nil {
				return ErrorUpdateSagaStatus
			}
		} else if saga.Saga_type == StatusFallBack { //	Если сага была в откатывании
			if err = regUC.registrationRepo.UpdateSagaStatus(ctxWithTrace, saga_uuid, StatusError); err != nil {
				return ErrorUpdateSagaStatus
			}
		}
	}

	return nil
}

func (regUC *registrationUC) FallBackSaga(ctx context.Context, saga_uuid uuid.UUID, event_name string) error {
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationUC.FallBackSaga")
	defer span.Finish()

	saga, err := regUC.registrationRepo.GetSagaById(ctxWithTrace, saga_uuid)
	if err != nil {
		return ErrorNoSaga
	}

	if saga.Saga_status == StatusFallBack {
		return ErrorSagaAlreadyFallBack
	} else if saga.Saga_status != StatusInProcess {
		return ErrorWrongFallBackStatus
	}

	list_of_events := models.SagaListEvents{}

	err = json.Unmarshal([]byte(saga.Saga_list_of_events), list_of_events)
	if err != nil {
		return ErrorUnmarshal
	}

	for i := 0; i < len(PossibleEventsListForSagaType[saga.Saga_type]); i++ {
		inv_event_name := "-" + PossibleEventsListForSagaType[saga.Saga_type][i]
		flag_exist := false
		for j := 0; j < len(list_of_events.EventList); j++ {
			if PossibleEventsListForSagaType[saga.Saga_type][i] == list_of_events.EventList[j] {
				flag_exist = true
				break
			} else if inv_event_name == list_of_events.EventList[j] {
				return ErrorSagaAlreadyFallBack
			}
		}
		//	Если нет такого event-а или обратного ему, то добавляем обратный ему
		if flag_exist == false {
			list_of_events.EventList = append(list_of_events.EventList, inv_event_name)
		}
	}

	list_of_events_bytes, err := json.Marshal(list_of_events)
	if err != nil {
		return ErrorMarshal
	}

	saga.Saga_list_of_events = string(list_of_events_bytes)
	saga.Saga_status = StatusFallBack

	saga, err = regUC.FallBackEvent(ctxWithTrace, saga, event_name)
	if err != nil {
		return ErrorFallBackEvent
	}

	err = regUC.registrationRepo.UpdateSagaEvents(ctxWithTrace, saga.Saga_uuid, saga.Saga_list_of_events)
	if err != nil {
		return ErrorUpdateSagaEvents
	}

	err = regUC.registrationRepo.UpdateSagaStatus(ctxWithTrace, saga.Saga_uuid, StatusFallBack)
	if err != nil {
		return ErrorUpdateSagaStatus
	}

	return nil
}

func (regUC *registrationUC) FallBackEvent(ctx context.Context, saga *models.Saga, event_name string) (*models.Saga, error) {
	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationUC.FallBackSaga")
	defer span.Finish()

	if saga.Saga_status != StatusFallBack {
		return nil, ErrorWrongFallBackStatusForFallBackEvent
	}

	err := regUC.ValidateEventName(ctxWithTrace, event_name)
	if err != nil {
		return nil, ErrorWrongEventName
	}

	err = regUC.ValidateEventNameForSagaType(ctxWithTrace, saga.Saga_type, event_name)
	if err != nil {
		return nil, ErrorWrongEventNameForSagaType
	}

	if event_name[0] == '-' {
		return nil, ErrorInversedEvent
	}

	list_of_events := models.SagaListEvents{}
	err = json.Unmarshal([]byte(saga.Saga_list_of_events), list_of_events)

	position := -1

	inverse_event_name := "-" + event_name
	for i := 0; i < len(list_of_events.EventList); i++ {
		if list_of_events.EventList[i] == inverse_event_name {
			return nil, ErrorEventExist
		} else if list_of_events.EventList[i] == event_name {
			position = i
		}
	}

	if position == -1 {
		return nil, ErrorNoEvent
	}

	list_of_events.EventList[position] = list_of_events.EventList[len(list_of_events.EventList)-1]
	list_of_events.EventList[len(list_of_events.EventList)-1] = ""
	list_of_events.EventList = list_of_events.EventList[:len(list_of_events.EventList)-1]
	list_of_events.EventList = append(list_of_events.EventList, inverse_event_name)

	list_of_events_bytes, err := json.Marshal(list_of_events)
	if err != nil {
		return nil, ErrorMarshal
	}

	saga.Saga_list_of_events = string(list_of_events_bytes)

	return saga, nil
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
