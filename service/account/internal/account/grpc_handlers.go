package account

import (
	acc_proto_platform "github.com/GCFactory/dbo-system/service/account/gen_proto/proto/platform"
	"github.com/GCFactory/dbo-system/service/account/pkg/kafka"
	"golang.org/x/net/context"
)

type GRPCHandlers interface {
	ReserveAccount(ctx context.Context, saga_uuid string, acc_data *acc_proto_platform.AccountDetails, kProducer *kafka.ProducerProvider) error
}
