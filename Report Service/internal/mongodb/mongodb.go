package mongodb

import (
	"context"
	"fmt"
	"log"
	"rabbitmq-topic/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepo struct {
	MongoClient *mongo.Client
}

func NewMongoRepo(mongo_url string) (*MongoRepo, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongo_url))
	if err != nil {
		return nil, fmt.Errorf("error: failed mongoDB connection: %v", err)
	}
	return &MongoRepo{MongoClient: client}, nil
}

func (m *MongoRepo) AddNewReportIntoMongoDB(report model.ReportDetails) error {
	reportCollection := m.MongoClient.Database("report").Collection("reports")
	reportStatusCollection := m.MongoClient.Database("report").Collection("report_status")

	result, err := reportCollection.InsertOne(context.Background(), report)
	if err != nil {
		return fmt.Errorf("failed to insert new report into mongoDB: %v", err)
	}
	reportID := result.InsertedID.(primitive.ObjectID).Hex()
	_, err = reportStatusCollection.InsertOne(context.Background(), bson.M{"status": "New Report Created with ID :" + reportID})
	if err != nil {
		log.Println("Unable to write report status on Create:", err)
	}
	return nil
}

func (m *MongoRepo) UpdateExistingReportInMongoDB(report model.ReportDetails, id string) error {
	reportCollection := m.MongoClient.Database("report").Collection("reports")
	reportStatusCollection := m.MongoClient.Database("report").Collection("report_status")

	objReportID, _ := primitive.ObjectIDFromHex(id)

	result, err := reportCollection.UpdateOne(context.Background(), bson.M{"_id": objReportID}, bson.M{"$set": report})
	if err != nil {
		return fmt.Errorf("failed to update report in mongoDB: %v", err)
	}
	if result.ModifiedCount == 0 {
		return fmt.Errorf("report does not exist with ID: %v", id)
	}

	_, err = reportStatusCollection.InsertOne(context.Background(), bson.M{"status": "Report Updated  with ID :" + id})
	if err != nil {
		log.Println("Unable to write report status on Update:", err)
	}
	return nil
}
func (m *MongoRepo) DeleteExistingReportFromMongoDB(id string) error {
	reportCollection := m.MongoClient.Database("report").Collection("reports")
	reportStatusCollection := m.MongoClient.Database("report").Collection("report_status")

	objReportID, _ := primitive.ObjectIDFromHex(id)
	result, err := reportCollection.DeleteOne(context.Background(), bson.M{"_id": objReportID})
	if err != nil {
		return fmt.Errorf("failed to delete report from mongoDB: %v", err)
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("report does not exist with ID: %v", id)
	}

	_, err = reportStatusCollection.InsertOne(context.Background(), bson.M{"status": "Report Deleted  with ID :" + id})
	if err != nil {
		log.Println("Unable to write report status on Delete:", err)
	}
	return nil
}

func (m *MongoRepo) GetAllReportsFromMongoDB()([]model.ReportDetails, error){
	reportCollection := m.MongoClient.Database("report").Collection("reports")

	cursor, err :=reportCollection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil ,fmt.Errorf("failed to get all reports from mongoDB:%v",err)
	}
	var listreports []model.ReportDetails
	if err = cursor.All(context.Background(), listreports); err != nil {
		return nil, fmt.Errorf("failed to cursor.All reports in mongoDB: %v",err)
	}
	return listreports, nil
}

func (m *MongoRepo) GetReportsByFilterFromMongoDB(filter model.ReportFilter)([]model.ReportDetails, error){
	reportCollection := m.MongoClient.Database("report").Collection("reports")

	cursor, err :=reportCollection.Find(context.Background(), bson.M{filter.FilterName : filter.FilterValue})
	if err != nil {
		return nil ,fmt.Errorf("failed to get  reports by filter from mongoDB:%v",err)
	}
	var listreports []model.ReportDetails
	if err = cursor.All(context.Background(), listreports); err != nil {
		return nil, fmt.Errorf("failed to cursor.All reports in mongoDB: %v",err)
	}
	return listreports, nil
}

func (m *MongoRepo) AddNewEmployeeIntoMongoDB( employee model.EmployeeDetail)(*model.EmployeeID, error){
	reportCollection := m.MongoClient.Database("report").Collection("employees")

	result, err :=reportCollection.InsertOne(context.Background(), employee)
	if err != nil {
		return nil, fmt.Errorf("failed to insert new employee : %v",err)
	}
	id :=result.InsertedID.(primitive.ObjectID).Hex()

	return &model.EmployeeID{Id: id, Status: "Active"}, nil
}

func (m *MongoRepo) GetReportsStatusByIdFromMongoDB(id string)([]model.ReportStatus, error){
	reportStatusCollection := m.MongoClient.Database("report").Collection("report_status")

	cursor, err :=reportStatusCollection.Find(context.Background(), bson.M{"_id" : id})
	if err != nil {
		return nil, fmt.Errorf("failed to get Report Status with ID %v  error : %v", id, err)
	}
	var listReportsStatus []model.ReportStatus
	err =cursor.All(context.Background(), listReportsStatus)
	if err != nil {
		return nil, fmt.Errorf("failed to cursor.ALL on getting reports status : %v",err)
	}
	return listReportsStatus, nil
}