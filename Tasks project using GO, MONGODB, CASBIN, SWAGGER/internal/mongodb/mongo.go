package mongodb

import (
	"context"
	"fmt"
	"tasks/internal/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMongo(mongo_url string) (*mongo.Client, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongo_url))
	if err != nil {
		return nil, fmt.Errorf("error: connecting to mongoDB :%v", err)
	}
	return client, nil
}

type MongoDB struct {
	mDB *mongo.Database
}

func NewMongoDB(m *mongo.Database) *MongoDB {
	return &MongoDB{mDB: m}
}

func (m *MongoDB) AddTasksToMongoDB(listTasks []model.TaskInfo) (*model.TaskResponse, error) {
	taskCollection := m.mDB.Collection("tasks")

	Tasks := make([]interface{}, len(listTasks))
	for index, task := range listTasks {
		Tasks[index] = task
	}

	_, err := taskCollection.InsertMany(context.TODO(), Tasks)
	if err != nil {
		return nil, fmt.Errorf("error : inserting tasks to mongoDB: %v", err)
	}

	return &model.TaskResponse{Message: "Tasks added succesfully ."}, nil
}

func (m *MongoDB) GetTasksByEmailFromMongoDB(email string) ([]model.TaskInfo, error) {
	taskCollection := m.mDB.Collection("tasks")

	cursor, err := taskCollection.Find(context.TODO(), bson.M{"assignedto.email": email})
	if err != nil {
		return nil, fmt.Errorf("error: finding tasks from mongoDB: %v", err)
	}
	defer cursor.Close(context.TODO())

	var listTasks []model.TaskInfo
	for cursor.Next(context.TODO()) {
		var task model.TaskInfo
		if err := cursor.Decode(&task); err != nil {
			return nil, fmt.Errorf("error : decoding tasks from mongoDB: %v", err)
		}
		listTasks = append(listTasks, task)
	}
	if len(listTasks) == 0 {
		return nil, nil
	}
	return listTasks, nil
}

func (m *MongoDB) GetAllTasksByDateFromMongoDB(date string) ([]model.TaskInfo, error) {
	taskCollection := m.mDB.Collection("tasks")

	dueDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, fmt.Errorf("error: parsing date into time format : %v", err)
	}
	cursor, err := taskCollection.Find(context.TODO(), bson.M{"duedate": bson.M{"$lte": dueDate.Format("2006-01-02")}})
	if err != nil {
		return nil, fmt.Errorf("error: finding tasks from mongoDB: %v", err)
	}
	defer cursor.Close(context.TODO())

	var listTasks []model.TaskInfo
	for cursor.Next(context.TODO()) {
		var task model.TaskInfo
		if err := cursor.Decode(&task); err != nil {
			return nil, fmt.Errorf("error : decoding tasks from mongoDB: %v", err)
		}
		listTasks = append(listTasks, task)
	}
	if len(listTasks) == 0 {
		return nil, nil
	}
	return listTasks, nil
}

func (m *MongoDB) GetUnfinishedSubtasksFromMongoDB() ([]model.TaskInfo, error) {
	taskCollection := m.mDB.Collection("tasks")

	cursor, err := taskCollection.Find(context.TODO(), bson.M{
		"subtasks": bson.M{
			"$elemMatch": bson.M{
				"status": bson.M{
					"$in": []string{"In Progress", "Pending"},
				},
			},
		}})
	if err != nil {
		return nil, fmt.Errorf("error: finding tasks from mongoDB: %v", err)
	}
	defer cursor.Close(context.TODO())

	var listTasks []model.TaskInfo
	for cursor.Next(context.TODO()) {
		var task model.TaskInfo
		if err := cursor.Decode(&task); err != nil {
			return nil, fmt.Errorf("error : decoding tasks from mongoDB: %v", err)
		}
		listTasks = append(listTasks, task)
	}
	if len(listTasks) == 0 {
		return nil, nil
	}
	return listTasks, nil
}

