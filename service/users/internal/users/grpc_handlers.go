package users

import (
	"github.com/GCFactory/dbo-system/service/users/gen_proto/proto/api"
	"github.com/GCFactory/dbo-system/service/users/pkg/kafka"
	"golang.org/x/net/context"
)

type GRPCHandlers interface {
	AddUser(ctx context.Context, saga_uuid string, event_uuid string, users_data *api.UserInfo, kProducer *kafka.ProducerProvider) error
	GetUserData(ctx context.Context, saga_uuid string, event_uuid string, operation_details *api.OperationDetails, kProducer *kafka.ProducerProvider) error
	UpdateUsersPassport(ctx context.Context, saga_uuid string, event_uuid string, operation_details *api.OperationDetails, kProducer *kafka.ProducerProvider) error
	AddUserAccount(ctx context.Context, saga_uuid string, event_uuid string, operation_details *api.OperationDetails, kProducer *kafka.ProducerProvider) error
	GetUsersAccounts(ctx context.Context, saga_uuid string, event_uuid string, operation_details *api.OperationDetails, kProducer *kafka.ProducerProvider) error
}
