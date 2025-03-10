package models

import (
	"github.com/google/uuid"
	"time"
)

type Token struct {
	ID          uuid.UUID
	Live_time   time.Duration
	Date_expire time.Time
}
