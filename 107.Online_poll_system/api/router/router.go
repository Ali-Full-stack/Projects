package router

import (
	"context"
	"log"
	"net/http"
	"os"
	"poll-service/api/handler"
	"poll-service/internal/storage/mongodb"
	"poll-service/internal/storage/postgres"
	"poll-service/internal/storage/redisdb"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

func Routes() {
	mux := http.NewServeMux()

	mongoRepo, err :=mongodb.NewMongoRepo(os.Getenv("mongo_url"))
	if err != nil {
		log.Fatal(err)
	}
	defer mongoRepo.Client.Disconnect(context.Background())

	postgres, err :=postgres.ConnectPostgres(os.Getenv("postgres_url"))
	if err != nil {
		log.Fatal(err)
	}
	defer postgres.Close()

	redisClient :=redisdb.ConnectRedis(os.Getenv("redis_url"))
	
	handler :=handler.NewHandler(mongoRepo, postgres, redisClient)
	
	//USERS
	mux.HandleFunc("POST /users/register", handler.UserRegister)
	mux.HandleFunc("POST /users/login", handler.UserLogin)


	//POLLS
	mux.HandleFunc("POST /polls", handler.CreatePoll)


	log.Println("Server is listening on port :", os.Getenv("server_url"))
	log.Fatal(http.ListenAndServe(os.Getenv("server_url"),mux))
}