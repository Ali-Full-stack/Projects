package mongodb

import (
	"context"
	"encoding/json"
	"fmt"
	"notification-service/internal/model"
	"time"

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

func (m *MongoRepo) AddOrderIntoMongoDB(orderByte []byte) error {
	orderCollection := m.MongoClient.Database("orders").Collection("order")
	var order model.OrderInfo
	err :=json.Unmarshal(orderByte, &order)
	if err != nil {
		return fmt.Errorf("failed to unmarshal order info:%v",err)
	}

	order.Time =time.Now().Format(time.ANSIC)

	_, err = orderCollection.InsertOne(context.Background(), order)
	if err != nil {
		return fmt.Errorf("failed to insert order into mongoDB: %v", err)
	}
	fmt.Println("Order Created Succesfully ..")
	return nil
}
