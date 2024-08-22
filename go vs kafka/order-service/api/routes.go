package api

import (
	"kafka-go/internal/mongodb"
	"kafka-go/kafka"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

func Routes() {
	r := gin.Default()

	mongoRepo, err := mongodb.NewMongoRepo(os.Getenv("mongo_url"))
	if err != nil {
		log.Fatal("failed MongoDB Connection:", err)
	}

	kafkaRepo, err := kafka.ConnectKafka(os.Getenv("kafka_url"), "order-events", 3, 1)
	if err != nil {
		log.Fatal(err)
	}
	defer kafkaRepo.Client.Close()

	orderHandler := NewOrderHandler(mongoRepo, kafkaRepo)

	r.POST("/orders", orderHandler.CreateOrder)
	r.PUT("/orders/:id", orderHandler.UpdateOrder)
	r.DELETE("/orders/:id", orderHandler.DeleteOrder)
	r.GET("/orders", orderHandler.GetAllOrders)
	r.GET("/orders/:id", orderHandler.GetOrderById)

	r.Run(os.Getenv("server_url"))

}
