package repository

const (
	createSaga = `INSERT INTO registration_service.saga_table (saga_uuid, sage_status, saga_type, saga_name, list_of_events)
					VALUES ($1, $2, $3, $4, $5)`
)
