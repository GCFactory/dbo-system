package registration

import (
	"context"
)

type RegistrationGRPCHandlers interface {
	StartOperation(ctx context.Context, operation_type uint8, operation_data []byte) error
	Process(ctx context.Context, event_type string, event_data []byte) error
}
