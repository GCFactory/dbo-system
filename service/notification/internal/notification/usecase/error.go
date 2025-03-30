package usecase

import "errors"

var (
	ErrorUserNotFound             = errors.New("User not found")
	ErrorUserSettingsAlreadyExist = errors.New("User settings already exist")
	ErrorNoUserIdHeader           = errors.New("No user id into rmq message headers")
	ErrorNoNotificationLvlHeader  = errors.New("No notification lvl into rmq message headers")
	ErrorInvalidNotificationLvl   = errors.New("Unknown notification lvl")
)
