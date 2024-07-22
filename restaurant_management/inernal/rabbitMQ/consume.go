package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"
	"restaurant-service/inernal/model"
	mongodb "restaurant-service/inernal/mongoDB"
	"time"
)

func (r *RabbitRepo) ConsumeOrderForPizza(m *mongodb.MongoRepo) {
	qName, err := r.DeclareQueue("pizza")
	if err != nil {
		log.Fatal(err)
	}
	msgs, err := r.RabbitChannel.Consume(
		qName, // queue
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		log.Fatalf("failed to register a consumer for pizza: %v", err)
	}
	forever := make(chan bool)

	go func() {
		for msg := range msgs {
			orderType, _ := msg.Headers["type"].(string)
			tableNumber, _ := msg.Headers["table"].(int32)
			switch orderType {
			case "order":
				var pizzaOrder model.Pizza
				if err := json.Unmarshal(msg.Body, &pizzaOrder); err != nil {
					log.Fatal("failed to unmarshal message body in pizza order:", err)
				}
				m.UpdateOrdersStatusInMongoDB(tableNumber, "pizza", "Pending")
				fmt.Println("Table Number:", tableNumber)
				fmt.Println("Pizzas :", pizzaOrder.Pizza)
			case "fire":
				var fireCourse model.FireCourse
				if err := json.Unmarshal(msg.Body, &fireCourse); err != nil {
					log.Fatal("failed to unmarshal message body in pizza firecoures")
				}
				m.UpdateOrdersStatusInMongoDB(fireCourse.TableNumber, fireCourse.Course, "Preparing")
				fmt.Println("Table Number :", fireCourse.TableNumber)
				fmt.Println("Waiter :", fireCourse.Waiter)
				fmt.Println("Course:",fireCourse.Course)
				fmt.Println("Time :", fireCourse.Time)
				go func(){
					time.Sleep(1 * time.Minute)
					m.UpdateOrdersStatusInMongoDB(fireCourse.TableNumber, fireCourse.Course, "Ready")
				}()				
			}

			log.Println("Received message for Pizza:", orderType)
		}
	}()
	fmt.Println("Pizza Waiting for Orders >>  ")
	<-forever
}

func (r *RabbitRepo) ConsumeOrderForBar(m *mongodb.MongoRepo) {
	qName, err := r.DeclareQueue("bar")
	if err != nil {
		log.Fatal(err)
	}
	msgs, err := r.RabbitChannel.Consume(
		qName, // queue
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		log.Fatalf("failed to register a consumer for bar: %v", err)
	}
	forever := make(chan bool)

	go func() {
		for msg := range msgs {
			orderType, _ := msg.Headers["type"].(string)
			tableNumber, _ := msg.Headers["table"].(int32)

			switch orderType {
			case "order":
				var drinkOrder model.Drink
				if err := json.Unmarshal(msg.Body, &drinkOrder); err != nil {
					log.Fatal("failed to unmarshal message body in drink order:", err)
				}
				m.UpdateOrdersStatusInMongoDB(tableNumber, "drink", "Pending")
				fmt.Println("Table Number:", tableNumber)
				fmt.Println("Drinks:", drinkOrder.Drink)
			case "fire":
				var fireCourse model.FireCourse
				if err := json.Unmarshal(msg.Body, &fireCourse); err != nil {
					log.Fatal("failed to unmarshal message body in drink firecoures")
				}
				m.UpdateOrdersStatusInMongoDB(fireCourse.TableNumber, fireCourse.Course, "Preparing")
				fmt.Println("Table Number :", fireCourse.TableNumber)
				fmt.Println("Waiter :", fireCourse.Waiter)
				fmt.Println("Course:",fireCourse.Course)
				fmt.Println("Time :", fireCourse.Time)
				go func(){
					time.Sleep(1 * time.Minute)
					m.UpdateOrdersStatusInMongoDB(fireCourse.TableNumber, fireCourse.Course, "Ready")
				}()	
			}

			log.Println("Received message for Bar:", orderType)
		}
	}()
	fmt.Println("Bar  Waiting for Orders >>  ")
	<-forever
}

