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
func (r *RabbitRepo) ExchangeDeclaration() error {
	err := r.RabbitChannel.ExchangeDeclare(
		"order-exchange", // name
		"direct",         // type
		false,            // durable
		false,            // auto-delete
		false,            // internal
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare an order - exchange: %v", err)
	}
	return nil
}

func (r *RabbitRepo) DeclareQueue(keyName string)(string, error){
	if err := r.ExchangeDeclaration(); err != nil {
		return "", err
	}
	q, err := r.RabbitChannel.QueueDeclare(
		keyName, // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		return "", fmt.Errorf("failed to declare a queue: %v", err)
	}

	if err := r.RabbitChannel.QueueBind(
		q.Name,           // queue name
		keyName,         // routing key
		"order-exchange", // exchange
		false,
		nil,
	); err != nil {
		return "", fmt.Errorf("failed to bind a queue: %v", err)
	}
	return q.Name, nil
}

	
