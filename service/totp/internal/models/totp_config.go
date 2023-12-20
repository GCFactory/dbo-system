package models

import (
	otpConfig "github.com/GCFactory/dbo-system/service/totp/pkg/otp/config"
	"github.com/google/uuid"
	"io"
	"time"
)

// TOTPConfig models
// @Description TOTP configuration
type TOTPConfig struct {
	// TOTP id
	Id uuid.UUID `json:"totp_id,omitempty" db:"totp_id" validate:"omitempty"`
	// USER id
	UserId   uuid.UUID `json:"user_id" db:"user_id"`
	IsActive bool      `json:"is_active" db:"is_active"`
	// What organization issued
	Issuer string `db:"issuer,omitempty" validate:"omitempty"`
	// Who organization issued
	AccountName string `json:"account_name" db:"account_name" validate:"omitempty,lte=60"`
	// TOTP secret
	Secret string `json:"secret,omitempty" db:"secret" validate:"omitempty,lte=40"`
	// When data created
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at" validate:"omitempty"`
	// When data updated
	UpdatedAt time.Time `json:"updated_at,omitempty" db:"updated_at" validate:"omitempty"`
}

// GenerateOpts models
// @Description GenerateOpts provides options for Generate().  The default values
// @Description are compatible with Google-Authenticator.
type GenerateOpts struct {
	// Name of the issuing Organization/Company.
	Issuer string
	// Name of the User's Account (eg, email address)
	AccountName string
	// Number of seconds a TOTP hash is valid for. Defaults to 30 seconds.
	Period uint
	// Size in size of the generated Secret. Defaults to 20 bytes.
	SecretSize uint
	// Secret to store. Defaults to a randomly generated secret of SecretSize.  You should generally leave this empty.
	Secret []byte
	// Digits to request. Defaults to 6.
	Digits otpConfig.Digits
	// Algorithm to use for HMAC. Defaults to SHA1.
	Algorithm otpConfig.Algorithm
	// Reader to use for generating TOTP Key.
	Rand io.Reader
}
