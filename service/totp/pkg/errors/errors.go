package errors

import "errors"

var (
	ActiveTotp       = errors.New("User's totp is active")
	NoSecretField    = errors.New("No secret field")
	NoPeriodField    = errors.New("No period field")
	NoIssuerField    = errors.New("No issure field")
	NoDigitsField    = errors.New("No digits field")
	NoAlgorithmField = errors.New("No algorithm field")
	NoUserId         = errors.New("No totp for this user id")
	WrongTotpCode    = errors.New("Wrong totp code")
	NoTotpId         = errors.New("No totp with this totp id")
	TotpIsDisabled   = errors.New("Totp is disabled yet")
	NoId             = errors.New("No such user and totp ids")
	TotpIsActive     = errors.New("Totp is active yet")
	NoUserName       = errors.New("No user name")
)
