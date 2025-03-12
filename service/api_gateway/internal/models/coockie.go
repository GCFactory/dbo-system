package models

import (
	"github.com/google/uuid"
	"time"
)

type Token struct {
	ID        uuid.UUID
	Data      uuid.UUID
	Live_time time.Duration
}
