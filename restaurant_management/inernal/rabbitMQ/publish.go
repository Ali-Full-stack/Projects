package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"
	"restaurant-service/inernal/model"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (r *RabbitRepo) PublishFullGuestOrder(order *model.GuestOrder) error {
	err := r.ExchangeDeclaration()
	if err != nil {
		return err
	}

	r.PublishToOne("kitchen", order, order.TableNumber)
	r.PublishToOne("bar", order.Drink, order.TableNumber)
	r.PublishToOne("pizza", order.Pizza, order.TableNumber)
	r.PublishToOne("dessert", order.Dessert, order.TableNumber)
	return nil
}

func (r *RabbitRepo) PublishToOne(keyName string, order any, tableNumber int32) {

	data, err := json.Marshal(order)
	if err != nil {
		log.Println("failed marshalling while publishing to kitchen:", err)
	}
	if err := r.RabbitChannel.Publish(
		"order-exchange", // exchange
		keyName,          // routing key
		false,            // mandatory
		false,            // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Headers:     amqp.Table{"type": "order", "table": tableNumber},
			Body:        data,
		}); err != nil {
		log.Println("failed to publish a message:", err)
	}
}

func (r *RabbitRepo) PublishForFire(fire model.FireCourse, keyName string) error {
	err := r.ExchangeDeclaration()
	if err != nil {
		return fmt.Errorf("failed to publish fire course: %v", err)
	}
	fire.Time = time.Now().Format(time.ANSIC)
	data, err := json.Marshal(fire)
	if err != nil {
		return fmt.Errorf("failed marshalling while publishing to kitchen: %v", err)
	}
	if err := r.RabbitChannel.Publish(
		"order-exchange", // exchange
		keyName,          // routing key
		false,            // mandatory
		false,            // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Headers:     amqp.Table{"type": "fire", "table": fire.TableNumber},
			Body:        data,
		}); err != nil {
		return fmt.Errorf("failed to publish a message in FireCourse : %v", err)
	}
	return nil
}
