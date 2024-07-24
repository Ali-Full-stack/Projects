package api

import (
	"log"
	"net/http"
	"rabbitmq-topic/auth"
	"rabbitmq-topic/internal/model"
	"rabbitmq-topic/internal/mongodb"
	"rabbitmq-topic/internal/rabbitmq"
	"rabbitmq-topic/internal/redisdb"

	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	MongoRepo *mongodb.MongoRepo
	RabbitRepo  *rabbitmq.RabbitRepo
	RedisClient   *redisdb.RedisClient
	RedisLog   *log.Logger
}

func NewReportHandler(m *mongodb.MongoRepo, rm *rabbitmq.RabbitRepo, rc *redisdb.RedisClient, rl *log.Logger )*ReportHandler{
	return &ReportHandler{MongoRepo: m, RabbitRepo: rm, RedisClient: rc,RedisLog: rl}
}
// @Router  				/reports  [post]
// @Summary 			Creates  Reports  
// @Description 		This method is responsible for creating  reports   
// @Security 				BearerAuth
// @Tags					 REPORT
// @accept					json
// @Produce				  json
// @Param 					token    header    string    true    "Reporter Token"
// @Param 					body    body    model.ReportDetails    true  "Report Details"
// @Success					201 	{object}   string		"Request Accepted .."
// @Failure					 400 {object}    error		"Invalid report details !"
// @Failure					 500 {object}    error	"Unable to publish report details !"
// @Failure					 403 {object}    error		"Permission Denied !"
func (r *ReportHandler) CreateNewReport(c *gin.Context){
	r.RedisLog.Println("INFO: received http request on CreatingNewReport. ")
	var reportDetails model.ReportDetails
	if err :=c.BindJSON(&reportDetails); err != nil {
		r.RedisLog.Println("ERROR: failed to decode request body on CreatingNewReport :",err)
		c.AbortWithStatusJSON(400, gin.H{
			"error" : "Invalid report details !",
		})
		return
	}
	routingKey:=reportDetails.Title+".created"
	if err :=r.RabbitRepo.PublishReport(routingKey, reportDetails, "created", ""); err != nil {
		r.RedisLog.Println("ERROR: failed to publish Report details on CreatingNewReport :",err)
		c.AbortWithStatusJSON(500, gin.H{
			"error" : "Unable to publish report details !",
		})
		return
	}
	r.RedisLog.Println("INFO: Report Published Succesfully with Title", reportDetails.Title," on CreatingNewReport ..")
	c.IndentedJSON(201, model.ReportResponse{Message: "Request Accepted .."})
}
// @Router  				/reports/:id  [put]
// @Summary 			Updates  Reports  
// @Description 		This method is responsible for updates  reports   
// @Security 				BearerAuth
// @Tags					 REPORT
// @accept					json
// @Produce				  json
// @Param 					token    header    string    true    "Reporter Token"
// @Param 					id    		path    string    true    "Report ID"
// @Param 					body    body    model.ReportDetails    true  "Report Details"
// @Success					202 	{object}   string		"Request Accepted .."
// @Failure					 400 {object}    error		"Invalid report details !"
// @Failure					 500 {object}    error	"Unable to publish report details !"
// @Failure					 403 {object}    error		"Permission Denied !"
func (r *ReportHandler) UpdateExistingReport(c *gin.Context){
	r.RedisLog.Println("INFO: received http request on UpdatingReport. ")
	
	id :=c.Param("id")

	var reportDetails model.ReportDetails
	if err :=c.BindJSON(&reportDetails); err != nil {
		r.RedisLog.Println("ERROR: failed to decode request body on UpdatingReport :",err)
		c.AbortWithStatusJSON(400, gin.H{
			"error" : "Invalid report details !",
		})
		return
	}
	routingKey:=reportDetails.Title+".updated"
	if err :=r.RabbitRepo.PublishReport(routingKey, reportDetails, "updated", id); err != nil {
		r.RedisLog.Println("ERROR: failed to publish Report details on UpdatingReport :",err)
		c.AbortWithStatusJSON(500, gin.H{
			"error" : "Unable to publish report details !",
		})
		return
	}
	r.RedisLog.Println("INFO: Report Published Succesfully with ID", id," on UpdatingReport ..")
	c.IndentedJSON(202, model.ReportResponse{Message: "Request Accepted .."})
}
// @Router  				/reports/:id  [delete]
// @Summary 			Deletes  Reports  
// @Description 		This method is responsible for deleting  reports   
// @Security 				BearerAuth
// @Tags					 REPORT
// @Produce				  json
// @Param 					token    header    string    true    "Reporter Token"
// @Param 					id    		path    string    true    "Report ID"
// @Success					202 	{object}   string		"Request Accepted .."
// @Failure					 500 {object}    error	"Unable to publish report details !"
// @Failure					 403 {object}    error		"Permission Denied !"
func (r *ReportHandler) DeleteExistingReport(c *gin.Context){
	r.RedisLog.Println("INFO: received http request on DeletingReport. ")	
	id :=c.Param("id")
	routingKey:=id+".deleted"
	if err :=r.RabbitRepo.PublishReport(routingKey, model.ReportDetails{} , "deleted", id); err != nil {
		r.RedisLog.Println("ERROR: failed to publish Report details on DeletingReport :",err)
		c.AbortWithStatusJSON(500, gin.H{
			"error" : "Unable to publish report details !",
		})
		return
	}
	r.RedisLog.Println("INFO: Report Published Succesfully with ID", id," on DeletingReport ..")
	c.IndentedJSON(202, model.ReportResponse{Message: "Request Accepted .."})
}
// @Router  				/reports  [get]
// @Summary 			Gets All  Reports  
// @Description 		This method is responsible for getting all  reports   
// @Security 				BearerAuth
// @Tags					 REPORT
// @Produce				  json
// @Param 					token    header    string    true    "Reporter Token"
// @Success					200 	{object}   []model.ReportDetails	
// @Failure					 500 {object}    error	"Unable to get all  reports details !"
// @Failure					 404 {object}    string		"There is No Reports Exists !"
// @Failure					 403 {object}    error		"Permission Denied !"
func (r *ReportHandler) GetAllReports(c *gin.Context){
	r.RedisLog.Println("INFO: received http request on GetAllReports. ")	

	listreports,  err :=r.MongoRepo.GetAllReportsFromMongoDB()
	if  err != nil {
		r.RedisLog.Println("ERROR: failed to get All Report details From mongoDB :",err)
		c.AbortWithStatusJSON(500, gin.H{
			"error" : "Unable to get all  reports details !",
		})
		return
	}
	if len(listreports) == 0 {
		c.IndentedJSON(http.StatusNotFound, model.ReportResponse{Message: "There is No Reports Exists !"})
	}
	r.RedisLog.Println("INFO: Report  Succesfully sent to client   on GetAllReports ..")
	c.IndentedJSON(200 , listreports)
}
// @Router  				/reports/search  [get]
// @Summary 			Gets  multiple  Reports  
// @Description 		This method is responsible for getting Multiple  reports by filter  
// @Security 				BearerAuth
// @Tags					 REPORT
// @accept				  	json
// @Produce				  json
// @Param 					token    header    string    true    "Reporter Token"
// @Param				    body 	body    model.ReportFilter	true "Filter details"
// @Success					200 	{object}   []model.ReportDetails	
// @Failure					 500 {object}    error	"Unable to get  reports details  by filter!"
// @Failure					 404 {object}    string		"There is No Reports Exists !"
// @Failure					 403 {object}    error		"Permission Denied !"
func (r *ReportHandler) GetReportsByFilter(c *gin.Context){
	r.RedisLog.Println("INFO: received http request on GetReportsByFilter. ")	

	var filter model.ReportFilter
	if err :=c.BindJSON(&filter); err != nil {
		r.RedisLog.Println("ERROR: failed to decode request body on GetReportsByFilter :",err)
		c.AbortWithStatusJSON(400, gin.H{
			"error" : "Invalid Filter Details !",
		})
		return
	}

	listreports,  err :=r.MongoRepo.GetReportsByFilterFromMongoDB(filter)
	if  err != nil {
		r.RedisLog.Println("ERROR: failed to get Reports By filter details From mongoDB :",err)
		c.AbortWithStatusJSON(500, gin.H{
			"error" : "Unable to get  reports details  by filter!",
		})
		return
	}
	if len(listreports) == 0 {
		c.IndentedJSON(404, model.ReportResponse{Message: "There is No Reports Exists !"})
	}
	r.RedisLog.Println("INFO: Report  Succesfully sent to client   on GetReportsByFilter ..")
	c.IndentedJSON(200 , listreports)
}
// @Router  				/employees/register  [post]
// @Summary 			Registers   New Employees    
// @Description 		This method is responsible for registering new  employees   
// @Security 				BearerAuth
// @Tags					 EMPLOYEE
// @accept					json
// @Produce				  json
// @Param 					body    body    model.EmployeeDetail    true  "Employee Details"
// @Success					201 	{object}   model.EmployeeID		
// @Failure					 400 {object}    error		"Invalid employee details !"
// @Failure					 500 {object}    error	"Unable to register  new employee!"
func (r *ReportHandler) RegisterNewReporters(c *gin.Context){
	r.RedisLog.Println("INFO: received http request on RegisterNewEmployee. ")	
	var employee model.EmployeeDetail
	if err :=c.BindJSON(&employee); err != nil {
		r.RedisLog.Println("ERROR: failed to decode request body on RegisterNewEmployee :",err)
		c.AbortWithStatusJSON(400, gin.H{
			"error" : "Invalid employee details !",
		})
		return
	}
	employeeID,  err :=r.MongoRepo.AddNewEmployeeIntoMongoDB(employee)
	if  err != nil {
		log.Println(err)
		r.RedisLog.Println("ERROR: failed to insert new employee  details into mongoDB :",err)
		c.AbortWithStatusJSON(500, gin.H{
			"error" : "Unable to register  new employee!",
		})
		return
	}

	err =r.RedisClient.AddNewEmployeeIntoRedis(employeeID.Id, employee)
	if  err != nil {
		log.Println(err)
		r.RedisLog.Println("ERROR: failed to insert new employee  details into Redis :",err)
		c.AbortWithStatusJSON(500, gin.H{
			"error" : "Unable to register  new employee!",
		})
		return
	}

	c.IndentedJSON(201, employeeID)
}
// @Router  				/employees/login  [get]
// @Summary 			Employee 		Login       
// @Description 		This method is responsible for managing   employees login   
// @Security 				BearerAuth
// @Tags					 EMPLOYEE
// @Produce				  json
// @Param 					id             header    string    true    "ID"
// @Param 					password    header    string    true    "Password"
// @Success					201 	{object}   model.EmployeeToken		
// @Failure					 400 {object}    error		"Login Failed : Invalid ID or Password !"
// @Failure					 500 {object}    error	"Unable to procces login!"
func (r *ReportHandler) LoginReporters(c *gin.Context){
	id :=c.GetHeader("id")
	password :=c.GetHeader("password")

	role, err :=r.RedisClient.VerifyEmployeeLoginFromRedis(id, password)
	if err != nil {
		r.RedisLog.Println("ERROR: failed to verify  employee  login  Redis :",err)
		c.AbortWithStatusJSON(400, gin.H{
			"error" : "Login Failed : Invalid ID or Password !",
		})
		return
	}

	jwtHandler :=auth.JWTHandler{Role: role, Id: id}
	token, err :=jwtHandler.GenerateToken()
	if err != nil {
		r.RedisLog.Println("ERROR: failed to generate token  for employee in  login :",err)
		c.AbortWithStatusJSON(500, gin.H{
			"error" : "Unable to procces login!",
		})
		return
	}
	c.IndentedJSON(200, model.EmployeeToken{Token: token})
}
// @Router  				/reports/status/:id  [get]
// @Summary 			Gets  multiple  Report's  statuses
// @Description 		This method is responsible for getting Multiple  report's statuses by ID  
// @Security 				BearerAuth
// @Tags					 REPORT
// @Produce				  json
// @Param 					token    header    string    true    "Reporter Token"
// @Param				    id 	   path       string 	 		true  "Report  ID"
// @Success					200 	{object}   []model.ReportStatus	
// @Failure					 500 {object}    error	"Unable to get  reports statuses !"
// @Failure					 404 {object}    string		"There is No Reports status Exists !"
// @Failure					 403 {object}    error		"Permission Denied !"
func (r *ReportHandler) GetReportsStatusById( c *gin.Context){
	r.RedisLog.Println("INFO: received http request on GetReportStatusByID. ")	
	
	id :=c.GetHeader("id")

	listReportStatuses,  err :=r.MongoRepo.GetReportsStatusByIdFromMongoDB(id)
	if  err != nil {
		r.RedisLog.Println("ERROR: failed to get  Reports status   From mongoDB :",err)
		c.AbortWithStatusJSON(500, gin.H{
			"error" : "Unable to get  reports statuses !",
		})
		return
	}
	if len(listReportStatuses) == 0 {
		c.IndentedJSON(http.StatusNotFound, model.ReportResponse{Message: "There is No Reports  status Exists !"})
	}
	r.RedisLog.Println("INFO: Reports status  Succesfully sent to client   on GetReportStatusById ..")
	c.IndentedJSON(200 , listReportStatuses)

}