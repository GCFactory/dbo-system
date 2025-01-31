package kafka

import (
	"fmt"
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

// Consumer represents a Sarama consumer group consumer
type Consumer struct {
	Ready       chan bool
	HandlerFunc func(*sarama.ConsumerMessage) error
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(consumer.Ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
// Once the Messages() channel is closed, the Handler must finish its processing
// loop and exit.
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/IBM/sarama/blob/main/consumer_group.go#L27-L29
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				fmt.Printf("message channel was closed")
				return nil
			}

			err := consumer.HandlerFunc(message)
			if err != nil {
				return err
			}

			// Mark message as consumed
			session.MarkMessage(message, "")

		// Should return when `session.Context()` is done.
		// If not, will raise `ErrRebalanceInProgress` or `read tcp <ip>:<port>: i/o timeout` when kafka rebalance. see:
		// https://github.com/IBM/sarama/issues/1192
		case <-session.Context().Done():
			return nil
		}
	}
}
