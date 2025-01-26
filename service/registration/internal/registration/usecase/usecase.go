package usecase

import (
	"context"
	"encoding/json"
	"github.com/GCFactory/dbo-system/platform/pkg/logger"
	"github.com/GCFactory/dbo-system/service/registration/config"
	"github.com/GCFactory/dbo-system/service/registration/internal/models"
	"github.com/GCFactory/dbo-system/service/registration/internal/registration"
	"github.com/GCFactory/dbo-system/service/registration/internal/registration/repository"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"slices"
)

type registrationUC struct {
	cfg              *config.Config
	registrationRepo registration.Repository
	logger           logger.Logger
}

func (regUC registrationUC) CreateEvent(ctx context.Context, event_type string, saga_uuid uuid.UUID, is_fall_back bool, revet_event_uuid uuid.UUID) (event *models.Event, err error) {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationUC.CreateEvent")
	defer span.Finish()

	if !ValidateEventType(event_type) {
		return nil, ErrorInvalidEventType
	} else if saga_uuid == uuid.Nil {
		return nil, ErrorSagaWasNotFound
	}

	if _, err := regUC.registrationRepo.GetSaga(ctxWithTrace, saga_uuid); err != nil {
		if err == repository.ErrorGetSaga {
			return nil, ErrorSagaWasNotFound
		} else {
			return nil, err
		}
	}

	var revert_event *models.Event

	if revet_event_uuid != uuid.Nil {
		revert_event, err = regUC.registrationRepo.GetEvent(ctxWithTrace, revet_event_uuid)
		if err != nil {
			if err == repository.ErrorGetEvent {
				return nil, ErrorEventWasNotFound
			}
		}

		if revert_event.Event_status != EventStatusCompleted &&
			revert_event.Event_status != EventStatusError {
			return nil, ErrorWrongEventStatusForRollBack
		}
	}

	event_data, ok := EventListOfData[event_type]
	if !ok {
		return nil, ErrorEventDataNotExist
	}

	local_event := &models.Event{
		Event_uuid:          uuid.New(),
		Saga_uuid:           saga_uuid,
		Event_is_roll_back:  is_fall_back,
		Event_rollback_uuid: revet_event_uuid,
		Event_required_data: event_data,
		Event_name:          event_type,
		Event_status:        EventStatusCreated,
		Event_result:        "{}",
	}

	err = regUC.registrationRepo.CreateEvent(ctxWithTrace, local_event)
	if err != nil {
		return nil, err
	}

	if revert_event != nil {
		revert_event.Event_rollback_uuid = local_event.Event_uuid
		err = regUC.registrationRepo.UpdateEvent(ctxWithTrace, revert_event)
		if err != nil {
			return nil, err
		}
	}

	return local_event, nil
}

func (regUC registrationUC) RevertEvent(ctx context.Context, event_uuid uuid.UUID) (*models.Event, error) {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationUC.RevertEvent")
	defer span.Finish()

	if event_uuid == uuid.Nil {
		return nil, ErrorInvalidEventType
	}

	event, err := regUC.registrationRepo.GetEvent(ctxWithTrace, event_uuid)
	if err != nil {
		if err == repository.ErrorGetEvent {
			return nil, ErrorEventWasNotFound
		} else {
			return nil, err
		}
	}

	if event.Event_status != EventStatusCreated &&
		event.Event_status != EventStatusCompleted {
		return nil, ErrorWrongEventStatusForRollBack
	}

	if event.Event_is_roll_back {
		return nil, ErrorWrongEventStatusForRollBack
	}

	event_type, ok := RevertEvent[event.Event_name]
	if ok {

		revert_event, err := regUC.CreateEvent(ctxWithTrace, event_type, event.Saga_uuid, true, event_uuid)
		if err != nil {
			return nil, err
		}
		event.Event_status = EventStatusFallBackInProcess
		err = regUC.registrationRepo.UpdateEvent(ctxWithTrace, event)
		if err != nil {
			return nil, err
		}
		return revert_event, nil
	}

	return nil, nil
}

func (regUC registrationUC) DeleteEvent(ctx context.Context, event_uuid uuid.UUID) error {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationUC.DeleteEvent")
	defer span.Finish()

	if event_uuid == uuid.Nil {
		return ErrorInvalidEventType
	}

	event, err := regUC.registrationRepo.GetEvent(ctxWithTrace, event_uuid)
	if err != nil {
		if err == repository.ErrorGetEvent {
			return ErrorEventWasNotFound
		} else {
			return err
		}
	}

	if event.Event_status != EventStatusCompleted &&
		event.Event_status != EventStatusError &&
		event.Event_status != EventStatusFallBackCompleted &&
		event.Event_status != EventStatusCreated {
		return ErrorWrongEventStatusForDelete
	}

	err = regUC.registrationRepo.DeleteEvent(ctxWithTrace, event_uuid)
	if err != nil {
		return err
	}

	return nil
}

