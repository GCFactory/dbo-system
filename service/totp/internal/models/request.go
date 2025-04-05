package models

import "github.com/google/uuid"

type TOTPEnroll struct {
	TotpId     string `json:"totp_id"`
	TotpSecret string `json:"totp_secret"`
	TotpUrl    string `json:"totp_url"`
}

type DefaultHttpRequest struct {
	Status int    `json:"status"`
	Info   string `json:"info"`
}

type TOTPVerify struct {
	Status string `json:"status" validate:"omitempty"`
}

type TOTPValidate struct {
	Status string `json:"status"`
}

type TOTPEnable struct {
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

type GetTotpUrlBody struct {
	UserId uuid.UUID `json:"user_id"`
}

type GetTotpUrlResponse struct {
	TotpUrl string `json:"totp_url"`
}
