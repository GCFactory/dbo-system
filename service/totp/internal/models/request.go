package models

import "github.com/google/uuid"

type TOTPEnroll struct {
	Base32    string `json:"base32"`
	OTPathURL string `json:"otpath_url"`
}

type TOTPVerify struct {
	Status string `json:"info" validate:"omitempty"`
}

type TOTPValidate struct {
	Status string `json:"status"`
}

type TOTPDisable struct {
	Status string `json:"status"`
}

type TOTPRequest struct {
	// USER id
	UserId uuid.UUID `json:"user_id" db:"user_id"`
	Code   string    `json:"code" validate:"omitempty,lte=12"`
}

type TOTPBase struct {
	OTPValid bool `json:"otp_valid"`
}
