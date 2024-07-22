package mongodb

import (
	"context"
	"fmt"
	"restaurant-service/inernal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (m *MongoRepo) AddNewEmployeeIntoMongoDB(employee model.EmployeeInfo)(*model.EmployeeID, error){
	employeeCollection :=m.MongoClient.Database("restaurant").Collection("employees")

	result , err :=employeeCollection.InsertOne(context.Background(), employee)
	if err != nil {
		return nil, fmt.Errorf("failed to insert new employee to mongoDB: %v",err)
	}
	id :=result.InsertedID.(primitive.ObjectID).Hex()
	return &model.EmployeeID{ID: id, Status: "Active"}, nil
}

func (m *MongoRepo) UpdateEmployeeInfoInMongoDb(employee model.EmployeeInfo, id string)(*model.EmployeeResponse, error){
	employeeCollection :=m.MongoClient.Database("restaurant").Collection("employees")
	idObj, _ :=primitive.ObjectIDFromHex(id)

	result, err :=employeeCollection.UpdateOne(context.Background(), bson.M{"_id": idObj}, bson.M{"$set" : employee})
	if err != nil {
		return nil, fmt.Errorf("failed to update  employee in mongoDB: %v",err)
	}
	if result.ModifiedCount == 0 {
		return nil, fmt.Errorf("no Employee found with ID: %v",id)
	}
	return &model.EmployeeResponse{Message: "Employee Info Updated Successfully ."}, nil
}

func (m *MongoRepo) DeleteEmployeeFromMongoDb( id string)(*model.EmployeeResponse, error){
	employeeCollection :=m.MongoClient.Database("restaurant").Collection("employees")
	idObj, _ :=primitive.ObjectIDFromHex(id)

	result, err :=employeeCollection.DeleteOne(context.Background(), bson.M{"_id": idObj})
	if err != nil {
		return nil, fmt.Errorf("failed to Delete  employee in mongoDB: %v",err)
	}
	if result.DeletedCount == 0 {
		return nil, fmt.Errorf("no Employee found with ID: %v",id)
	}
	return &model.EmployeeResponse{Message: "Employee Deleted Successfully ."}, nil
}