func (regUC registrationUC) CreateSaga(ctx context.Context, saga_type string, saga_group uint8, saga_data map[string]interface{}) (saga *models.Saga, err error) {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationUC.CreateSaga")
	defer span.Finish()

	if !ValidateSagaType(saga_type) {
		return nil, ErrorInvalidSagaType
	} else if !ValidateSagaGroup(saga_group) {
		return nil, ErrorInvalidSagGroup
	}

	saga = &models.Saga{
		Saga_uuid:   uuid.New(),
		Saga_status: SagaStatusUndefined,
		Saga_type:   saga_group,
		Saga_name:   saga_type,
		Saga_data:   saga_data,
	}

	err = regUC.registrationRepo.CreateSaga(ctxWithTrace, saga)
	if err != nil {
		return nil, err
	}

	events, ok := PossibleEventsListForSagaType[saga_type]
	if !ok {
		err = regUC.DeleteSaga(ctxWithTrace, saga.Saga_uuid)
		if err != nil {
			return nil, err
		}
		return nil, ErrorInvalidSagaType
	}

	// Fills saga by events
	for _, event_type := range events {
		if !ValidateEventType(event_type) {
			err = regUC.DeleteSaga(ctxWithTrace, saga.Saga_uuid)
			if err != nil {
				return nil, err
			}
			return nil, ErrorInvalidEventType
		}
		_, local_err := regUC.CreateEvent(ctxWithTrace, event_type, saga.Saga_uuid, false, uuid.Nil)
		if local_err != nil {
			err = regUC.DeleteSaga(ctxWithTrace, saga.Saga_uuid)
			if err != nil {
				return nil, err
			}
			return nil, local_err
		}
	}

	//saga_depend, ok := ListOfSagaDepend[saga_type]
	//if !ok {
	//	err = regUC.DeleteSaga(ctxWithTrace, saga.Saga_uuid)
	//	if err != nil {
	//		return nil, err
	//	}
	//	return nil, ErrorNoDependsData
	//}

	//if saga_depend.Children != nil {
	//	children := saga_depend.Children
	//	for _, child_type := range children {
	//		child_saga, err := regUC.CreateSaga(ctxWithTrace, child_type, saga_group)
	//		if err != nil {
	//			err = regUC.DeleteSaga(ctxWithTrace, saga.Saga_uuid)
	//			if err != nil {
	//				return nil, err
	//			}
	//			return nil, err
	//		} else {
	//			saga_connection := &models.SagaConnection{
	//				Current_saga_uuid:     saga.Saga_uuid,
	//				Next_saga_uuid:        child_saga.Saga_uuid,
	//				Acc_connection_status: SagaConnectionStatusWaiting,
	//			}
	//			err = regUC.registrationRepo.CreateSagaConnection(ctxWithTrace, saga_connection)
	//			if err != nil {
	//				err = regUC.DeleteSaga(ctxWithTrace, saga.Saga_uuid)
	//				if err != nil {
	//					return nil, err
	//				}
	//				return nil, err
	//			}
	//		}
	//	}
	//}

	saga.Saga_status = SagaStatusCreated
	err = regUC.registrationRepo.UpdateSaga(ctxWithTrace, saga)
	if err != nil {
		err = regUC.DeleteSaga(ctxWithTrace, saga.Saga_uuid)
		if err != nil {
			return nil, err
		}
		return nil, err
	}

	return saga, nil
}

