package models

import "github.com/google/uuid"

type DefaultHttpResponse struct {
	Status int    `json:"status"`
	Info   string `json:"info"`
}

type UpdateUserSettings struct {
	UserId            uuid.UUID `json:"user_id" validate:"required"`
	EmailNotification bool      `json:"email_notification" validate:"omitempty"`
	Email             string    `json:"email" validate:"required,email"`
}

type GetUserSettings struct {
	UserId string `json:"user_id" validate:"required"`
}
