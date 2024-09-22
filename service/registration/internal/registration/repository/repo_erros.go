package repository

import "errors"

var (
	ErrorCreateSaga                = errors.New("registrationRepo.CreateSaga")
	ErrorDeleteSaga                = errors.New("registrationRepo.DeleteSaga")
	ErrorGetSaga                   = errors.New("registrationRepo.GetSaga")
	ErrorUpdateSaga                = errors.New("registrationRepo.UpdateSaga.ExecContext")
	ErrorUpdateSagaConnection      = errors.New("registrationRepo.UpdateSagaConnection.ExecContext")
	ErrorCreateSagaConnection      = errors.New("registrationRepo.CreateSagaConnection")
	ErrorGetSagaCurrentConnections = errors.New("registrationRepo.GetSagaCurrentConnections")
	ErrorGetSagaNextConnections    = errors.New("registrationRepo.GetSagaNextConnections")
	ErrorDeleteSagaConnection      = errors.New("registrationRepo.DeleteSagaConnection")
	ErrorCreateEvent               = errors.New("registrationRepo.CreateEvent")
	ErrorGetEvent                  = errors.New("registrationRepo.GetEvent")
	ErrorDeleteEvent               = errors.New("registrationRepo.DeleteEvent")
	ErrorUpdateEvent               = errors.New("registrationRepo.UpdateEvent")
	ErrorGetListOfSagaEvents       = errors.New("registrationRepo.GetListOfSagaEvents")
)
