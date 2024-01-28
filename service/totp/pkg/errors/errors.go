package errors

import "errors"

var (
	ActiveTotp       = errors.New("User's totp is active")
	EmptyTotpUrl     = errors.New("Empty totp url")
	NoSecretField    = errors.New("No secret field")
	NoPeriodField    = errors.New("No period field")
	NoIssuerField    = errors.New("No issure field")
	NoDigitsField    = errors.New("No digits field")
	NoAlgorithmField = errors.New("No algorithm field")
)
