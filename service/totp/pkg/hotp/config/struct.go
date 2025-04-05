package config

import "github.com/GCFactory/dbo-system/service/totp/pkg/otp/config"

const Debug = false

// ValidateOpts provides options for ValidateCustom().
type ValidateOpts struct {
	// Digits as part of the input. Defaults to 6.
	Digits config.Digits
	// Algorithm to use for HMAC. Defaults to SHA1.
	Algorithm config.Algorithm
}
