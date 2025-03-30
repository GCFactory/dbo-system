package models

import "github.com/google/uuid"

type DefaultHttpResponse struct {
	Status int    `json:"status"`
	Info   string `json:"info"`
}
type AddUserSettings struct {
	UserId            uuid.UUID `json:"user_id" validate:"required"`
	EmailNotification bool      `json:"email_notification" validate:"required"`
	Email             string    `json:"email" validate:"required,email"`
}

type UpdateUserSettings struct {
	UserId            uuid.UUID `json:"user_id" validate:"required"`
	EmailNotification bool      `json:"email_notification" validate:"omitempty"`
	Email             string    `json:"email" validate:"required,email"`
}

type DeleteUserSettings struct {
	UserId uuid.UUID `json:"user_id" validate:"required"`
}