func (regUC registrationUC) RevertSaga(ctx context.Context, saga_uuid uuid.UUID) (result []*models.Event, err error) {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationUC.CreateSaga")
	defer span.Finish()

	saga, err := regUC.registrationRepo.GetSaga(ctxWithTrace, saga_uuid)
	if err != nil {
		return nil, err
	}

	if saga.Saga_status != SagaStatusInProcess &&
		saga.Saga_status != SagaStatusCompleted {
		return nil, ErrorWrongSagaStatusForRollBack
	}

	list_of_child_connections, err := regUC.registrationRepo.GetSagaConnectionsCurrentSaga(ctxWithTrace, saga_uuid)
	if err != nil {
		return nil, err
	}

	var child_revert bool = false

	for _, connection := range list_of_child_connections.List_of_connetcions {
		child_saga_uuid := connection.Next_saga_uuid
		connection_status := connection.Acc_connection_status
		if connection_status == SagaConnectionStatusSuccess {
			new_events, err := regUC.RevertSaga(ctxWithTrace, child_saga_uuid)
			if err != nil {
				return nil, err
			}
			result = append(result, new_events...)
		}
	}

	if child_revert {
		return result, nil
	}

	list_of_events_uuid, err := regUC.registrationRepo.GetListOfSagaEvents(ctxWithTrace, saga_uuid)
	if err != nil {
		return result, err
	}

	var list_of_events []*models.Event = nil
	for _, event_uuid := range list_of_events_uuid.EventList {
		event, err := regUC.registrationRepo.GetEvent(ctxWithTrace, event_uuid)
		if err != nil {
			return result, err
		}
		list_of_events = append(list_of_events, event)
	}

	if list_of_events != nil {
		for _, event := range list_of_events {
			if event.Event_status == EventStatusCompleted && !event.Event_is_roll_back {
				reverted_saga, err := regUC.RevertEvent(ctxWithTrace, event.Event_uuid)
				if err != nil {
					return result, err
				}
				result = append(result, reverted_saga)
			}
		}
	}

	saga.Saga_status = SagaStatusFallBackInProcess

	err = regUC.registrationRepo.UpdateSaga(ctxWithTrace, saga)
	if err != nil {
		return result, err
	}

	return result, nil

}

func (regUC registrationUC) DeleteSaga(ctx context.Context, saga_uuid uuid.UUID) (err error) {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationUC.DeleteSaga")
	defer span.Finish()

	_, err = regUC.registrationRepo.GetSaga(ctxWithTrace, saga_uuid)
	if err != nil {
		return ErrorSagaWasNotFound
	}

	list_of_connections, err := regUC.registrationRepo.GetSagaConnectionsCurrentSaga(ctx, saga_uuid)
	if err != nil {
		return err
	}

	for _, connection := range list_of_connections.List_of_connetcions {
		child_saga_uuid := connection.Next_saga_uuid
		err = regUC.DeleteSaga(ctx, child_saga_uuid)
		if err != nil {
			return err
		}
	}

	err = regUC.registrationRepo.DeleteSaga(ctxWithTrace, saga_uuid)
	if err != nil {
		return err
	}

	return nil

}

