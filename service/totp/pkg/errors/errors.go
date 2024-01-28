package errors

import "errors"

var (
	ActiveTotp       = errors.New("User's totp is active")
	NoSecretField    = errors.New("No secret field")
	NoPeriodField    = errors.New("No period field")
	NoIssuerField    = errors.New("No issure field")
	NoDigitsField    = errors.New("No digits field")
	NoAlgorithmField = errors.New("No algorithm field")
	NoUserTotp       = errors.New("No totp for this user id")
	WrongTotpCode    = errors.New("Wrong totp code")
)
