package models

import "github.com/google/uuid"

type Saga struct {
	Saga_uuid           uuid.UUID      `json:"saga_uuid" db:"saga_uuid" validate:"len=36 required uuid"`
	Saga_status         uint           `json:"saga_status" db:"saga_status" validate:"min=0 max=255 required"`
	Saga_type           uint           `json:"saga_type" db:"saga_type" validate:"min=0 max=9 required"`
	Saga_name           string         `json:"saga_name" db:"saga_name" validate:"max=50 required"`
	Saga_list_of_events SagaListEvents `json:"list_of_events" db:"list_of_events" validate:"json required"`
}

type Event struct {
	Event_uuid          uuid.UUID `json:"event_uuid" db:"event_uuid" validate:"len=36 required uuid"`
	Event_status        uint      `json:"event_status" db:"event_status" validate:"min=0 max=255 required"`
	Event_name          string    `json:"event_name" db:"event_name" validate:"max=64 required"`
	Event_result        string    `json:"event_result" db:"event_result" validate:"json required"`
	Event_rollback_uuid uuid.UUID `json:"event_rollback_uuid" db:"event_rollback_uuid" validate:"len=36 required uuid"`
}

type SagaListEvents struct {
	EventList []uuid.UUID `json:"list_of_events" validate:"required"`
}

type SagaConnection struct {
	Current_saga_uuid uuid.UUID `json:"current_saga_uuid" validate:"len=36 required uuid"`
	Next_saga_uuid    uuid.UUID `json:"next_saga_uuid" validate:"len=36 required uuid"`
}

type ListOfSagaConnections struct {
	List_of_connetcions []SagaConnection `json:"list_of_connetcions" validate:"json required"`
}
