package errors

import "errors"

var (
	ActiveTotp = errors.New("User's totp is active")
)
