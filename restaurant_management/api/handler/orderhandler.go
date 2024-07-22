package handler

import (
	"log"
	"restaurant-service/inernal/model"
	mongodb "restaurant-service/inernal/mongoDB"
	rabbitmq "restaurant-service/inernal/rabbitMQ"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	RabbitMQ *rabbitmq.RabbitRepo
	MongoDB  *mongodb.MongoRepo
}

func NewOrderHandler(r *rabbitmq.RabbitRepo, m *mongodb.MongoRepo) *OrderHandler {
	return &OrderHandler{RabbitMQ: r, MongoDB: m}
}
//////////////////////////////////////////////////////////////////////////////////////////////////////////
// @Router  				/order/create  [post]
// @Summary 			Creates Guest Orders  
// @Description 		This method is responsible for creating guest  orders and sends to kitchen  
// @Security 				BearerAuth
// @Tags					 ORDER
// @accept					json
// @Produce				  json
// @Param 					id    header    string    true    "Employee ID"
// @Param 					password    header    string    true    "Password"
// @Param 					body    body    model.GuestOrder    true  "Order Details"
// @Success					201 	{object}   string		"Order accepted . "
// @Failure					 400 {object}    error		"Incorrect Order Details !",
// @Failure					 500 {object}    error	"Unable to Send  an Order To kitchen !"
// @Failure					 403 {object}    error		"Unauthorized access"
func (o *OrderHandler) CreateNewOrder(c *gin.Context) {
	var guestOrder model.GuestOrder
	if err := c.BindJSON(&guestOrder); err != nil {
		log.Println("Failed encoding request body:", err)
		c.AbortWithStatusJSON(400, gin.H{
			"error": "Incorrect Order Details !",
		})
	}
	if err := o.RabbitMQ.PublishFullGuestOrder(&guestOrder); err != nil {
		log.Println("Failed Publishing guest order:", err)
		c.AbortWithStatusJSON(500, gin.H{
			"error": "Unable to Send  an Order To kitchen !",
		})
	}

	c.IndentedJSON(201, gin.H{
		"message": "Order accepted .",
	})

}
//////////////////////////////////////////////////////////////////////////////////////////////////////////
// @Router  				/order/fire  [post]
// @Summary 			Fires  Order course  
// @Description 		This method is responsible for firing guest  orders  courses  
// @Security 				BearerAuth
// @Tags					 ORDER
// @accept					json
// @Produce				  json
// @Param 					id    header    string    true    "Employee ID"
// @Param 					password    header    string    true    "Password"
// @Param 					body    body    model.FireCourse    true  "Fire Details"
// @Success					201 	{object}   string		"Fire Course accepted . "
// @Failure					 400 {object}    error		"Incorrect  FireCourse information !"
// @Failure					 500 {object}    error		"Unable to Send  fire course to kitchen !"
// @Failure					 403 {object}    error		"Unauthorized access"
func (o *OrderHandler) FireTheCourse(c *gin.Context) {
	var fireCourse model.FireCourse
	if err := c.BindJSON(&fireCourse); err != nil {
		log.Println("Failed encoding request body:", err)
		c.AbortWithStatusJSON(400, gin.H{
			"error": "Incorrect  FireCourse information !",
		})
	}
	if fireCourse.Course == "main"{
		fireCourse.Course = "kitchen"
	}
	err :=o.RabbitMQ.PublishForFire(fireCourse, fireCourse.Course)
	if err != nil {
		log.Println("Failed Publishing  fire course:", err)
		c.AbortWithStatusJSON(500, gin.H{
			"error": "Unable to Send  fire course to kitchen !",
		})
	}
	c.IndentedJSON(201, gin.H{
		"message": "Fire Course accepted.",
	})
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////
// @Router  				/order/fire  [get]
// @Summary 			Gets  Order Status  
// @Description 		This method is responsible for getting  guest  orders  status 
// @Security 				BearerAuth
// @Tags					 ORDER
// @accept					json
// @Produce				  json
// @Param 					id    header    string    true    "Employee ID"
// @Param 					password    header    string    true    "Password"
// @Param 					body    body    model.GetStatus    true  "Course Details"
// @Success 				200    {object}    model.StatusResponse     "Status Details"
// @Failure					 400 {object}    error		"Incorrect  Course information !",
// @Failure					 500 {object}    error	"Unable to get status information !"
// @Failure					 403 {object}    error		"Unauthorized access"
func (o *OrderHandler) GetOrderStatus(c *gin.Context) {
	var courseStatus  model.GetStatus
	if err := c.BindJSON(&courseStatus); err != nil {
		log.Println("Failed encoding request body:", err)
		c.AbortWithStatusJSON(400, gin.H{
			"error": "Incorrect Course information !",
		})
	}

	statusResponse, err :=o.MongoDB.GetOrdersStatusFromMongoDB(courseStatus)
	if err != nil {
		log.Println("Failed Getting  status response  from mongoDB:", err)
		c.AbortWithStatusJSON(500, gin.H{
			"error": "Unable to get status information !",
		})
	}
	c.IndentedJSON(200, statusResponse)
	
}