func (m *MongoDB) UpdatesubtasksStatusInMongoDB(Filtertask model.FilterSubtask) (*model.TaskResponse, error) {
	taskCollection := m.mDB.Collection("tasks")

	filter := bson.M{
		"title":          Filtertask.Title,
		"subtasks.title": Filtertask.SubTask.Title,
	}
	updateStatus := bson.M{"$set": bson.M{"subtasks.$.title": Filtertask.SubTask.Status}}

	result, err := taskCollection.UpdateOne(context.TODO(), filter, updateStatus)
	if err != nil {
		return nil, fmt.Errorf("error: updating subtask status: %v", err)
	}
	if result.ModifiedCount == 0 {
		return nil, fmt.Errorf("no task found ")
	}
	return &model.TaskResponse{Message: "Subtask status updated succesfully ."}, nil
}

func (m *MongoDB) ChangeEmployeeOfTaskInMongoDB(newEmployee model.ChangeEmployee) (*model.TaskResponse, error) {
	taskCollection := m.mDB.Collection("tasks")

	update := bson.M{"&set": bson.M{
		"assignedto": newEmployee.Title,
	}}
	result, err := taskCollection.UpdateOne(context.TODO(), bson.M{"title": newEmployee.Title}, update)
	if err != nil {
		return nil, fmt.Errorf("error: updating assigned_to: %v", err)
	}
	if result.ModifiedCount == 0 {
		return nil, fmt.Errorf("no task found ")
	}
	return &model.TaskResponse{Message: "Responsible Employee Changed  succesfully ."}, nil
}

func (m *MongoDB) AddNewSubtaskToTaskInMongoDB(newSubtask model.FilterSubtask) (*model.TaskResponse, error) {
	taskCollection := m.mDB.Collection("tasks")

	addSubtask := bson.M{"$push": bson.M{
		"subtasks": newSubtask.SubTask,
	}}
	result, err := taskCollection.UpdateOne(context.TODO(), bson.M{"title": newSubtask.Title}, addSubtask)
	if err != nil {
		return nil, fmt.Errorf("error: adding new subtask: %v", err)
	}
	if result.ModifiedCount == 0 {
		return nil, fmt.Errorf("no task found ")
	}
	return &model.TaskResponse{Message: "New Subtask added  succesfully ."}, nil
}

func (m *MongoDB) DeleteSubtaskOfTaskFromMongoDB(subtask model.FilterSubtask) (*model.TaskResponse, error) {
	taskCollection := m.mDB.Collection("tasks")

	filter := bson.M{"title": subtask.Title}
	deleteSubtask := bson.M{"$pull": bson.M{
		"subtasks": bson.M{
			"$elemMatch": bson.M{
				"title": subtask.SubTask.Title,
			},
		},
	}}

	result, err := taskCollection.UpdateOne(context.TODO(), filter, deleteSubtask)
	if err != nil {
		return nil, fmt.Errorf("error: Deleting subtask: %v", err)
	}
	if result.ModifiedCount == 0 {
		return nil, fmt.Errorf("no task found")
	}
	return &model.TaskResponse{Message: "Subtask deleted  succesfully ."}, nil
}

func (m *MongoDB) DeleteEmployeesAllTasksFromMongoDB(email string) (*model.TaskResponse, error) {
	taskCollection := m.mDB.Collection("tasks")

	result, err := taskCollection.DeleteOne(context.TODO(), bson.M{"assignedto.email": email})
	if err != nil {
		return nil, fmt.Errorf("error: Deleting subtask: %v", err)
	}
	if result.DeletedCount == 0 {
		return nil, fmt.Errorf("no task found")
	}
	return &model.TaskResponse{Message: "All tasks  deleted  succesfully  for employee."}, nil
}

func (m *MongoDB) DeleteExpiredAllTasksFromMongoDB(givendate string) (*model.TaskResponse, error) {
	taskCollection := m.mDB.Collection("tasks")

	date, err := time.Parse("2006-01-02", givendate)
	if err != nil {
		return nil, fmt.Errorf("error: parsing date into time format : %v", err)
	}
	filter :=bson.M{"$and" : []bson.M{
		{"duedate" : bson.M{
			"$lte" : date.Format("2006-01-02")}},
		{"status" : "Completed"},
	}}
	result, err :=taskCollection.DeleteMany(context.TODO(),filter)
	if err != nil {
		return nil, fmt.Errorf("error: Deleting Tasks: %v", err)
	}
	if result.DeletedCount == 0 {
		return nil, fmt.Errorf("no task found")
	}
	return &model.TaskResponse{Message: "All expiredtasks  deleted  succesfully  By DATE."}, nil
}
