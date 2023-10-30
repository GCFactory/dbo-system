package amqp

import (
	"fmt"
	"github.com/GCFactory/dbo-system/platform/config"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

const (
	reconnectDelay = 5 * time.Second
	reInitDelay    = 2 * time.Second
	resendDelay    = 5 * time.Second
)

func NewAMQP(c *config.Config) (*amqp.Connection, error) {
	connectionURL := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		c.RabbitMQ.User,
		c.RabbitMQ.Password,
		c.RabbitMQ.Host,
		c.RabbitMQ.Port,
	)

	conn, err := amqp.Dial(connectionURL)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func CreateChannel(conn *amqp.Connection) (*amqp.Channel, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	return ch, nil
}
func DeclareQueue(ch *amqp.Channel, cfg *config.Config) (*amqp.Queue, error) {
	queue, err := ch.QueueDeclare(
		cfg.RabbitMQ.Queue,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}
	return &queue, nil
}
