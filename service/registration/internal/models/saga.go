package models

import (
	"github.com/google/uuid"
)

type Saga struct {
	Saga_uuid   uuid.UUID `json:"saga_uuid" db:"saga_uuid" validate:"len=36 required uuid"`
	Saga_status uint      `json:"saga_status" db:"saga_status" validate:"min=0 max=255 required"`
	Saga_type   uint8     `json:"saga_type" db:"saga_type" validate:"min=0 max=9 required"`
	Saga_name   string    `json:"saga_name" db:"saga_name" validate:"max=50 required"`
}

type Event struct {
	Event_uuid          uuid.UUID `json:"event_uuid" db:"event_uuid" validate:"len=36 required uuid"`
	Saga_uuid           uuid.UUID `json:"saga_uuid" db:"saga_uuid" validate:"len=36 required uuid"`
	Event_status        uint8     `json:"event_status" db:"event_status" validate:"min=0 max=255 required"`
	Event_name          string    `json:"event_name" db:"event_name" validate:"max=64 required"`
	Event_is_roll_back  bool      `json:"event_is_roll_back" db:"event_is_roll_back" validate:"bool required"`
	Event_result        string    `json:"event_result" db:"event_result" validate:"json required"`
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

type SagaDepend struct {
	Parents  []string
	Children []string
}

//type EventInfo struct {
//	EventType            uint
//	EventRealisationPath string
//	EventRealisation     *usecase.EventRealisation
//}

//type Passport struct {
//	Passport_uuid       uuid.UUID `json:"passport_uuid" db:"passport_uuid" validate:"required,uuid4"`
//	Passport_series     string    `json:"passport_series" db:"passport_series" validate:"required,len=4,number"`
//	Passport_number     string    `json:"passport_number" db:"passport_number" validate:"required,len=6,number"`
//	Name                string    `json:"name" db:"name" validate:"required"`
//	Surname             string    `json:"surname" db:"surname" validate:"required"`
//	Patronimic          string    `json:"patronimic" db:"patronimic" validate:"required"`
//	Birth_date          string    `json:"birth_date" db:"birth_date" validate:"required,datetime"`
//	Birth_location      string    `json:"birth_location" db:"birth_location" validate:"required"`
//	Pick_up_point       string    `json:"pick_up_point" db:"pick_up_point" validate:"required"`
//	Authority           string    `json:"authority" db:"authority" validate:"required,len=7"`
//	Authority_date      string    `json:"authority_date" db:"authority_date" validate:"required,datetime"`
//	Registration_adress string    `json:"registration_adress" db:"registration_adress" validate:"required"`
//}
//
//type ListOfAccounts struct {
//	Data []uuid.UUID `json:"accounts" db:"accounts"`
//}
//
//type User struct {
//	User_uuid     uuid.UUID      `json:"user_uuid" db:"user_uuid" validate:"required,uuid4"`
//	Passport_uuid uuid.UUID      `json:"passport_uuid" db:"passport_uuid" validate:"required,uuid4"`
//	User_inn      string         `json:"user_inn" db:"user_inn" validate:"required,number,len=20"`
//	User_accounts ListOfAccounts `json:"user_accounts" db:"user_accounts" validate:"required,json"`
//}