func (regUC registrationUC) ProcessingSagaAndEvents(ctx context.Context, saga_uuid uuid.UUID, event_uuid uuid.UUID, success bool, data map[string]interface{}) (result []*models.Event, err error) {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationUC.ProcessingSagaAndEvents")
	defer span.Finish()

	result = nil

	if event_uuid != uuid.Nil {

		event, err := regUC.registrationRepo.GetEvent(ctxWithTrace, event_uuid)
		if err != nil {
			return result, err
		}

		switch event.Event_status {
		case EventStatusUndefined:
			{
				return result, ErrorInvalidEventStatus
			}
		case EventStatusCreated:
			{
				if regUC.CheckEventDataIsReady(ctxWithTrace, event_uuid, event.Saga_uuid) {
					event.Event_status = EventStatusInProgress
					err = regUC.registrationRepo.UpdateEvent(ctxWithTrace, event)
					if err != nil {
						return result, err
					}
					result = append(result, event)
				} else {
					event.Event_status = EventStatusError
					event.Event_result = "{ \"Local error\": \"Event hasn't got all data from saga!\"}"
					err = regUC.registrationRepo.UpdateEvent(ctxWithTrace, event)
					if err != nil {
						return result, err
					}
					new_events, err := regUC.ProcessingSagaAndEvents(ctxWithTrace, event.Saga_uuid, uuid.Nil, false, nil)
					if err != nil {
						return result, err
					}
					result = append(result, new_events...)
				}
			}
		case EventStatusInProgress:
			{
				event_result, err := json.Marshal(data)
				if err != nil {
					return result, err
				}
				event.Event_result = string(event_result)

				if success {
					event.Event_status = EventStatusCompleted
					err = regUC.registrationRepo.UpdateEvent(ctxWithTrace, event)
					if err != nil {
						return result, err
					}
					if event.Event_is_roll_back {
						new_events, err := regUC.ProcessingSagaAndEvents(ctxWithTrace, uuid.Nil, event.Event_rollback_uuid, true, nil)
						if err != nil {
							return result, err
						}
						result = slices.Concat(result, new_events)
					} else {
						saga, local_err := regUC.registrationRepo.GetSaga(ctxWithTrace, saga_uuid)
						if local_err != nil {
							return result, local_err
						}
						for field_name, field_value := range data {
							saga.Saga_data[field_name] = field_value
						}
						local_err = regUC.registrationRepo.UpdateSaga(ctxWithTrace, saga)
						if local_err != nil {
							return result, local_err
						}
					}
				} else {
					event.Event_status = EventStatusError
					err = regUC.registrationRepo.UpdateEvent(ctxWithTrace, event)
					if err != nil {
						return result, err
					}
					if event.Event_is_roll_back {
						new_events, err := regUC.ProcessingSagaAndEvents(ctxWithTrace, uuid.Nil, event.Event_rollback_uuid, false, nil)
						if err != nil {
							return result, err
						}
						result = slices.Concat(result, new_events)
					}
				}
			}
		case EventStatusCompleted:
			{
				if !success {
					new_event, err := regUC.RevertEvent(ctxWithTrace, event.Event_uuid)
					if err != nil {
						return result, err
					}
					if new_event == nil {
						return result, ErrorNoReventEvent
					}
					result = append(result, new_event)
				} else {
					err = regUC.SetOrUpdateSagaData(ctxWithTrace, event.Saga_uuid, data)
					if err != nil {
						return result, err
					}
				}
			}
		case EventStatusFallBackInProcess:
			{
				if success {
					event.Event_status = EventStatusFallBackCompleted
					err = regUC.registrationRepo.UpdateEvent(ctxWithTrace, event)
					if err != nil {
						return result, err
					}
				} else {
					event.Event_status = EventStatusFallBackError
					err = regUC.registrationRepo.UpdateEvent(ctxWithTrace, event)
					if err != nil {
						return result, err
					}
				}
			}
		case EventStatusFallBackCompleted:
			{
				// NONE
			}
		case EventStatusFallBackError:
			{
				// NONE
			}
		case EventStatusError:
			{
				// NONE
			}
		default:
			{
				return result, ErrorInvalidEventStatus
			}
		}
	}

	if saga_uuid != uuid.Nil {

		saga, err := regUC.registrationRepo.GetSaga(ctxWithTrace, saga_uuid)
		if err != nil {
			return result, err
		}
		list_of_events, err := regUC.registrationRepo.GetListOfSagaEvents(ctxWithTrace, saga_uuid)
		if err != nil {
			return result, err
		}

		// List of Saga's events
		var events []*models.Event = nil

		if list_of_events != nil {
			for _, event_uuid := range list_of_events.EventList {
				event, err := regUC.registrationRepo.GetEvent(ctxWithTrace, event_uuid)
				if err != nil {
					return result, err
				}
				events = append(events, event)
			}
		}

		switch saga.Saga_status {
		case SagaStatusUndefined:
			{
				return result, ErrorInvalidSagaStatus
			}
		case SagaStatusCreated:
			{
				all_is_ok := true
				for _, event := range events {
					if event.Event_status == EventStatusError {
						all_is_ok = false
					}
					new_events, err := regUC.ProcessingSagaAndEvents(ctxWithTrace, uuid.Nil, event.Event_uuid, true, nil)
					if err != nil {
						return result, err
					}
					result = append(result, new_events...)
				}
				check_saga, err := regUC.registrationRepo.GetSaga(ctxWithTrace, saga.Saga_uuid)
				if err != nil {
					return result, err
				}
				if check_saga.Saga_status == saga.Saga_status {
					if all_is_ok {
						saga.Saga_status = SagaStatusInProcess
					} else {
						saga.Saga_status = SagaStatusError
					}
					err = regUC.registrationRepo.UpdateSaga(ctxWithTrace, saga)
					if err != nil {
						return result, err
					}
				}
			}
		case SagaStatusInProcess:
			{
				all_completed := true
				is_error := false
				for _, event := range events {
					if event.Event_status != EventStatusCompleted {
						all_completed = false
						if event.Event_status == EventStatusError {
							is_error = true
						}
					}
				}

				if is_error {
					saga.Saga_status = SagaStatusError
					err = regUC.registrationRepo.UpdateSaga(ctxWithTrace, saga)
					if err != nil {
						return result, err
					}
					new_events, err := regUC.ProcessingSagaAndEvents(ctxWithTrace, saga_uuid, uuid.Nil, true, nil)
					if err != nil {
						return result, err
					}
					result = append(result, new_events...)
				}

				if all_completed {
					saga.Saga_status = SagaStatusCompleted
					err = regUC.registrationRepo.UpdateSaga(ctxWithTrace, saga)
					if err != nil {
						return result, err
					}
					new_events, err := regUC.ProcessingSagaAndEvents(ctxWithTrace, saga_uuid, uuid.Nil, true, nil)
					if err != nil {
						return result, err
					}
					result = append(result, new_events...)
				}
			}
		case SagaStatusCompleted:
			{
				list_of_connections, err := regUC.registrationRepo.GetSagaConnectionsCurrentSaga(ctxWithTrace, saga_uuid)
				if err != nil {
					return result, err
				}

				if success {
					for _, connection := range list_of_connections.List_of_connetcions {
						connection.Acc_connection_status = SagaConnectionStatusSuccess
						err = regUC.registrationRepo.UpdateSagaConnection(ctxWithTrace, connection)
						if err != nil {
							return result, err
						}
						err = regUC.SetOrUpdateSagaData(ctxWithTrace, connection.Next_saga_uuid, data)
						if err != nil {
							return result, err
						}

						new_events, err := regUC.ProcessingSagaAndEvents(ctxWithTrace, connection.Next_saga_uuid, uuid.Nil, true, nil)
						if err != nil {
							return result, err
						}
						result = append(result, new_events...)
					}
				} else {
					all_child_reverts := true
					for _, connection := range list_of_connections.List_of_connetcions {
						if connection.Acc_connection_status != SagaConnectionStatusFallBack {
							all_child_reverts = false
							if connection.Acc_connection_status == SagaConnectionStatusSuccess {
								new_events, err := regUC.ProcessingSagaAndEvents(ctxWithTrace, connection.Next_saga_uuid, uuid.Nil, false, nil)
								if err != nil {
									return result, err
								}
								result = append(result, new_events...)
							}
						}
					}

					if all_child_reverts {
						new_events, err := regUC.RevertSaga(ctxWithTrace, saga_uuid)
						if err != nil {
							return result, err
						}
						result = append(result, new_events...)
					}
				}
			}
		case SagaStatusFallBackInProcess:
			{
				all_is_completed := true
				is_error := false
				for _, event := range events {
					if event.Event_is_roll_back && event.Event_status != EventStatusCompleted ||
						!event.Event_is_roll_back && event.Event_status != EventStatusFallBackCompleted {
						all_is_completed = false
						if event.Event_status == EventStatusError {
							is_error = true
							break
						} else if !event.Event_is_roll_back && event.Event_status == EventStatusCompleted {
							reverted_event, err := regUC.RevertEvent(ctxWithTrace, event.Event_uuid)
							if err != nil {
								return result, err
							}
							result = append(result, reverted_event)
						}
					}
				}

				if is_error {
					saga.Saga_status = SagaStatusFallBackError
					err = regUC.registrationRepo.UpdateSaga(ctxWithTrace, saga)
					if err != nil {
						return result, err
					}
					new_events, err := regUC.ProcessingSagaAndEvents(ctxWithTrace, saga_uuid, uuid.Nil, true, nil)
					if err != nil {
						return result, err
					}
					result = append(result, new_events...)
				}

				if all_is_completed {
					saga.Saga_status = SagaStatusFallBackSuccess
					err = regUC.registrationRepo.UpdateSaga(ctxWithTrace, saga)
					if err != nil {
						return result, err
					}

					new_events, err := regUC.ProcessingSagaAndEvents(ctxWithTrace, saga_uuid, uuid.Nil, true, nil)
					if err != nil {
						return result, err
					}
					result = append(result, new_events...)
				}
			}
		case SagaStatusFallBackSuccess:
			{
				list_of_connections, err := regUC.registrationRepo.GetSagaConnectionsNextSaga(ctxWithTrace, saga_uuid)
				if err != nil {
					return result, err
				}
				for _, connection := range list_of_connections.List_of_connetcions {
					connection.Acc_connection_status = SagaConnectionStatusFallBack
					err = regUC.registrationRepo.UpdateSagaConnection(ctxWithTrace, connection)
					if err != nil {
						return result, err
					}
					new_events, err := regUC.ProcessingSagaAndEvents(ctxWithTrace, saga_uuid, uuid.Nil, false, nil)
					if err != nil {
						return result, err
					}
					result = append(result, new_events...)
				}

			}
		case SagaStatusFallBackError:
			{
				all_is_completed := true
				for _, event := range events {
					if !event.Event_is_roll_back && event.Event_status != EventStatusFallBackCompleted && event.Event_status != EventStatusFallBackError ||
						event.Event_is_roll_back && event.Event_status != EventStatusCompleted && event.Event_status != EventStatusError {
						all_is_completed = false
						if !event.Event_is_roll_back && event.Event_status == EventStatusCompleted {
							reverted_event, err := regUC.RevertEvent(ctxWithTrace, event.Event_uuid)
							if err != nil {
								return result, err
							}
							result = append(result, reverted_event)
						}
					}
				}

				if all_is_completed {
					list_of_connetions, err := regUC.registrationRepo.GetSagaConnectionsNextSaga(ctxWithTrace, saga_uuid)
					if err != nil {
						return result, err
					}
					for _, connetion := range list_of_connetions.List_of_connetcions {
						connetion.Acc_connection_status = SagaConnectionStatusFallBack
						err = regUC.registrationRepo.UpdateSagaConnection(ctxWithTrace, connetion)
						if err != nil {
							return result, err
						}
						new_events, err := regUC.ProcessingSagaAndEvents(ctxWithTrace, saga_uuid, uuid.Nil, false, nil)
						if err != nil {
							return result, err
						}
						result = append(result, new_events...)
					}
				}
			}
		case SagaStatusError:
			{
				all_completed := true
				for _, event := range events {
					if event.Event_is_roll_back && event.Event_status != EventStatusCompleted && event.Event_status != EventStatusError ||
						!event.Event_is_roll_back && event.Event_status != EventStatusFallBackError && event.Event_status != EventStatusFallBackCompleted {
						all_completed = false
						if !event.Event_is_roll_back && event.Event_status == EventStatusCompleted {
							reverted_event, err := regUC.RevertEvent(ctxWithTrace, event.Event_uuid)
							if err != nil {
								return result, err
							}
							if reverted_event == nil {
								return result, ErrorNoReventEvent
							}
							result = append(result, reverted_event)
							new_events, err := regUC.ProcessingSagaAndEvents(ctxWithTrace, uuid.Nil, reverted_event.Event_uuid, true, nil)
							if err != nil {
								return result, err
							}
							result = append(result, new_events...)
						}
					}
				}

				if all_completed {
					list_of_connections, err := regUC.registrationRepo.GetSagaConnectionsNextSaga(ctxWithTrace, saga_uuid)
					if err != nil {
						return result, err
					}
					for _, connection := range list_of_connections.List_of_connetcions {
						connection.Acc_connection_status = SagaConnectionStatusFailed
						err = regUC.registrationRepo.UpdateSagaConnection(ctxWithTrace, connection)
						if err != nil {
							return result, err
						}
						new_events, err := regUC.ProcessingSagaAndEvents(ctxWithTrace, connection.Current_saga_uuid, uuid.Nil, false, nil)
						if err != nil {
							return result, err
						}
						result = append(result, new_events...)
					}
				}

			}
		default:
			{
				return result, ErrorInvalidSagaStatus
			}
		}

	}

	return result, nil

}

