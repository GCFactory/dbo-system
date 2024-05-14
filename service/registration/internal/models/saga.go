package models

import "github.com/google/uuid"

type Saga struct {
	Saga_uuid           uuid.UUID `json:"saga_uuid" db:"saga_uuid" validate:"len=36 required"`
	Saga_status         uint      `json:"saga_status" db:"saga_status" validate:"min=0 max=255 required"`
	Saga_type           uint      `json:"saga_type" db:"saga_type" validate:"min=0 max=9 required"`
	Saga_name           string    `json:"saga_name" db:"saga_name" validate:"max=50 required"`
	Saga_list_of_events string    `json:"list_of_events" db:"list_of_events" validate:"json required"`
}

type SagaListEvents struct {
	EventList []string `json:"list_of_events" validate:"required"`
}
