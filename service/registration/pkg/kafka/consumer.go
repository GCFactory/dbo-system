package kafka

import "github.com/IBM/sarama"

type KafkaConsumerGroup struct {
	config   *sarama.Config
	Consumer sarama.ConsumerGroup
}

func NewKafkaConsumer(brokerlist []string, groupid string) (*KafkaConsumerGroup, error) {
	cg, err := sarama.NewConsumerGroup(brokerlist, groupid, func() *sarama.Config {
		consumer := sarama.NewConfig()
		consumer.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
		consumer.Consumer.Return.Errors = false
		return consumer
	}())
	if err != nil {
		return nil, err
	}

	return &KafkaConsumerGroup{
		Consumer: cg,
	}, nil
}