func (regUC registrationUC) GetEventData(ctx context.Context, event_uuid uuid.UUID) (result map[string]interface{}, err error) {
	result = make(map[string]interface{})
	err = nil

	event, err := regUC.registrationRepo.GetEvent(ctx, event_uuid)
	if err != nil {
		return result, err
	}
	saga, err := regUC.registrationRepo.GetSaga(ctx, event.Saga_uuid)
	if err != nil {
		return result, err
	}

	saga_data := saga.Saga_data
	if regUC.CheckEventDataIsReady(ctx, saga.Saga_uuid, event_uuid) {
		for _, field := range event.Event_required_data {
			result[field] = saga_data[field]
		}
	} else {
		return result, ErrorEventDataNotExist
	}

	return result, err
}

func (regUC registrationUC) SetOrUpdateSagaData(ctx context.Context, saga_uuid uuid.UUID, new_saga_data map[string]interface{}) (err error) {
	err = nil

	if new_saga_data == nil {
		return err
	}

	saga, err := regUC.registrationRepo.GetSaga(ctx, saga_uuid)
	if err != nil {
		return err
	}

	saga_data := saga.Saga_data

	for key, value := range new_saga_data {
		saga_data[key] = value
	}

	saga.Saga_data = saga_data
	err = regUC.registrationRepo.UpdateSaga(ctx, saga)
	if err != nil {
		return err
	}

	return err
}

