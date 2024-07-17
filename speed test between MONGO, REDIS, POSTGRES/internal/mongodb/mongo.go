package mongodb

import (
	"context"
	"fmt"
	"speed-test/internal/model"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo struct {
	Client *mongo.Client
}

func NewMongo(mc *mongo.Client) *Mongo {
	return &Mongo{Client: mc}
}

func ConnectMongoDb(mongo_url string) (*mongo.Client, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongo_url))
	if err != nil {
		return nil, fmt.Errorf("error: connecting to mongoDB:%v", err)
	}
	return client, nil
}
func (m *Mongo) AddFilmsToMongo(req model.ListOfFilms, wg *sync.WaitGroup) (*model.Response, error) {
	defer wg.Done()
	start := time.Now()

	collectionFilm := m.Client.Database("film").Collection("films")

	for _, v := range req.Films {
		_, err := collectionFilm.InsertOne(context.TODO(), v)
		if err != nil {
			return nil, fmt.Errorf("error: inserting films to mongoDB: %v", err)
		}
	}

	fmt.Println("MongoDB speed [ POST ]: ", time.Since(start))
	return &model.Response{Message: "succesfully added to mongoDB"}, nil
}
