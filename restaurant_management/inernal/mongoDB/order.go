package mongodb

import (
	"context"
	"fmt"
	"restaurant-service/inernal/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (m *MongoRepo) AddNewOrderIntoMongoDB(order *model.GuestOrder) error {
	orderCollection :=m.MongoClient.Database("restaurant").Collection("orders")
	order.Starter.Status = "Preparing"
	order.Main.Status = "Pending"
	order.Time = time.Now().Format(time.ANSIC)

	_, err :=orderCollection.InsertOne(context.Background(), order)
	if err != nil {
		return fmt.Errorf("failed to insert order into mongoDB: %v",err)
	}
	return nil
}
func (m *MongoRepo) UpdateOrdersStatusInMongoDB(tableNumber int32 , order,  status string)(error){
	orderCollection :=m.MongoClient.Database("restaurant").Collection("orders")
	
	_, err :=orderCollection.UpdateOne(context.Background(), bson.M{ "tableNumber": tableNumber}, bson.M{"$set" : bson.M{order+".status" : status}})
	if err != nil {
		return fmt.Errorf("failed to update order  status in mongoDB: %v",err)
	}
	return nil
}

func (m *MongoRepo) GetOrdersStatusFromMongoDB(status model.GetStatus)(*model.StatusResponse, error){
	orderCollection :=m.MongoClient.Database("restaurant").Collection("orders")

	result :=orderCollection.FindOne(context.Background(),
	bson.M{"tableNumber" : status.TableNumber},
	options.FindOne().SetProjection(bson.M{status.Course + ".status" : 1, "_id" : 0 }))

	if result.Err() != nil{
		return nil, fmt.Errorf("no orders found in MongoDB: %v",result.Err())
	}
	var statusResponse model.StatusResponse
	statusResponse.TableNumber = status.TableNumber
	statusResponse.Course = status.Course
	var results bson.M
	if err :=result.Decode(&results); err !=nil {
		return nil, fmt.Errorf("failed to decode order status: %v",err)
	}
	statusResponse.Status =results[status.Course].(bson.M)["status"].(string)
	return &statusResponse, nil
}


