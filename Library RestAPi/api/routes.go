package api

import (
	"bookstore/api/handler"
	"bookstore/api/middleware"
	redis "bookstore/internal/Redis"
	mongodb "bookstore/internal/mongoDB"
	"bookstore/internal/repo"
	_ "bookstore/swagger/docs"
	"context"
	"log"
	"net/http"
	"os"

	_ "github.com/joho/godotenv/autoload"
	swag "github.com/swaggo/http-swagger"
)

// New ...
// @title  Project: LIBRARY
// @description This swagger UI was created to manage LIBRARY
// @version 1.0

func Routes() {
	mux := http.NewServeMux()
	mongoClient, err := mongodb.ConnectMongo(os.Getenv("mongo_url"))
	if err != nil {
		log.Fatal("failed mongoDB connection :", err)
	}
	defer mongoClient.Disconnect(context.Background())
	/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	redisClient := redis.ConnectRedis()
	authRepo := repo.NewAuthRepo(mongoClient)
	authHandler := handler.NewAuthHandler(authRepo, redisClient)

	mux.Handle("POST /register", middleware.RolePasswordMiddleware(http.HandlerFunc(authHandler.RegisterClient)))
	mux.Handle("POST /verification",middleware.RolePasswordMiddleware(http.HandlerFunc(authHandler.VerifyClientCode)))
	mux.Handle("GET /login", middleware.CheckJwtTokenMiddleware(http.HandlerFunc(authHandler.LoginClient)))

	///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	bookRepo := repo.NewBookRepo(mongoClient)
	bookHandler := handler.NewBookHandler(bookRepo)

	mux.Handle("POST /books", middleware.RolePasswordMiddleware(http.HandlerFunc(bookHandler.CreateMultipleBooks)))
	mux.Handle("GET /books",middleware.RolePasswordMiddleware(http.HandlerFunc( bookHandler.GetMultipleBooks)))
	mux.Handle("GET /books/category/{category}",middleware.RolePasswordMiddleware(http.HandlerFunc( bookHandler.GetMultipleBooksByCategory)))
	mux.Handle("GET /books/author/{author}", middleware.RolePasswordMiddleware(http.HandlerFunc(bookHandler.GetMultipleBooksByAuthor)))
	mux.Handle("PUT /books/status", middleware.RolePasswordMiddleware(http.HandlerFunc(bookHandler.UpdateStatusOfBook)))
	mux.Handle("DELETE /books/{isbn}", middleware.RolePasswordMiddleware(http.HandlerFunc(bookHandler.DeleteBookByISBN)))
	mux.Handle("GET /books/count", middleware.RolePasswordMiddleware(http.HandlerFunc(bookHandler.CountAuthorsBook)))

	///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	clientrRepo := repo.NewClientRepo(mongoClient)
	clientHandler := handler.NewClientHandler(clientrRepo)

	mux.Handle("POST /rent/{id}",  middleware.RolePasswordMiddleware(http.HandlerFunc(clientHandler.RentNewBookFromLibrary)))
	// mux.HandleFunc("PUT /rent", clientHandler.ReturnBookToLibrary)
	// mux.HandleFunc("PUT /rent/pay", clientHandler.PayForRentedBook)
	////////////////////////////////////////////////////////////////////////////////////////
	mux.Handle("/swagger/", swag.WrapHandler)
	////////////////////////////////////////////////////////////////////////////////////////


	log.Println("Server is Listening on port", os.Getenv("server_url"))
	log.Fatal(http.ListenAndServe(os.Getenv("server_url"), mux))
}
