package models

import "github.com/google/uuid"

type UserNotificationInfo struct {
	UserUuid   uuid.UUID `json:"user_uuid" db:"user_uuid"`
	EmailUsage bool      `json:"email_usage" db:"email_usage"`
	Email      string    `json:"email" db:"email"`
}
