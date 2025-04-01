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
	"sync"
	"time"
)

type registrationUC struct {
	cfg              *config.Config
	registrationRepo registration.Repository
	logger           logger.Logger
}

var mutex sync.Mutex

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

	event_data, ok := RequiredEventListOfData[event_type]
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

	saga, err := regUC.registrationRepo.GetSaga(ctxWithTrace, event.Saga_uuid)
	if err != nil {
		return nil, err
	}

	saga_group_rollback, ok := RevertEvent[saga.Saga_type]
	if ok {
		saga_roll_back, ok := saga_group_rollback[saga.Saga_name]
		if ok {
			event_type_roll_back, exists := saga_roll_back[event.Event_name]
			if exists {

				revert_event, err := regUC.CreateEvent(ctxWithTrace, event_type_roll_back, event.Saga_uuid, true, event_uuid)
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
		}
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

func (regUC registrationUC) CreateOperation(ctx context.Context, operation_type uint8) (*models.Operation, error) {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationUC.CreateOperation")
	defer span.Finish()

	operation_name, err := OperationNameFromCode(operation_type)
	if err != nil {
		return nil, err
	}

	operation := &models.Operation{
		Operation_uuid:   uuid.New(),
		Operation_name:   operation_name,
		Create_time:      time.Now(),
		Last_time_update: time.Now(),
	}

	err = regUC.registrationRepo.CreateOperation(ctxWithTrace, operation)
	if err != nil {
		return nil, err
	}

	return operation, nil

}

func (regUC registrationUC) UpdateOperationTimeEdit(ctx context.Context, operation_id uuid.UUID, new_time time.Time) error {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationUC.UpdateOperationTimeEdit")
	defer span.Finish()

	operation, err := regUC.registrationRepo.GetOperation(ctxWithTrace, operation_id)
	if err != nil {
		return ErrorNoOperationFound
	}

	operation.Last_time_update = new_time
	err = regUC.registrationRepo.UpdateOperation(ctxWithTrace, operation)
	if err != nil {
		return err
	}

	return nil

}

func (regUC registrationUC) DeleteOperation(ctx context.Context, operation_id uuid.UUID) error {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationUC.DeleteOperation")
	defer span.Finish()

	_, err := regUC.registrationRepo.GetOperation(ctxWithTrace, operation_id)
	if err != nil {
		return ErrorNoOperationFound
	}

	err = regUC.registrationRepo.DeleteOperation(ctxWithTrace, operation_id)
	if err != nil {
		return err
	}

	return nil

}

func (regUC registrationUC) CreateSaga(ctx context.Context, operation_id uuid.UUID, saga_type string, saga_group uint8, saga_data map[string]interface{}) (saga *models.Saga, err error) {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationUC.CreateSaga")
	defer span.Finish()

	if !ValidateSagaType(saga_type) {
		return nil, ErrorInvalidSagaType
	} else if !ValidateSagaGroup(saga_group) {
		return nil, ErrorInvalidSagGroup
	}

	saga = &models.Saga{
		Saga_uuid:      uuid.New(),
		Saga_status:    SagaStatusUndefined,
		Saga_type:      saga_group,
		Saga_name:      saga_type,
		Saga_data:      saga_data,
		Operation_uuid: operation_id,
	}

	err = regUC.registrationRepo.CreateSaga(ctxWithTrace, saga)
	if err != nil {
		return nil, err
	}

	events, ok := StartEventsListForSagaType[saga_type]
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

	result = nil
	err = nil

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

	var child_revert bool = true

	for _, connection := range list_of_child_connections.List_of_connetcions {
		//child_saga_uuid := connection.Next_saga_uuid
		connection_status := connection.Acc_connection_status
		if connection_status == SagaConnectionStatusSuccess {
			child_revert = false
			//new_events, err := regUC.RevertSaga(ctxWithTrace, child_saga_uuid)
			//if err != nil {
			//	return nil, err
			//}
			//result = append(result, new_events...)
		}
	}

	if !child_revert {
		return result, ErrorNotAllChildReverted
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

	var roll_back_event_exists bool = false

	if list_of_events != nil {
		for _, event := range list_of_events {
			if event.Event_status == EventStatusCompleted && !event.Event_is_roll_back {
				reverted_event, err := regUC.RevertEvent(ctxWithTrace, event.Event_uuid)
				if err != nil {
					return result, err
				}
				if reverted_event != nil {
					roll_back_event_exists = true
					result = append(result, reverted_event)
				}
			}
		}
	}

	if roll_back_event_exists {
		saga.Saga_status = SagaStatusFallBackInProcess
	} else {
		saga.Saga_status = SagaStatusFallBackSuccess
	}

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
					event_saga, err := regUC.registrationRepo.GetSaga(ctxWithTrace, event.Saga_uuid)
					if err != nil {
						return nil, err
					}
					if validation_err := AdditionalValidation(event_saga.Saga_type, event.Event_name, event_saga.Saga_data); validation_err != nil {

						event.Event_status = EventStatusError

						error_info := make(map[string]interface{})
						error_info["local validation error"] = validation_err.Error()

						error_info_str, err := json.Marshal(&error_info)
						if err != nil {

							return nil, err

						}

						event.Event_result = string(error_info_str)

						err = regUC.registrationRepo.UpdateEvent(ctxWithTrace, event)
						if err != nil {
							return nil, err
						}

						new_events, err := regUC.ProcessingSagaAndEvents(ctxWithTrace, event.Saga_uuid, uuid.Nil, false, nil)
						if err != nil {
							return nil, err
						}
						result = append(result, new_events...)

					} else {
						event.Event_status = EventStatusInProgress
						err = regUC.registrationRepo.UpdateEvent(ctxWithTrace, event)
						if err != nil {
							return result, err
						}
						result = append(result, event)
					}
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

						saga_group := saga.Saga_type
						saga_type := saga.Saga_name

						update_fileds, ok := SagaGroupResultDataUpdate[saga_group]
						if ok {

							saga_update_fields, is_exist := update_fileds[saga_type]
							if is_exist {

								for field_name, field_value := range data {

									if slices.Contains(saga_update_fields, field_name) {
										saga.Saga_data[field_name] = field_value
									}

								}
							}

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
						new_events, err := regUC.ProcessingSagaAndEvents(ctxWithTrace, event.Saga_uuid, event.Event_rollback_uuid, false, nil)
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
					//regUC.logger.Debug("ROLLBACK")
					new_event, err := regUC.RevertEvent(ctxWithTrace, event.Event_uuid)
					if err != nil {
						return result, err
					}
					if new_event == nil {

						event.Event_status = EventStatusFallBackCompleted
						err = regUC.registrationRepo.UpdateEvent(ctxWithTrace, event)
						if err != nil {
							return result, err
						}
						new_events, err := regUC.ProcessingSagaAndEvents(
							ctxWithTrace,
							event.Saga_uuid,
							uuid.Nil,
							true,
							data,
						)
						if err != nil {
							return result, err
						}
						result = append(result, new_events...)

					} else {
						//regUC.logger.Debug("ROLLBACK END")
						new_events, err := regUC.ProcessingSagaAndEvents(
							ctxWithTrace,
							uuid.Nil,
							new_event.Event_uuid,
							true,
							nil)
						if err != nil {
							return nil, err
						}
						result = append(result, new_events...)
					}
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

		saga, err := regUC.registrationRepo.GetSaga(ctxWithTrace, event.Saga_uuid)
		if err != nil {
			return nil, err
		}

		err = regUC.UpdateOperationTimeEdit(ctxWithTrace, saga.Operation_uuid, time.Now())
		if err != nil {
			return nil, err
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
				saga_connections, local_err := regUC.registrationRepo.GetSagaConnectionsNextSaga(ctxWithTrace, saga_uuid)
				if local_err != nil {
					return result, local_err
				}
				saga_data := saga.Saga_data
				for _, saga_connection := range saga_connections.List_of_connetcions {

					parent_saga, local_err := regUC.registrationRepo.GetSaga(ctxWithTrace, saga_connection.Current_saga_uuid)
					if local_err != nil {
						return result, local_err
					}

					for parent_saga_data_filed_name, parent_saga_data_field_data := range parent_saga.Saga_data {

						if saga_data == nil {
							saga_data = make(map[string]interface{})
						}
						saga_data[parent_saga_data_filed_name] = parent_saga_data_field_data

					}

					saga.Saga_data = saga_data
					err = regUC.registrationRepo.UpdateSaga(ctxWithTrace, saga)
					if err != nil {
						return result, err
					}

				}

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
					new_events, err := regUC.ProcessingSagaAndEvents(ctxWithTrace, saga.Saga_uuid, uuid.Nil, false, nil)
					if err != nil {
						return result, err
					}
					result = append(result, new_events...)
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
						if connection.Acc_connection_status != SagaConnectionStatusFallBack &&
							connection.Acc_connection_status != SagaConnectionStatusFailed {
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
						for _, new_event := range new_events {
							if new_event != nil {
								new_events_process, err := regUC.ProcessingSagaAndEvents(
									ctxWithTrace,
									uuid.Nil,
									new_event.Event_uuid,
									true,
									data,
								)
								if err != nil {
									return nil, err
								}
								result = append(result, new_events_process...)
							}
						}
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
					new_events, err := regUC.ProcessingSagaAndEvents(ctxWithTrace, connection.Current_saga_uuid, uuid.Nil, false, nil)
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
						new_events, err := regUC.ProcessingSagaAndEvents(ctxWithTrace, connetion.Current_saga_uuid, uuid.Nil, false, nil)
						if err != nil {
							return result, err
						}
						result = append(result, new_events...)
					}
				}
			}
		case SagaStatusError:
			{
				//regUC.logger.Debug("SagaStatusError")
				all_completed := true
				for _, event := range events {
					if event.Event_is_roll_back &&
						(event.Event_status == EventStatusCreated ||
							event.Event_status == EventStatusInProgress) ||
						!event.Event_is_roll_back &&
							(event.Event_status == EventStatusCreated ||
								event.Event_status == EventStatusInProgress ||
								event.Event_status == EventStatusCompleted ||
								event.Event_status == EventStatusFallBackInProcess) {
						//all_completed = false
						if !event.Event_is_roll_back && event.Event_status == EventStatusCompleted {
							new_events, err := regUC.ProcessingSagaAndEvents(ctxWithTrace, uuid.Nil, event.Event_uuid, false, nil)
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

		err = regUC.UpdateOperationTimeEdit(ctxWithTrace, saga.Operation_uuid, time.Now())
		if err != nil {
			return nil, err
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
	if regUC.CheckEventDataIsReady(ctx, event_uuid, saga.Saga_uuid) {
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

func (regUC registrationUC) CreateSagaTree(ctx context.Context, operation_type uint8, list_root_saga_types []string, saga_group uint8, start_data map[string]interface{}) (list_root_saga []*models.Saga, operation_id uuid.UUID, err error) {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationUC.CreateSagaTree")
	defer span.Finish()

	operation_id = uuid.Nil

	operation, err := regUC.CreateOperation(ctxWithTrace, operation_type)
	if err != nil {
		return nil, operation_id, err
	}
	operation_id = operation.Operation_uuid

	list_root_saga = make([]*models.Saga, 0)

	// create root layer of saga tree
	for _, saga_type := range list_root_saga_types {
		saga, err := regUC.CreateSaga(ctxWithTrace, operation.Operation_uuid, saga_type, saga_group, start_data)
		if err != nil {
			if list_root_saga != nil {
				for _, saga := range list_root_saga {
					err = regUC.DeleteSaga(ctxWithTrace, saga.Saga_uuid)
					if err != nil {
						return nil, operation_id, err
					}
				}
			}
			return nil, operation_id, err
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

	ListOfSagaDependGroup, ok := ListOfSagaDepend[saga_group]
	if !ok {
		return nil, operation_id, ErrorNoSagaGroupDepend
	}

	for {
		for saga_type, _ := range previous_saga_tree_layer {

			saga_depends, is_exist := ListOfSagaDependGroup[saga_type]
			if !is_exist {
				saga_tree = regUC.FillNewSagaTreeLayer(ctxWithTrace, saga_tree, layer_number, current_saga_tree_layer)
				err = regUC.ClearSagaTree(ctxWithTrace, saga_tree)
				if err != nil {
					return nil, operation_id, err
				}
				return nil, operation_id, ErrorNoDependsData
			}
			//if saga_depends.Children != nil {
			for _, child_type := range saga_depends.Children {
				if !ValidateSagaType(child_type) {
					saga_tree = regUC.FillNewSagaTreeLayer(ctxWithTrace, saga_tree, layer_number, current_saga_tree_layer)
					err = regUC.ClearSagaTree(ctxWithTrace, saga_tree)
					if err != nil {
						return nil, operation_id, err
					}
					return nil, operation_id, ErrorInvalidSagaType
				}
				_, is_exist = current_saga_tree_layer[child_type]
				if !is_exist {
					saga_child, local_err := regUC.CreateSaga(ctxWithTrace, operation.Operation_uuid, child_type, saga_group, nil)
					if local_err != nil {
						saga_tree = regUC.FillNewSagaTreeLayer(ctxWithTrace, saga_tree, layer_number, current_saga_tree_layer)
						err = regUC.ClearSagaTree(ctxWithTrace, saga_tree)
						if err != nil {
							return nil, operation_id, err
						}
						return nil, operation_id, local_err
					}
					current_saga_tree_layer[child_type] = saga_child

					child_depends, is_exist := ListOfSagaDependGroup[child_type]
					if !is_exist {
						saga_tree = regUC.FillNewSagaTreeLayer(ctxWithTrace, saga_tree, layer_number, current_saga_tree_layer)
						err = regUC.ClearSagaTree(ctxWithTrace, saga_tree)
						if err != nil {
							return nil, operation_id, err
						}
						return nil, operation_id, ErrorNoDependsData
					}

					for _, parent_type := range child_depends.Parents {
						if !ValidateSagaType(parent_type) {
							saga_tree = regUC.FillNewSagaTreeLayer(ctxWithTrace, saga_tree, layer_number, current_saga_tree_layer)
							err = regUC.ClearSagaTree(ctxWithTrace, saga_tree)
							if err != nil {
								return nil, operation_id, err
							}
							return nil, operation_id, ErrorInvalidSagaType
						}
						parent_saga, is_exist := previous_saga_tree_layer[parent_type]
						if !is_exist {
							saga_tree = regUC.FillNewSagaTreeLayer(ctxWithTrace, saga_tree, layer_number, current_saga_tree_layer)
							err = regUC.ClearSagaTree(ctxWithTrace, saga_tree)
							if err != nil {
								return nil, operation_id, err
							}
							return nil, operation_id, ErrorSagaWasNotFoundInPreviousLayer
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
								return nil, operation_id, err
							}
							return nil, operation_id, ErrorInvalidSagaType
						}
					}
				}
			}
			//}
			saga_tree = regUC.FillNewSagaTreeLayer(ctxWithTrace, saga_tree, layer_number, current_saga_tree_layer)
		}

		if len(current_saga_tree_layer) == 0 {
			break
		}

		saga_tree = regUC.FillNewSagaTreeLayer(ctxWithTrace, saga_tree, layer_number, current_saga_tree_layer)
		layer_number++

		previous_saga_tree_layer = current_saga_tree_layer
		current_saga_tree_layer = nil
		current_saga_tree_layer = make(map[string]*models.Saga)

	}

	return list_root_saga, operation_id, nil

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

func (regUC registrationUC) StartOperation(ctx context.Context, operation_type uint8, start_data map[string]interface{}) (result []*models.Event, operation_uuid uuid.UUID, err error) {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationUC.StartOperation")
	defer span.Finish()

	operation_id := uuid.Nil

	if !ValidateOperation(operation_type) {
		return nil, operation_id, ErrorWrongOperation
	}

	var list_of_root_saga []*models.Saga = nil

	switch operation_type {
	case OperationCreateUser,
		OperationAddAccount,
		OperationAddAccountCache,
		OperationWidthAccountCache,
		OperationCloseAccount,
		OperationGetUserData,
		OperationGetAccountData,
		OperationGroupUpdateUserPassword,
		OperationGetUserDataByLogin,
		OperationCheckUserPassword:
		{
			list_of_root_saga_types, is_exist := OperationsRootsSagas[operation_type]
			if !is_exist {
				return nil, operation_id, ErrorNoRootsSagas
			}
			list_of_root_saga, operation_id, err = regUC.CreateSagaTree(ctxWithTrace, operation_type, list_of_root_saga_types, operation_type, start_data)
			if err != nil {
				return nil, operation_id, err
			}
		}
	default:
		return nil, operation_id, ErrorWrongOperation
	}

	result = nil

	for _, saga := range list_of_root_saga {
		new_events, err := regUC.ProcessingSagaAndEvents(ctxWithTrace, saga.Saga_uuid, uuid.Nil, true, nil)
		if err != nil {
			return nil, operation_id, err
		}
		result = append(result, new_events...)
	}

	return result, operation_id, nil
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

func (regUC registrationUC) GetOperationStatus(ctx context.Context, operation_id uuid.UUID) (status string, err error) {
	err = nil
	status = OperationStatusUnknown

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationUC.GetOperationStatus")
	defer span.Finish()

	_, err = regUC.registrationRepo.GetOperation(ctxWithTrace, operation_id)
	if err != nil {
		return status, err
	}

	saga_uuids, err := regUC.registrationRepo.GetOperationSaga(ctxWithTrace, operation_id)
	if err != nil {
		return status, err
	}

	saga_list := make([]*models.Saga, 0)

	for _, saga_uuid := range saga_uuids.ListId {
		saga, err := regUC.registrationRepo.GetSaga(ctxWithTrace, saga_uuid)
		if err != nil {
			return status, err
		}
		saga_list = append(saga_list, saga)

	}

	for _, saga := range saga_list {
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
	}

	return status, err
}

func (regUC registrationUC) GetSagaChildrenWithoutChildren(ctx context.Context, saga_uuid uuid.UUID) (children []*models.Saga, err error) {
	children = nil
	err = nil

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationUC.GetSagaChildrenWithoutChildren")
	defer span.Finish()

	saga, err := regUC.registrationRepo.GetSaga(ctxWithTrace, saga_uuid)
	if err != nil {
		err = ErrorSagaWasNotFound
	} else {

		saga_children, local_err := regUC.registrationRepo.GetSagaConnectionsCurrentSaga(ctxWithTrace, saga_uuid)
		if local_err != nil {
			return children, local_err
		} else {
			if len(saga_children.List_of_connetcions) == 0 {
				children = append(children, saga)
			} else {
				for _, child := range saga_children.List_of_connetcions {
					child_children, local_err := regUC.registrationRepo.GetSagaConnectionsCurrentSaga(ctxWithTrace, child.Next_saga_uuid)
					if local_err != nil {
						return children, local_err
					}
					if len(child_children.List_of_connetcions) == 0 {
						child_saga, err := regUC.registrationRepo.GetSaga(ctxWithTrace, child.Next_saga_uuid)
						if err != nil {
							return children, err
						}
						children = append(children, child_saga)
					} else {
						sub_childs, err := regUC.GetSagaChildrenWithoutChildren(ctxWithTrace, child.Next_saga_uuid)
						if err != nil {
							return children, err
						}
						children = append(children, sub_childs...)
					}
				}
			}
		}
	}

	return children, err
}

func (regUC registrationUC) GetOperationResultData(ctx context.Context, operation_uuid uuid.UUID) (data map[string]interface{}, err error) {

	data = make(map[string]interface{})
	err = nil

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationUC.GetOperationResultData")
	defer span.Finish()

	operation, err := regUC.registrationRepo.GetOperation(ctxWithTrace, operation_uuid)
	if err != nil {
		err = ErrorNoOperationFound
	} else {
		saga_uuids, err := regUC.registrationRepo.GetOperationSaga(ctxWithTrace, operation_uuid)
		if err != nil {
			return data, err
		}

		for _, saga_uuid := range saga_uuids.ListId {
			saga, err := regUC.registrationRepo.GetSaga(ctxWithTrace, saga_uuid)
			if err != nil {
				return data, err
			}

			events_uuid, err := regUC.registrationRepo.GetListOfSagaEvents(ctxWithTrace, saga.Saga_uuid)
			if err != nil {
				return nil, err
			}

			var events []*models.Event

			for _, event_uuid := range events_uuid.EventList {

				event, err := regUC.registrationRepo.GetEvent(ctxWithTrace, event_uuid)
				if err != nil {

					return nil, err

				}

				events = append(events, event)

			}

			if saga.Saga_status == SagaStatusError ||
				saga.Saga_status == SagaStatusFallBackError {

				for _, event := range events {
					if event.Event_status == EventStatusError ||
						event.Event_status == EventStatusFallBackError {
						var event_result_data map[string]interface{}
						err := json.Unmarshal([]byte(event.Event_result), &event_result_data)
						if err != nil {
							return data, err
						}

						errors, ok := data["errors"].(map[string]interface{})
						if !ok {
							data["errors"] = make(map[string]interface{})
							errors = data["errors"].(map[string]interface{})
						}

						// Очистка всех данных кроме данных об ошибках
						data = nil
						data = make(map[string]interface{})

						saga_errors, exist := errors[saga.Saga_name].(map[string]interface{})
						if !exist {
							errors[saga.Saga_name] = make(map[string]interface{})
							saga_errors = errors[saga.Saga_name].(map[string]interface{})
						}

						event_name := event.Event_name
						error_data := make(map[string]interface{})

						for field_name, field_data := range event_result_data {
							error_data[field_name] = field_data
						}

						error_data["is_roll_back"] = event.Event_is_roll_back
						saga_errors[event_name] = error_data
						errors[saga.Saga_name] = saga_errors
						data["errors"] = errors
					}

				}

			} else {

				saga_group_result_data, ok := SagaGroupDataIsResult[saga.Saga_type]
				if ok {
					saga_result_data, is_exist := saga_group_result_data[saga.Saga_name]
					if is_exist {
						for _, event := range events {
							event_result_data, exists := saga_result_data[event.Event_name]
							if exists {

								var result_data map[string]interface{}
								err := json.Unmarshal([]byte(event.Event_result), &result_data)
								if err != nil {
									return data, err
								}

								for field_name, field_data := range result_data {
									if slices.Contains(event_result_data, field_name) {
										data[field_name] = field_data
									}
								}

							}

						}
					}
				}

			}

		}
	}

	data["time_begin"] = operation.Create_time
	data["time_end"] = operation.Last_time_update

	return data, err
}

func (regUC registrationUC) GetOperationTree(ctx context.Context, operation_uuid uuid.UUID) (data map[string]interface{}, err error) {
	data = make(map[string]interface{})
	err = nil

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationUC.GetOperationTree")
	defer span.Finish()

	operation, err := regUC.registrationRepo.GetOperation(ctxWithTrace, operation_uuid)
	if err != nil {
		return nil, err
	}

	data["operation_name"] = operation.Operation_name

	var saga_list_node []*models.SagaTreeNode
	var event_tree_node []*models.EventTreeNode
	var saga_connection_list_node []*models.SagaConnectionTree

	saga_uuids, err := regUC.registrationRepo.GetOperationSaga(ctxWithTrace, operation_uuid)
	if err != nil {
		return data, err
	}

	for _, saga_uuid := range saga_uuids.ListId {
		saga, err := regUC.registrationRepo.GetSaga(ctxWithTrace, saga_uuid)
		if err != nil {
			return data, err
		}

		events_uuid, err := regUC.registrationRepo.GetListOfSagaEvents(ctxWithTrace, saga.Saga_uuid)
		if err != nil {
			return nil, err
		}

		saga_node := &models.SagaTreeNode{
			Name:   saga.Saga_name,
			Status: uint8(saga.Saga_status),
			Id:     saga.Saga_uuid,
			Events: events_uuid.EventList,
		}

		saga_list_node = append(saga_list_node, saga_node)

		for _, event_uuid := range events_uuid.EventList {

			event, err := regUC.registrationRepo.GetEvent(ctxWithTrace, event_uuid)
			if err != nil {

				return nil, err

			}

			event_node := &models.EventTreeNode{
				Name:         event.Event_name,
				Id:           event.Event_uuid,
				Status:       event.Event_status,
				Roll_back_id: event.Event_rollback_uuid,
			}

			event_tree_node = append(event_tree_node, event_node)

		}

		saga_connections, err := regUC.registrationRepo.GetSagaConnectionsCurrentSaga(ctxWithTrace, saga_uuid)
		if err == nil {
			for _, saga_connection := range saga_connections.List_of_connetcions {

				saga_connection_node := &models.SagaConnectionTree{
					Parent_id: saga_connection.Current_saga_uuid,
					Child_id:  saga_connection.Next_saga_uuid,
				}

				saga_connection_list_node = append(saga_connection_list_node, saga_connection_node)

			}
		}
	}

	data["saga"] = saga_list_node
	data["events"] = event_tree_node
	data["saga_depend"] = saga_connection_list_node

	return data, err
}

func (regUC registrationUC) GerOperationListBetween(ctx context.Context, begin time.Time, end time.Time) ([]uuid.UUID, error) {

	span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "registrationUC.GerOperationListBetween")
	defer span.Finish()

	result, err := regUC.registrationRepo.GetOperationBetweenInterval(ctxWithTrace, begin, end)
	if err != nil {
		return nil, err
	}

	return result, nil

}

func NewRegistrationUseCase(cfg *config.Config, registration_repo registration.Repository, log logger.Logger) registration.UseCase {
	return &registrationUC{cfg: cfg, registrationRepo: registration_repo, logger: log}
}
