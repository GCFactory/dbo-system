package kafka

import (
	"github.com/IBM/sarama"
)

type ConsumerGroup struct {
	Consumer sarama.ConsumerGroup
}

func NewKafkaConsumer(brokerlist []string, groupid string) (*ConsumerGroup, error) {
	cg, err := sarama.NewConsumerGroup(brokerlist, groupid, func() *sarama.Config {
		consumer := sarama.NewConfig()
		consumer.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
		consumer.Consumer.Return.Errors = false
		return consumer
	}())
	if err != nil {
		return nil, err
	}

	return &ConsumerGroup{
		Consumer: cg,
	}, nil
}
