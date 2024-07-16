package api

import (
	"encoding/json"
	"log"
	"net/http"
	"tasks/internal/model"
	"tasks/internal/mongodb"
)

type Handler struct {
	MongoDB *mongodb.MongoDB
}

func NewHandler(m *mongodb.MongoDB) *Handler {
	return &Handler{MongoDB: m}
}
// @Router  				/tasks [post]
// @Summary 			Create Multiple Tasks
// @Description 		This method creates Multiple Tasks
// @Security 				BearerAuth
// @Tags						 TASK
// @accept					json
// @Produce				  json
// @Param 					body    body    []model.TaskInfo    true  "Tasks"
// @Success					201 	{object}   model.TaskResponse		"Tasks added succesfully ."
// @Failure					 400 {object} error "Invalid Request: incorrect tasks information !!"
// @Failure					403 {object} error "Unauthorized access"
func (h *Handler) CreateMultipleTasks(w http.ResponseWriter, r *http.Request) {
	var listTasks []model.TaskInfo
	if err := json.NewDecoder(r.Body).Decode(&listTasks); err != nil {
		http.Error(w, "Invalid Request: incorrect tasks information !!", http.StatusBadRequest)
		return
	}

	resp, err := h.MongoDB.AddTasksToMongoDB(listTasks)
	if err != nil {
		log.Println(err)
		log.Println("error:",err)
		http.Error(w, "Request Denied:  Unable to insert tasks !!", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(resp)
}
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// @Router  				/tasks/email/{email} [get]
// @Summary 			Get Client Tasks 
// @Description 		This method is responsible to get Clients All Tasks by employee email
// @Security 				BearerAuth
// @Tags					 TASK
// @Produce				  json
// @Param 					email    path    	string    true  "Employee Email"
// @Success					200 	{object}   []model.TaskInfo 
// @Failure					 500 {object} error "Request Denied: Unable To Find tasks in MongoDB !"
// @Failure					403 {object} error "Unauthorized access"
func (h *Handler) GetAllTasksByEmpoyeeEmail(w http.ResponseWriter, r *http.Request) {
	email := r.PathValue("email")

	listTasks, err := h.MongoDB.GetTasksByEmailFromMongoDB(email)
	if err != nil {
		log.Println("error:",err)
		http.Error(w, "Request Denied: Unable To Find tasks in MongoDB !", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(listTasks)
}
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// @Router  				/tasks/date/{date} [get]
// @Summary 			Get All Tasks By Date 
// @Description 		This method is responsible to get  All Tasks Before The Given Date 
// @Security 				BearerAuth
// @Tags					 TASK
// @Produce				  json
// @Param 					date    path    	string    true  "Date"
// @Success					200	{object}   []model.TaskInfo 
// @Failure					 500 {object} error "Request Denied: Unable To Find tasks in MongoDB !"
// @Failure					403 {object} error "Unauthorized access"
func (h *Handler) GetAllTasksByDate(w http.ResponseWriter, r *http.Request) {
	date := r.PathValue("date")
	listTasks, err := h.MongoDB.GetAllTasksByDateFromMongoDB(date)
	if err != nil {
		log.Println("error:",err)
		http.Error(w, "Request Denied: Unable To Find tasks in MongoDB !", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(listTasks)
}
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// @Router  				/tasks/subtasks/pending [get]
// @Summary 			Get  SubtasksTasks by Specific Criteria
// @Description 		This method is responsible to get  All subtasks Whch is not finished 
// @Security 				BearerAuth
// @Tags					 TASK
// @Produce				  json
// @Success					200	{object}   []model.TaskInfo 
// @Failure					 500 {object} error "Request Denied: Unable To Find tasks in MongoDB !"
// @Failure					403 {object} error "Unauthorized access"
func (h *Handler) GetUnfinishedSubtasks(w http.ResponseWriter, r *http.Request) {
	listTasks, err := h.MongoDB.GetUnfinishedSubtasksFromMongoDB()
	if err != nil {
		log.Println("error:",err)
		http.Error(w, "Request Denied: Unable To Find tasks in MongoDB !", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(listTasks)
}
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// @Router  				/tasks/subtasks/status [put]
// @Summary 			Update  Subtask's status 
// @Description 		This method is responsible to update  subtask's status by task title 
// @Security 				BearerAuth
// @Tags					 TASK
// @accept 					json
// @Produce				  json
// @Param 					body    body    model.FilterSubtask    true  "Subtask"
// @Success					202 	{object}   model.TaskResponse		"Subtask status updated succesfully ."
// @Failure					 500 {object} error "Request Denied: Unable To update subtasks status  in MongoDB !"
// @Failure					403 {object} error "Unauthorized access"
func (h *Handler) UpdatesubtasksStatus(w http.ResponseWriter, r *http.Request) {
	var FilterTask model.FilterSubtask
	if err := json.NewDecoder(r.Body).Decode(&FilterTask); err != nil {
		http.Error(w, "Invalid Request: incorrect Subtask information !!", http.StatusBadRequest)
		return
	}
	resp, err := h.MongoDB.UpdatesubtasksStatusInMongoDB(FilterTask)
	if err != nil {
		log.Println("error:",err)
		http.Error(w, "Request Denied: Unable To update subtasks status  in MongoDB !", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(202)
	json.NewEncoder(w).Encode(resp)
}
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// @Router  				/tasks/employee  [put]
// @Summary 			Update  employee of the task 
// @Description 		This method is responsible to update  assigned employee  od the task
// @Security 				BearerAuth
// @Tags					 TASK
// @accept 					json
// @Produce				  json
// @Param 					body    body    model.ChangeEmployee    true  "Employee"
// @Failure					403 {object} error "Unauthorized access"
// @Success					202 	{object}   model.TaskResponse		"Responsible Employee Changed  succesfully ."
// @Failure					 500 {object} error "Request Denied: Unable To update employee  in MongoDB !"
func (h *Handler) ChangeEmployeeOfTask(w http.ResponseWriter, r *http.Request) {
	var newEmployee model.ChangeEmployee
	if err := json.NewDecoder(r.Body).Decode(&newEmployee); err != nil {
		http.Error(w, "Invalid Request: incorrect Employee information !!", http.StatusBadRequest)
		return
	}
	resp, err := h.MongoDB.ChangeEmployeeOfTaskInMongoDB(newEmployee)
	if err != nil {
		log.Println("error:",err)
		http.Error(w, "Request Denied: Unable To update employee  in MongoDB !", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(202)
	json.NewEncoder(w).Encode(resp)
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// @Router  				/tasks/subtasks  [put]
// @Summary 			Update  Subtasks of the task 
// @Description 		This method is responsible to Add new Subtasks to task 
// @Security 				BearerAuth
// @Tags					 TASK
// @accept 					json
// @Produce				  json
// @Param 					body    body    model.FilterSubtask    true  "Subtask"
// @Success					202 	{object}   model.TaskResponse		"New Subtask added  succesfully ."
// @Failure					 500 {object} error  "Request Denied: Unable To add new Subtask  in MongoDB !"
// @Failure					403 {object} error "Unauthorized access"
func (h *Handler) AddNewSubtaskToTask(w http.ResponseWriter, r *http.Request) {
	var filterSubtask model.FilterSubtask
	if err := json.NewDecoder(r.Body).Decode(&filterSubtask); err != nil {
		http.Error(w, "Invalid Request: incorrect subtasks information !!", http.StatusBadRequest)
		return
	}
	resp, err := h.MongoDB.AddNewSubtaskToTaskInMongoDB(filterSubtask)
	if err != nil {
		log.Println("error:", err)
		http.Error(w, "Request Denied: Unable To add new Subtask  in MongoDB !", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(202)
	json.NewEncoder(w).Encode(resp)
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// @Router  				/tasks/subtasks  [delete]
// @Summary 			Delete  Subtasks of the task 
// @Description 		This method is responsible to Delete  Subtasks of the task by status 
// @Security 				BearerAuth
// @Tags					 TASK
// @accept 					json
// @Produce				  json
// @Param 					body    body    model.FilterSubtask    true  "Subtask"
// @Success					202 	{object}   model.TaskResponse		 "Subtask deleted  succesfully ."
// @Failure					 500 {object} error "Request Denied: Unable to Delete subtask  in MongoDB !"
// @Failure					403 {object} error "Unauthorized access"
func (h *Handler) DeleteSubtaskOfTask(w http.ResponseWriter, r *http.Request) {
	var filterSubtask model.FilterSubtask
	if err := json.NewDecoder(r.Body).Decode(&filterSubtask); err != nil {
		http.Error(w, "Invalid Request: incorrect subtasks information !!", http.StatusBadRequest)
		return
	}
	resp, err := h.MongoDB.DeleteSubtaskOfTaskFromMongoDB(filterSubtask)
	if err != nil {
		log.Println("error:", err)
		http.Error(w, "Request Denied: Unable to Delete subtask  in MongoDB !", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(202)
	json.NewEncoder(w).Encode(resp)
}
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// @Router  				/tasks/email/{email}  [delete]
// @Summary 			Delete  clients  tasks 
// @Description 		This method is responsible to Delete  All  tasks by assigned employee's email 
// @Security 				BearerAuth
// @Tags					 TASK
// @Produce				  json
// @Param 					email    path    	string    true  "Employee Email"
// @Success					202 	{object}   model.TaskResponse		"All tasks  deleted  succesfully  for employee."
// @Failure					 500 {object} error  "Request Denied: Unable to Delete tasks  in MongoDB !"
// @Failure					403 {object} error "Unauthorized access"
func (h *Handler) DeleteEmployeesAllTasks(w http.ResponseWriter, r *http.Request){
	email :=r.PathValue("email")

	resp, err := h.MongoDB.DeleteEmployeesAllTasksFromMongoDB(email)
	if err != nil {
		log.Println("error:", err)
		http.Error(w, "Request Denied: Unable to Delete tasks  in MongoDB !", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(202)
	json.NewEncoder(w).Encode(resp)
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// @Router  				/tasks/date/{date}  [delete]
// @Summary 			Delete  all  expired tasks 
// @Description 		This method is responsible to Delete  All  Completed tasks by given date 
// @Security 				BearerAuth
// @Tags					 TASK
// @Produce				  json
// @Param 					date    path    	string    true  "Date"
// @Success					202 	{object}   model.TaskResponse		"All expiredtasks  deleted  succesfully  By DATE."
// @Failure					 500 {object} error  "Request Denied: Unable to Delete tasks  in MongoDB !"
// @Failure					403 {object} error "Unauthorized access"
func (h *Handler) DeleteExpiredAllTasks(w http.ResponseWriter, r * http.Request){
	date :=r.PathValue("date")

	resp, err :=h.MongoDB.DeleteExpiredAllTasksFromMongoDB(date)
	if err != nil {
		log.Println("error:", err)
		http.Error(w, "Request Denied: Unable to Delete tasks  in MongoDB !", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(202)
	json.NewEncoder(w).Encode(resp)
}
