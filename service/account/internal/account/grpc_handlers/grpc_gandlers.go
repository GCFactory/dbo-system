package grpc_handlers

import (
	"fmt"
	acc_proto_platform "github.com/GCFactory/dbo-system/service/account/gen_proto/proto/platform"
	"github.com/GCFactory/dbo-system/service/account/internal/account"
	"github.com/GCFactory/dbo-system/service/account/pkg/kafka"
	"golang.org/x/net/context"
)

type AccountGRPCHandlers struct {
	kProducer *kafka.ProducerProvider
}

func (accH AccountGRPCHandlers) ReserveAccount(ctx context.Context, acc_data *acc_proto_platform.AccountDetails, kProducer *kafka.ProducerProvider) error {
	//	TODO: реализовать
	fmt.Printf("Message claimed: value = %s\n", acc_data)
	return nil
}

func NewAccountGRPCHandlers(kProducer *kafka.ProducerProvider) account.GRPCHandlers {
	return &AccountGRPCHandlers{kProducer: kProducer}
}
