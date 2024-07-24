package rabbitmq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func RabbitMQConnection(rabbit_url string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(rabbit_url)
	if err != nil {
		return nil, fmt.Errorf("error: failed rabbitMQ connection: %v", err)
	}
	return conn, nil
}

type RabbitRepo struct {
	RabbitChannel *amqp.Channel
}

func NewRabbitRepo(rch *amqp.Channel) *RabbitRepo {
	return &RabbitRepo{RabbitChannel: rch}
}