func (regUC registrationUC) CreateSagaTree(ctx context.Context, list_root_saga_types []string, saga_group uint8, start_data map[string]interface{}) (list_root_saga []*models.Saga, err error) {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationUC.CreateSagaTree")
	defer span.Finish()

	list_root_saga = nil

	// create root layer of saga tree
	for _, saga_type := range list_root_saga_types {
		saga, err := regUC.CreateSaga(ctxWithTrace, saga_type, saga_group, start_data)
		if err != nil {
			if list_root_saga != nil {
				for _, saga := range list_root_saga {
					err = regUC.DeleteSaga(ctxWithTrace, saga.Saga_uuid)
					if err != nil {
						return nil, err
					}
				}
			}
			return nil, err
		}
		list_root_saga = append(list_root_saga, saga)
	}

	var saga_tree map[uint8][]*models.Saga = make(map[uint8][]*models.Saga)

	layer_number := uint8(0)
	saga_tree[layer_number] = list_root_saga

	layer_number++

	var previous_saga_tree_layer map[string]*models.Saga = make(map[string]*models.Saga)
	var current_saga_tree_layer map[string]*models.Saga = make(map[string]*models.Saga)

	// fill previous_saga_tree_layer_types
	for _, saga := range list_root_saga {
		previous_saga_tree_layer[saga.Saga_name] = saga
	}

	for {
		for saga_type, _ := range previous_saga_tree_layer {

			saga_depends, is_exist := ListOfSagaDepend[saga_type]
			if !is_exist {
				saga_tree = regUC.FillNewSagaTreeLayer(ctxWithTrace, saga_tree, layer_number, current_saga_tree_layer)
				err = regUC.ClearSagaTree(ctxWithTrace, saga_tree)
				if err != nil {
					return nil, err
				}
				return nil, ErrorNoDependsData
			}
			for _, child_type := range saga_depends.Children {
				if !ValidateSagaType(child_type) {
					saga_tree = regUC.FillNewSagaTreeLayer(ctxWithTrace, saga_tree, layer_number, current_saga_tree_layer)
					err = regUC.ClearSagaTree(ctxWithTrace, saga_tree)
					if err != nil {
						return nil, err
					}
					return nil, ErrorInvalidSagaType
				}
				_, is_exist = current_saga_tree_layer[child_type]
				if !is_exist {
					saga_child, local_err := regUC.CreateSaga(ctxWithTrace, child_type, saga_group, nil)
					if local_err != nil {
						saga_tree = regUC.FillNewSagaTreeLayer(ctxWithTrace, saga_tree, layer_number, current_saga_tree_layer)
						err = regUC.ClearSagaTree(ctxWithTrace, saga_tree)
						if err != nil {
							return nil, err
						}
						return nil, local_err
					}
					current_saga_tree_layer[child_type] = saga_child

					child_depends, is_exist := ListOfSagaDepend[child_type]
					if !is_exist {
						saga_tree = regUC.FillNewSagaTreeLayer(ctxWithTrace, saga_tree, layer_number, current_saga_tree_layer)
						err = regUC.ClearSagaTree(ctxWithTrace, saga_tree)
						if err != nil {
							return nil, err
						}
						return nil, ErrorNoDependsData
					}

					for _, parent_type := range child_depends.Parents {
						if !ValidateSagaType(parent_type) {
							saga_tree = regUC.FillNewSagaTreeLayer(ctxWithTrace, saga_tree, layer_number, current_saga_tree_layer)
							err = regUC.ClearSagaTree(ctxWithTrace, saga_tree)
							if err != nil {
								return nil, err
							}
							return nil, ErrorInvalidSagaType
						}
						parent_saga, is_exist := previous_saga_tree_layer[parent_type]
						if !is_exist {
							saga_tree = regUC.FillNewSagaTreeLayer(ctxWithTrace, saga_tree, layer_number, current_saga_tree_layer)
							err = regUC.ClearSagaTree(ctxWithTrace, saga_tree)
							if err != nil {
								return nil, err
							}
							return nil, ErrorSagaWasNotFoundInPreviousLayer
						}
						saga_connection := &models.SagaConnection{
							Next_saga_uuid:        saga_child.Saga_uuid,
							Current_saga_uuid:     parent_saga.Saga_uuid,
							Acc_connection_status: SagaConnectionStatusWaiting,
						}

						err = regUC.registrationRepo.CreateSagaConnection(ctxWithTrace, saga_connection)
						if err != nil {
							err = regUC.ClearSagaTree(ctxWithTrace, saga_tree)
							if err != nil {
								return nil, err
							}
							return nil, ErrorInvalidSagaType
						}
					}
				}
			}
			saga_tree = regUC.FillNewSagaTreeLayer(ctxWithTrace, saga_tree, layer_number, current_saga_tree_layer)
		}

		if current_saga_tree_layer == nil {
			break
		}

		saga_tree = regUC.FillNewSagaTreeLayer(ctxWithTrace, saga_tree, layer_number, current_saga_tree_layer)
		layer_number++

		previous_saga_tree_layer = current_saga_tree_layer
		current_saga_tree_layer = nil

	}

	return list_root_saga, nil

}

