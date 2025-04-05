package grpc

import (
	"errors"
	"github.com/GCFactory/dbo-system/service/notification/internal/notification/repo"
	"github.com/GCFactory/dbo-system/service/notification/internal/notification/usecase"
)

var (
	ErrorNoUserEmail          = errors.New("No user email")
	ErrorUnknownOperationType = errors.New("unknown operation type")
	ErrorInvalidInputData     = errors.New("invalid input data")
)

var GRPCErrors = map[error]uint32{

	repo.ErrorNoUserSettingsFound:         001,
	usecase.ErrorUserNotFound:             101,
	usecase.ErrorUserSettingsAlreadyExist: 110,
	usecase.ErrorNoUserIdHeader:           120,
	usecase.ErrorNoNotificationLvlHeader:  130,
	usecase.ErrorInvalidNotificationLvl:   140,
	ErrorNoUserEmail:                      201,
	ErrorUnknownOperationType:             210,
	ErrorInvalidInputData:                 220,
}

func GetErrorCode(err error) uint32 {
	var result uint32
	result, ok := GRPCErrors[err]

	if !ok {
		result = 4294967295
	}

	return result
}
