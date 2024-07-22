package main

import (
	"log"
	"os"
	mongodb "restaurant-service/inernal/mongoDB"
	rabbitmq "restaurant-service/inernal/rabbitMQ"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	conn, err := rabbitmq.RabbitMQConnection(os.Getenv("rabbit_url"))
	if err != nil {
		log.Fatal("Failed RabbitMQ Connection:", err)
	}
	defer conn.Close()
	channel, err := conn.Channel()
	if err != nil {
		log.Fatal("failed Channel in RabbitMQ :", err)
	}
	rabbitRepo := rabbitmq.NewRabbitRepo(channel)
	mongoRepo, err := mongodb.NewMongoRepo(os.Getenv("mongo_url"))
	if err != nil {
		log.Fatal("failed MongoDB Connection:", err)
	}

	rabbitRepo.ConsumeOrderForKitchen(mongoRepo)
}