func (regUC registrationUC) ClearSagaTree(ctx context.Context, saga_tree map[uint8][]*models.Saga) (err error) {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationUC.ClearSagaTree")
	defer span.Finish()

	for _, list_of_saga := range saga_tree {
		for _, saga := range list_of_saga {
			err = regUC.DeleteSaga(ctxWithTrace, saga.Saga_uuid)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (regUC registrationUC) FillNewSagaTreeLayer(ctx context.Context, saga_tree map[uint8][]*models.Saga, layer_number uint8, new_layer map[string]*models.Saga) map[uint8][]*models.Saga {

	span, _ := opentracing.StartSpanFromContext(ctx, "registrationUC.FillNewSagaTreeLayer")
	defer span.Finish()

	result := saga_tree

	var list_of_saga []*models.Saga = nil

	for _, saga := range new_layer {
		list_of_saga = append(list_of_saga, saga)
	}

	result[layer_number] = list_of_saga

	return result
}

func (regUC registrationUC) StartOperation(ctx context.Context, operation_type uint8, start_data map[string]interface{}) (result []*models.Event, err error) {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationUC.StartOperation")
	defer span.Finish()

	if !ValidateOperation(operation_type) {
		return nil, ErrorWrongOperation
	}

	var list_of_root_saga []*models.Saga = nil

	switch operation_type {
	case OperationCreateUser | OperationAddAccount:
		{
			list_of_root_saga_types, is_exist := OperationsRootsSagas[operation_type]
			if !is_exist {
				return nil, ErrorNoRootsSagas
			}
			list_of_root_saga, err = regUC.CreateSagaTree(ctxWithTrace, list_of_root_saga_types, operation_type, start_data)
			if err != nil {
				return nil, err
			}
		}
	default:
		return nil, ErrorWrongOperation
	}

	result = nil

	for _, saga := range list_of_root_saga {
		new_events, err := regUC.ProcessingSagaAndEvents(ctxWithTrace, saga.Saga_uuid, uuid.Nil, true, nil)
		if err != nil {
			return nil, err
		}
		result = append(result, new_events...)
	}

	return result, nil
}

func (regUC registrationUC) CheckEventDataIsReady(ctx context.Context, event_uuid uuid.UUID, saga_uuid uuid.UUID) (result bool) {
	result = true

	saga, err := regUC.registrationRepo.GetSaga(ctx, saga_uuid)
	if err != nil {
		result = false
	} else {
		event, err := regUC.registrationRepo.GetEvent(ctx, event_uuid)
		if err != nil {
			result = false
		} else {
			saga_data := saga.Saga_data
			event_fields := event.Event_required_data

			for _, field := range event_fields {
				if _, ok := saga_data[field]; !ok {
					result = false
					break
				}
			}
		}
	}

	return result
}

func (regUC registrationUC) GetOperationStatus(ctx context.Context, saga_uuid uuid.UUID) (status string, err error) {
	err = nil
	status = OperationStatusUnknown

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationUC.GetOperationStatus")
	defer span.Finish()

	saga, err := regUC.registrationRepo.GetSaga(ctx, saga_uuid)
	if err != nil {
		return status, ErrorSagaWasNotFound
	}

	switch saga.Saga_status {
	case SagaStatusError, SagaStatusFallBackError, SagaStatusFallBackInProcess, SagaStatusFallBackSuccess:
		status = OperationStatusFailed
		return status, nil
	case SagaStatusInProcess, SagaStatusCreated:
		status = OperationStatusInProgress
		return status, nil
	case SagaStatusCompleted:
		status = OperationStatusSuccess
		break
	}

	if status == OperationStatusSuccess {
		saga_connections, local_err := regUC.registrationRepo.GetSagaConnectionsCurrentSaga(ctxWithTrace, saga_uuid)
		if local_err != nil {
			return OperationStatusUnknown, local_err
		}
		if saga_connections != nil && len(saga_connections.List_of_connetcions) > 0 {
			for _, saga_connection := range saga_connections.List_of_connetcions {
				saga_status, sub_saga_err := regUC.GetOperationStatus(ctxWithTrace, saga_connection.Next_saga_uuid)
				if sub_saga_err != nil {
					return OperationStatusUnknown, sub_saga_err
				}
				if saga_status == OperationStatusFailed {
					return OperationStatusFailed, nil
				} else if saga_status == OperationStatusInProgress {
					return OperationStatusInProgress, nil
				}
			}
		}
	}

	return status, err
}

func NewRegistrationUseCase(cfg *config.Config, registration_repo registration.Repository, log logger.Logger) registration.UseCase {
	return &registrationUC{cfg: cfg, registrationRepo: registration_repo, logger: log}
}
