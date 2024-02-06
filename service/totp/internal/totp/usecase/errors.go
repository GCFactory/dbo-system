package usecase

import (
	"github.com/pkg/errors"
)

var (
	ErrorUpdateActivityByTotpId = errors.New("totpUC.totpRepo.UpdateTotpActivityByTotpId")
	ErrorGenTotp                = errors.New("totpUC.totpPkg.Generate")
	ErrorCreateConfig           = errors.New("totpUC.totpRepo.CreateConfig")
)
