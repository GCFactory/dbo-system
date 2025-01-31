package repository

import "errors"

var (
	ErrorCreateConfig               = errors.New("totpRepo.CreateConfig.QueryRowxContext")
	ErrorGetActiveConfig            = errors.New("totpRepo.GetActiveConfig.QueryRowContext")
	ErrorGetConfigByTotpId          = errors.New("totpRepo.GetConfigByTotpId.QueryRowContext")
	ErrorUpdateTotpActivityByTotpId = errors.New("totpRepo.GetConfigByTotpId.QueryRowContext")
	ErrorGetConfiByUserId           = errors.New("totpRepo.GetConfigByUserId.QueryRowContext")
	ErrorGetLastDisabledConfig      = errors.New("totpRepo.GetLastDisabledConfig.QueryRowContext")
)
