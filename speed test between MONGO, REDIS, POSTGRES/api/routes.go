package api

import (
	"context"
	"log"
	"net/http"
	"os"
	"speed-test/internal/mongodb"
	"speed-test/internal/postgres"
	"speed-test/internal/redis"

	_ "github.com/joho/godotenv/autoload"
	_"github.com/lib/pq"
)

func Routes() {
	mux := http.NewServeMux()

	db, err := postgres.OpenSql("postgres", (os.Getenv("postgres_url")))
	if err != nil {
		log.Fatal("Failed postgres connection:", err)
	}
	defer db.Close()
	postgresDB := postgres.NewPostgres(db)
	//////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	mongoClient, err := mongodb.ConnectMongoDb(os.Getenv("mongo_url"))
	if err != nil {
		log.Fatal("failed MongoDB :", err)
	}
	defer mongoClient.Disconnect(context.TODO())
	mongoDB := mongodb.NewMongo(mongoClient)
	//////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	redisClient := redis.ConnectRedis(os.Getenv("redis_url"))
	//////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	handler :=NewHandler(postgresDB, mongoDB, redisClient)

	mux.HandleFunc("POST /films", handler.CreateListOfFilms)
	mux.HandleFunc("GET /films", handler.GetAllFilmsByFilter)

	log.Println("Server is listening on port:", os.Getenv("server_url"))
	log.Fatal(http.ListenAndServe(os.Getenv("server_url"), mux))
}