func (r *RabbitRepo) ConsumeOrderForDessert(m *mongodb.MongoRepo) {
	qName, err := r.DeclareQueue("dessert")
	if err != nil {
		log.Fatal(err)
	}
	msgs, err := r.RabbitChannel.Consume(
		qName, // queue
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		log.Fatalf("failed to register a consumer for dessert: %v", err)
	}
	forever := make(chan bool)

	go func() {
		for msg := range msgs {
			orderType, _ := msg.Headers["type"].(string)
			tableNumber, _ := msg.Headers["table"].(int32)
			
			switch orderType {
			case "order":
				var dessertOrder model.Dessert
				if err := json.Unmarshal(msg.Body, &dessertOrder); err != nil {
					log.Fatal("failed to unmarshal message body in dessert order:", err)
				}
				m.UpdateOrdersStatusInMongoDB(tableNumber, "dessert", "Pending")
				fmt.Println("Table Number:", tableNumber)
				fmt.Println("Desserts :", dessertOrder.Dessert)
				fmt.Println("Coffees :", dessertOrder.Coffee)
			case "fire":
				var fireCourse model.FireCourse
				if err := json.Unmarshal(msg.Body, &fireCourse); err != nil {
					log.Fatal("failed to unmarshal message body in dessert firecoures")
				}
				m.UpdateOrdersStatusInMongoDB(fireCourse.TableNumber, fireCourse.Course, "Preparing")
				fmt.Println("Table Number :", fireCourse.TableNumber)
				fmt.Println("Waiter :", fireCourse.Waiter)
				fmt.Println("Course:",fireCourse.Course)
				fmt.Println("Time :", fireCourse.Time)
				go func(){
					time.Sleep(1 * time.Minute)
					m.UpdateOrdersStatusInMongoDB(fireCourse.TableNumber, fireCourse.Course, "Ready")
				}()	
			}

			log.Println("Received message for Dessert:", orderType)
		}
	}()
	fmt.Println("Dessert  Waiting for Orders >>  ")
	<-forever
}

func (r *RabbitRepo) ConsumeOrderForKitchen(m *mongodb.MongoRepo) {
	qName, err := r.DeclareQueue("kitchen")
	if err != nil {
		log.Fatal(err)
	}
	msgs, err := r.RabbitChannel.Consume(
		qName, // queue
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		log.Fatalf("failed to register a consumer for kitchen: %v", err)
	}
	forever := make(chan bool)

	go func() {
		for msg := range msgs {
			orderType, _ := msg.Headers["type"].(string)
			
			switch orderType {
			case "order":
				var kitchenOrder model.GuestOrder
				if err := json.Unmarshal(msg.Body, &kitchenOrder); err != nil {
					log.Fatal("failed to unmarshal message body in Kitchen order:", err)
				}
				m.AddNewOrderIntoMongoDB(&kitchenOrder)
			case "fire":
				var fireCourse model.FireCourse
				if err := json.Unmarshal(msg.Body, &fireCourse); err != nil {
					log.Fatal("failed to unmarshal message body in kitchen firecoures")
				}
				fireCourse.Course = "main"
				m.UpdateOrdersStatusInMongoDB(fireCourse.TableNumber, fireCourse.Course, "Preparing")
				fmt.Println("Table Number :", fireCourse.TableNumber)
				fmt.Println("Waiter :", fireCourse.Waiter)
				fmt.Println("Course:",fireCourse.Course)
				fmt.Println("Time :", fireCourse.Time)
				go func(){
					time.Sleep(1 * time.Minute)
					m.UpdateOrdersStatusInMongoDB(fireCourse.TableNumber, fireCourse.Course, "Ready")
				}()	
			}

			log.Println("Received message for Kitchen:", orderType)
		}
	}()
	fmt.Println("Kitchen  Waiting for Orders >>  ")
	<-forever
}
