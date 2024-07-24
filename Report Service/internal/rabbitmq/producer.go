package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"rabbitmq-topic/internal/model"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (r *RabbitRepo) PublishReport(routingKey string, report model.ReportDetails, status string, id string) error {
	err := r.RabbitChannel.ExchangeDeclare(
		"report-exchange", // name
		"topic",           // type
		true,              // durable
		false,             // auto-deleted
		false,             // internal
		false,             // no-wait
		nil,               // arguments
	)
	if err != nil {
		log.Fatalf("%s: %s", "Failed to declare an exchange", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err :=json.Marshal(report)
	if err != nil {
		return fmt.Errorf("error: unable to marshal request body: %v",err)
	}
	err = r.RabbitChannel.PublishWithContext(ctx,
		"report-exchange", // exchange
		routingKey,        // routing key
		false,             // mandatory
		false,             // immediate
		amqp.Publishing{
			Headers: amqp.Table{
				"type" : status,
				"id" : id,
			},
			ContentType: "text/plain",
			Body:        body,
			Expiration: "10000",
		})
	if err != nil {
		log.Fatalf("%s: %s", "Failed to publish a message", err)
	}
	return nil
}
