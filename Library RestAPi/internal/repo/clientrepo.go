package repo

import (
	"bookstore/internal/model"
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ClientRepo struct {
	MongoClient *mongo.Client
}

func NewClientRepo(mc *mongo.Client) *ClientRepo {
	return &ClientRepo{MongoClient: mc}
}

func (c *ClientRepo) AddRentedBookToClient(clientID string, book model.RentBook) (*model.Response, error) {
	clientCollection := c.MongoClient.Database(os.Getenv("mongo_db")).Collection("client")
	id, _ :=primitive.ObjectIDFromHex(clientID)

	result, err :=clientCollection.UpdateOne(context.Background(), bson.M{"_id" :id}, bson.M{
		"$push" :bson.M{"rentedbooks" : book},
	})
	if err != nil {
		return nil, fmt.Errorf("error: unable to update client rentbooks: %v",err)
	}
	if result.ModifiedCount == 0 {
		return nil, fmt.Errorf("error: could not  find client  : %v",err)
	}
	return &model.Response{Message: "Rented book added to client succesfully ."}, nil
}
