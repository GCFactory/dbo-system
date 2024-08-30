package account

import (
	acc_proto_api "github.com/GCFactory/dbo-system/service/account/gen_proto/proto/api"
	acc_proto_platform "github.com/GCFactory/dbo-system/service/account/gen_proto/proto/platform"
	"github.com/GCFactory/dbo-system/service/account/pkg/kafka"
	"golang.org/x/net/context"
)

type GRPCHandlers interface {
	ReserveAccount(ctx context.Context, saga_uuid string, event_uuid string, acc_data *acc_proto_platform.AccountDetails, kProducer *kafka.ProducerProvider) error
	ChangeAccountStatus(ctx context.Context, saga_uuid string, event_uuid string, operation_type string, acc_data *acc_proto_api.OperationDetails, kProducer *kafka.ProducerProvider) error
	GetAccountData(ctx context.Context, saga_uuid string, event_uuid string, acc_data *acc_proto_api.OperationDetails, kProducer *kafka.ProducerProvider) error
	OperationWithAccAmount(ctx context.Context, saga_uuid string, event_uuid string, operation_type string, acc_data *acc_proto_api.OperationDetails, kProducer *kafka.ProducerProvider) error
}
