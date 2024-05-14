package repository

const (
	CreateSaga = `INSERT INTO registration_service.saga_table (saga_uuid, saga_status, saga_type, saga_name, list_of_events)
					VALUES ($1, $2, $3, $4, $5) 
					RETURNING *`
	GetSaga          = `SELECT * FROM registration_service.saga_table WHERE saga_uuid=$1 LIMIT 1`
	UpdateSagaStatus = `UPDATE registration_service.saga_table SET sage_status = $2 WHERE saga_uuid = $1`
	UpdateSagaEvents = `UPDATE registration_service.saga_table SET list_of_events = $2 WHERE saga_uuid = $1`
)
