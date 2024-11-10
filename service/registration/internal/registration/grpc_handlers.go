package registration

import (
	"context"
)

type RegistrationGRPCHandlers interface {
	StartOperation(ctx context.Context, operation_type uint8, start_data map[string]interface{}) error
	Process(ctx context.Context, event_type string, event_data []byte) error
}
