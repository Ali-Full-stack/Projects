package api

import (
	"log"
	"os"
	"restaurant-service/api/handler"
	"restaurant-service/api/middleware"
	mongodb "restaurant-service/inernal/mongoDB"
	rabbitmq "restaurant-service/inernal/rabbitMQ"
	redisdb "restaurant-service/inernal/redisDB"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	_"restaurant-service/api/swagger/docs"
	swaggerFiles "github.com/swaggo/files"
	swag "github.com/swaggo/gin-swagger"
)
// New ...
// @title  Project: RESTAURANT MANAGEMENT
// @description This swagger UI was created to manage Restaurant operation
// @version 1.0
// @host localhost:8888
// @contact.email  ali.team@gmail.com

func Routes() {
	r := gin.Default()
	/////////////////////////////////////////////////////////////////////////////////////////////////////////
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
	/////////////////////////////////////////////////////////////////////////////////////////////////////////
	mongoRepo, err := mongodb.NewMongoRepo(os.Getenv("mongo_url"))
	if err != nil {
		log.Fatal("failed MongoDB Connection:", err)
	}
	/////////////////////////////////////////////////////////////////////////////////////////////////////////
	redisClient := redisdb.ConnectRedis()
	/////////////////////////////////////////////////////////////////////////////////////////////////////////

	orderHandler := handler.NewOrderHandler(rabbitRepo, mongoRepo)

	order := r.Group("/order")
	order.Use(middleware.EmployeePasswordMiddleware())
	{
		order.POST("/create", orderHandler.CreateNewOrder)
		order.POST("fire", orderHandler.FireTheCourse)
		order.GET("/status", orderHandler.GetOrderStatus)
	}
	/////////////////////////////////////////////////////////////////////////////////////////////////////////
	adminhandler := handler.NewAdminHandler(mongoRepo, redisClient)

	admin := r.Group("/admin")
	admin.Use(middleware.AdminPasswordMiddleware())
	{
		admin.POST("/register", adminhandler.RegisterNewEmployee)
		admin.PUT("/update/:id", adminhandler.UpdateEmployeeInformation)
		admin.DELETE("/delete/:id", adminhandler.DeleteEmployee)
	}
	/////////////////////////////////////////////////////////////////////////////////////////////////////////
	r.GET("/swagger/*any", swag.WrapHandler(swaggerFiles.Handler))
	////////////////////////////////////////////////////////////////////////////////////////////////////

	r.Run(os.Getenv("server_url"))
}
