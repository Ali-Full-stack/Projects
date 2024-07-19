package repo

import (
	"bookstore/internal/model"
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthRepo struct {
	MongoClient *mongo.Client
}

func NewAuthRepo(mc *mongo.Client) *AuthRepo {
	return &AuthRepo{MongoClient: mc}
}

func (a *AuthRepo) AddNewClientToMongoDB(client model.ClientInfo) (*model.ClientID, error) {
	authCollection := a.MongoClient.Database(os.Getenv("mongo_db")).Collection("client")

	result, err := authCollection.InsertOne(context.Background(), client)
	if err != nil {
		return nil, fmt.Errorf("error: inserting New Client  to MongoDB: %v", err)
	}
	id :=result.InsertedID.(primitive.ObjectID).Hex()
	return &model.ClientID{Id: id}, nil
}
