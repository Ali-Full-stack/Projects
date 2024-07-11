package api

import (
	"auth/internal/redis"
	"auth/internal/repo"
	"auth/internal/storage"
	"log"
	"net/http"
	"os"

	_ "auth/api/docs"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	swagger "github.com/swaggo/http-swagger"
)

// New ...
// @title  Project: AUTHENTICATION
// @description This swagger UI : Authenticate New User
// @version 1.0
// @name Authorization

func Routes() {
	mux := http.NewServeMux()

	db, err := storage.OpenSql(os.Getenv("driver_name"), os.Getenv("postgres_url"))
	if err != nil {
		log.Fatal("failed DB:", err)
	}
	defer db.Close()

	repo := repo.NewUserRepo(db)

	redis := redis.ConnectRedis()

	handler := NewUserHandler(repo, redis)

	mux.HandleFunc("POST /register", handler.RegisterNewUser)
	// mux.HandleFunc("POST /verify", handler.VerifyCode)
	// mux.HandleFunc("POST /login", handler.LoginUser)
	// mux.HandleFunc("GET /password/{email}", handler.ForgetPassword)
	// mux.HandleFunc("POST /password/{code}", handler.NewPassword)

	mux.Handle("/swagger/", swagger.WrapHandler)

	log.Println("Server is listening on port", os.Getenv("server_url"))
	log.Fatal(http.ListenAndServe(os.Getenv("server_url"), mux))

}
