package notification

import (
	api "github.com/GCFactory/dbo-system/service/notification/gen_proto/proto/notification_api"
	"github.com/GCFactory/dbo-system/service/notification/pkg/kafka"
	"golang.org/x/net/context"
)

type GRPCHandlers interface {
	AddUserSettings(ctx context.Context, saga_uuid string, event_uuid string, users_data *api.AdditionalInfo, kProducer *kafka.ProducerProvider) error
	RemoveUserSettings(ctx context.Context, saga_uuid string, event_uuid string, users_data *api.AdditionalInfo, kProducer *kafka.ProducerProvider) error
}
