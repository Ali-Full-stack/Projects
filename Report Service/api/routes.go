package api

import (
	"log"
	"os"
	"rabbitmq-topic/api/middleware"
	_ "rabbitmq-topic/api/swagger/docs"
	"rabbitmq-topic/internal/mongodb"
	"rabbitmq-topic/internal/rabbitmq"
	"rabbitmq-topic/internal/redisdb"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	swaggerFiles "github.com/swaggo/files"
	swag "github.com/swaggo/gin-swagger"
)

// New ...
// @title  Project: Report Service
// @description This swagger UI was created to manage Reports
// @version 1.0
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @host localhost:8888
// @contact.email  ali.team@gmail.com
func Routes() {
	r := gin.Default()

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

	redisRepo := redisdb.ConnectRedis()
	redisLogger := log.New(redisRepo, "", log.LstdFlags)

	reportHandler := NewReportHandler(mongoRepo, rabbitRepo, redisRepo, redisLogger)

	report :=r.Group("/reports")
	report.Use(middleware.VerifyTokenMIddleware)
	{
			report.POST("", reportHandler.CreateNewReport)
			report.PUT("/:id", reportHandler.UpdateExistingReport)
			report.DELETE("/:id", reportHandler.DeleteExistingReport)
			report.GET("", reportHandler.GetAllReports)
			report.GET("/search", reportHandler.GetReportsByFilter)
			report.GET("/status/:id", reportHandler.GetReportsStatusById)
	}

	r.POST("/employees/register", reportHandler.RegisterNewReporters)
	r.GET("/employees/login",  reportHandler.LoginReporters)

	r.GET("/swagger/*any", swag.WrapHandler(swaggerFiles.Handler))

	r.Run(os.Getenv("server_url"))

}
