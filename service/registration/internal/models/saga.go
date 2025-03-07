package models

import (
	"github.com/google/uuid"
	"time"
)

type Operation struct {
	Operation_uuid   uuid.UUID `json:"operation_uuid" db:"operation_uuid" validate:"len=36 required uuid"`
	Operation_name   string    `json:"operation_name" db:"operation_name" validate:"required"`
	Create_time      time.Time `json:"create_time" db:"create_time" validate:"required,datetime"`
	Last_time_update time.Time `json:"last_time_update" db:"last_time_update" validate:"required,datetime"`
}

type SagaFromDB struct {
	Saga_uuid      uuid.UUID `json:"saga_uuid" db:"saga_uuid" validate:"len=36 required uuid"`
	Saga_status    uint      `json:"saga_status" db:"saga_status" validate:"min=0 max=255 required"`
	Saga_type      uint8     `json:"saga_type" db:"saga_type" validate:"min=0 max=9 required"`
	Saga_name      string    `json:"saga_name" db:"saga_name" validate:"max=50 required"`
	Saga_data      string    `json:"saga_data" db:"saga_data" validate:"required_with"`
	Operation_uuid uuid.UUID `json:"operation_uuid" db:"operation_uuid" validate:"len=36 required uuid"`
}

type Saga struct {
	Saga_uuid      uuid.UUID              `json:"saga_uuid" db:"saga_uuid" validate:"len=36 required uuid"`
	Saga_status    uint                   `json:"saga_status" db:"saga_status" validate:"min=0 max=255 required"`
	Saga_type      uint8                  `json:"saga_type" db:"saga_type" validate:"min=0 max=9 required"`
	Saga_name      string                 `json:"saga_name" db:"saga_name" validate:"max=50 required"`
	Saga_data      map[string]interface{} `json:"saga_data" db:"saga_data" validate:"required_with"`
	Operation_uuid uuid.UUID              `json:"operation_uuid" db:"operation_uuid" validate:"len=36 required uuid"`
}

type Event struct {
	Event_uuid          uuid.UUID `json:"event_uuid" db:"event_uuid" validate:"len=36 required uuid"`
	Saga_uuid           uuid.UUID `json:"saga_uuid" db:"saga_uuid" validate:"len=36 required uuid"`
	Event_status        uint8     `json:"event_status" db:"event_status" validate:"min=0 max=255 required"`
	Event_name          string    `json:"event_name" db:"event_name" validate:"max=64 required"`
	Event_is_roll_back  bool      `json:"event_is_roll_back" db:"event_is_roll_back" validate:"bool required"`
	Event_result        string    `json:"event_result" db:"event_result" validate:"json required"`
	Event_required_data []string  `json:"event_required_data" db:"event_required_data" validate:"json required"`
	Event_rollback_uuid uuid.UUID `json:"event_rollback_uuid" db:"event_rollback_uuid" validate:"len=36 required uuid"`
}
type EventFromDB struct {
	Event_uuid          uuid.UUID `json:"event_uuid" db:"event_uuid" validate:"len=36 required uuid"`
	Saga_uuid           uuid.UUID `json:"saga_uuid" db:"saga_uuid" validate:"len=36 required uuid"`
	Event_status        uint8     `json:"event_status" db:"event_status" validate:"min=0 max=255 required"`
	Event_name          string    `json:"event_name" db:"event_name" validate:"max=64 required"`
	Event_is_roll_back  bool      `json:"event_is_roll_back" db:"event_is_roll_back" validate:"bool required"`
	Event_result        string    `json:"event_result" db:"event_result" validate:"json required"`
	Event_required_data string    `json:"event_required_data" db:"event_required_data" validate:"json required"`
	Event_rollback_uuid uuid.UUID `json:"event_rollback_uuid" db:"event_rollback_uuid" validate:"len=36 required uuid"`
}

type SagaListEvents struct {
	EventList []uuid.UUID `json:"list_of_events" validate:"required"`
}

type SagaConnection struct {
	Current_saga_uuid     uuid.UUID `json:"current_saga_uuid" validate:"len=36 required uuid"`
	Next_saga_uuid        uuid.UUID `json:"next_saga_uuid" validate:"len=36 required uuid"`
	Acc_connection_status uint8     `json:"acc_connection_status" validate:"min=0 max=255 required"`
}

type ListOfSagaConnections struct {
	List_of_connetcions []*SagaConnection `json:"list_of_connections" validate:"json required"`
}

type ListOfSaga struct {
	ListId []uuid.UUID `json:"list_of_saga" validate:"json required"`
}

type SagaDepend struct {
	Parents  []string
	Children []string
}
