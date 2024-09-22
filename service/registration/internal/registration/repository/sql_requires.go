package repository

const (
	CreateSaga = `INSERT INTO saga (
                  saga_uuid, 
                  saga_status, 
                  saga_type, 
                  saga_name)
				VALUES ($1, $2, $3, $4, $5)`
	DeleteSaga = `DELETE FROM saga 
       			WHERE saga_uuid=$1`
	GetSaga = `SELECT * FROM saga 
         		WHERE saga_uuid=$1`
	UpdateSaga = `UPDATE saga 
				SET saga_status = $2,
				    saga_type = $3,
				    saga_name = $4
				WHERE saga_uuid = $1`
	CreateSagaConnection = `INSERT INTO saga_connection (
                        	current_saga_uuid,
							next_saga_uuid,
                             acc_connection_status)
							Values ($1, $2, $3)`
	GetSagaConnectionsCurrentSaga = `SELECT * as list_of_connections
						FROM saga_connection
						WHERE current_saga_uuid = $1`
	GetSagaConnectionsNextSaga = `SELECT * as list_of_connections
						FROM saga_connection
						WHERE next_saga_uuid = $1`
	DeleteSagaConnection = `DELETE FROM saga_connection
							WHERE current_saga_uuid = $1
								AND
							next_saga_uuid = $2`
	UpdateSagaConnection = `UPDATE saga_connection
							SET
								acc_connection_status = $3
							WHERE current_saga_uuid = $1
								AND	
							next_saga_uuid = $2`
	CreateEvent = `INSERT INTO event (
	   			event_uuid,
				saga_uuid,
				event_status,
				event_name,
				event_result,
				event_rollback_uuid)
				VALUES ($1, $2, $3, $4, $5, $6)`
	GetEvent = `SELECT * 
				FROM event
				WHERE event_uuid = $1`
	DeleteEvent = `DELETE FROM event
				WHERE event_uuid = $1`
	UpdateEvent = `UPDATE event
					SET event_status = $2,
					    event_name = $3,
					    event_result = $4,
					    event_rollback_uuid = $5
					WHERE event_uuid = $1`
	GetListOfSagaEvents = `SELECT event_uuid as list_of_events
						FROM event
						WHERE saga_uuid = $1`
)
