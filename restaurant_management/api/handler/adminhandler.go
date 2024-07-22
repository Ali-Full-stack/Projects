package handler

import (
	"log"
	"restaurant-service/inernal/model"
	mongodb "restaurant-service/inernal/mongoDB"
	redisdb "restaurant-service/inernal/redisDB"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	MongoDB *mongodb.MongoRepo
	RedisClient  *redisdb.RedisClient
}

func NewAdminHandler(m *mongodb.MongoRepo, r *redisdb.RedisClient) *AdminHandler {
	return &AdminHandler{MongoDB: m, RedisClient: r}
}
//////////////////////////////////////////////////////////////////////////////////////////////////////////
// @Router  				/admin/register  [post]
// @Summary 			Registers New Employee 
// @Description 		This method is responsible for Registering   New Employees  
// @Security 				BearerAuth
// @Tags					 REGISTRATION
// @accept					json
// @Produce				  json
// @Param 					id    header    string    true    "Role"
// @Param 					password    header    string    true    "Password"
// @Param 					body    body    model.EmployeeInfo    true  "Employee Details"
// @Success                 201 {object}     model.EmployeeInfo    "New employee registered"
// @Failure					 400 {object}    error		"Incorrect Employee Information !"
// @Failure					 500 {object}    error		"Failed Registering New Employee !"
// @Failure					 403 {object}    error		"Unauthorized access"
func (a *AdminHandler) RegisterNewEmployee(c *gin.Context) {
	var employeeInfo model.EmployeeInfo
	if err := c.BindJSON(&employeeInfo); err != nil {
		log.Println("Failed encoding request body:", err)
		c.AbortWithStatusJSON(400, gin.H{
			"error": "Incorrect Employee Information !",
		})
	}
	resp, err := a.MongoDB.AddNewEmployeeIntoMongoDB(employeeInfo)
	if err != nil {
		log.Println("Failed Registering New Employee :", err)
		c.AbortWithStatusJSON(500, gin.H{
			"error": "Failed Registering New Employee !",
		})
	}
	err =a.RedisClient.AddEmployeeForLogin(resp.ID, employeeInfo)
	if err != nil {
		log.Println("Failed Adding  New Employee to Redis  :", err)
		c.AbortWithStatusJSON(500, gin.H{
			"error": "Unable to Add  New Employee For Login !",
		})
	}
	
	c.IndentedJSON(201, resp)
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////
// @Router  				/admin/update/:id  [put]
// @Summary 			Updates  Employee Information 
// @Description 		This method is responsible for Updating    Employee Information  
// @Security 				BearerAuth
// @Tags					 REGISTRATION
// @accept					json
// @Produce				  json
// @Param 					id    header    string    true    "Role"
// @Param 					password    header    string    true    "Password"
// @Param                   id                 path        string   true     "Employee ID"    
// @Param 					body    body    model.EmployeeInfo    true  "Employee Details"
// @Success					202    {object}     model.EmployeeResponse     "Response"
// @Failure					 400 {object}    error	"Incorrect Employee Information !"
// @Failure					 500 {object}   error		"Unable to Update employee info ! Try again later !"
// @Failure					 403 {object}   error	"Unauthorized access"
func (a *AdminHandler) UpdateEmployeeInformation(c *gin.Context) {
	id :=c.Param("id")
	var employeeInfo model.EmployeeInfo
	if err := c.BindJSON(&employeeInfo); err != nil {
		log.Println("Failed encoding request body:", err)
		c.AbortWithStatusJSON(400, gin.H{
			"error": "Incorrect Employee Information !",
		})
	}
	resp, err := a.MongoDB.UpdateEmployeeInfoInMongoDb(employeeInfo, id)
	if err != nil {
		log.Println("Failed Updating  Employee Info :", err)
		c.AbortWithStatusJSON(500, gin.H{
			"error": "Unable to Update employee info ! Try again later !",
		})
	}
	c.IndentedJSON(202, resp)
}
//////////////////////////////////////////////////////////////////////////////////////////////////////////
// @Router  				/admin/delete/:id  [delete]
// @Summary 			Deletes  Employee Information 
// @Description 		This method is responsible for Deleting    Employee Information  
// @Security 				BearerAuth
// @Tags					 REGISTRATION
// @accept					json
// @Produce				  json
// @Param 					id    header    string    true    "Role"
// @Param 					password    header    string    true    "Password"
// @Param                   id                 path        string   true     "Employee ID"    
// @Success				   202   {object}    model.EmployeeResponse      "Response"
// @Failure					 500 {object}    error		"Unable to Delete  Employee  ! Try again later !"
// @Failure					 403 {object}   error		"Unauthorized access"
func (a *AdminHandler) DeleteEmployee( c *gin.Context){
	id :=c.Param("id")

	resp, err :=a.MongoDB.DeleteEmployeeFromMongoDb(id)
	if err != nil {
		log.Println("Failed Deleting  Employee From mongoDB :", err)
		c.AbortWithStatusJSON(500, gin.H{
			"error": "Unable to Delete  Employee  ! Try again later !",
		})
	}
	c.IndentedJSON(202, resp)
}


