package usecase

import "errors"

var (
	ErrorNoUserIdHeader          = errors.New("No user id into rmq message headers")
	ErrorNoNotificationLvlHeader = errors.New("No notification lvl into rmq message headers")
	ErrorInvalidNotificationLvl  = errors.New("Unknown notification lvl")
)
