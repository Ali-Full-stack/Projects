package api

import (
	"context"
	"log"
	"net/http"
	"os"
	"tasks/internal/mongodb"

	_ "tasks/swagger/docs"

	_ "github.com/joho/godotenv/autoload"
	swag "github.com/swaggo/http-swagger"
)

// New ...
// @title  Project: TASKS
// @description This swagger UI was created to manage client tasks
// @version 1.0

func Routes() {
	mux := http.NewServeMux()

	mongoClient, err := mongodb.ConnectMongo(os.Getenv("mongo_url"))
	if err != nil {
		log.Fatal(err)
	}
	defer mongoClient.Disconnect(context.TODO())
	mongoDatabase := mongoClient.Database("task")
	/////////////////////////////////////////////////////////////////////////////////////////
	mongoDB := mongodb.NewMongoDB(mongoDatabase)
	handler := NewHandler(mongoDB)
	////////////////////////////////////////////////////////////////////////////////////////
	mux.Handle("POST /tasks", RolePasswordMiddleware(http.HandlerFunc(handler.CreateMultipleTasks)))
	mux.Handle("GET /tasks/email/{email}", RolePasswordMiddleware(http.HandlerFunc(handler.GetAllTasksByEmpoyeeEmail)))
	mux.Handle("GET /tasks/date/{date}", RolePasswordMiddleware(http.HandlerFunc(handler.GetAllTasksByDate)))
	mux.Handle("GET /tasks/subtasks/pending", RolePasswordMiddleware(http.HandlerFunc(handler.GetUnfinishedSubtasks)))
	mux.Handle("PUT /tasks/subtasks/status", RolePasswordMiddleware(http.HandlerFunc(handler.UpdatesubtasksStatus)))
	mux.Handle("PUT /tasks/employee", RolePasswordMiddleware(http.HandlerFunc(handler.ChangeEmployeeOfTask)))
	mux.Handle("PUT /tasks/subtasks", RolePasswordMiddleware(http.HandlerFunc(handler.AddNewSubtaskToTask)))
	mux.Handle("DELETE /tasks/subtasks", RolePasswordMiddleware(http.HandlerFunc(handler.DeleteSubtaskOfTask)))
	mux.Handle("DELETE /tasks/email/{email}", RolePasswordMiddleware(http.HandlerFunc(handler.DeleteEmployeesAllTasks)))
	mux.Handle("DELETE /tasks/date/{date}", RolePasswordMiddleware(http.HandlerFunc(handler.DeleteExpiredAllTasks)))
	/////////////////////////////////////////////////////////////////////////////////////////

	mux.Handle("/swagger/", swag.WrapHandler)

	/////////////////////////////////////////////////////////////////////////////////////////

	log.Println("Server is listening on port", os.Getenv("server_url"))
	log.Fatal(http.ListenAndServe(os.Getenv("server_url"), mux))
}
