package repository

import "errors"

var (
	ErrorCreateSaga       = errors.New("registrationRepo.CreateSaga.QueryRowxContext")
	ErrorGetSaga          = errors.New("registrationRepo.GetSagaById.QueryRowxContext")
	ErrorUpdateSagaStatus = errors.New("registrationRepo.UpdateSagaStatus.ExecContext")
	ErrorUpdateSagaEvents = errors.New("registrationRepo.UpdateSagaEvents.ExecContext")
)
