package grpc

import "errors"

var (
	ErrorEmptyStartEventList     = errors.New("Empty start event list")
	ErrorInvalidServer           = errors.New("Invalid server")
	ErrorInvalidServersOperation = errors.New("Invalid servers operation")
	ErrorInvalidOperationsData   = errors.New("Invalid operations data")
	ErrorInvalidServersTopic     = errors.New("Invalid servers topic")
	ErrorConvertTimestamp        = errors.New("Convert timestamp error")
	ErrorInvalidUUIDs            = errors.New("Invalid uuids")
)
