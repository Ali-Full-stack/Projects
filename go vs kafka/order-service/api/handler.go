package api

import (
	"kafka-go/internal/model"
	"kafka-go/internal/mongodb"
	"kafka-go/kafka"
	"log"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	MongoRepo *mongodb.MongoRepo
	KafkaRepo *kafka.KafkaClient
}

func NewOrderHandler(m *mongodb.MongoRepo, k *kafka.KafkaClient) *OrderHandler {
	return &OrderHandler{MongoRepo: m, KafkaRepo: k}
}

func (o *OrderHandler) CreateOrder(c *gin.Context) {
	var order model.OrderInfo
	if err := c.BindJSON(&order); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(400, gin.H{
			"error": "invalid order information",
		})
		return
	}
	if err := o.KafkaRepo.CreateOrderKafka(c.Request.Context(), order); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{
			"error": "Unable to produce message",
		})
		return
	}

	c.IndentedJSON(201, gin.H{
		"message" : "request accepted ..",
	})
	
}

func (o *OrderHandler) UpdateOrder(c *gin.Context) {}

func (o *OrderHandler) DeleteOrder(c *gin.Context) {}

func (o *OrderHandler) GetAllOrders(c *gin.Context) {}

func (o *OrderHandler) GetOrderById(c *gin.Context) {}
