package router

import (
	"log"
	"net/http"
	"os"
	"upload-file/api/handler"
	"upload-file/api/middleware"
	"upload-file/internal/postgres"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

func Routes(){
	mux :=http.NewServeMux()

	db, err :=postgres.ConnectPostgres(os.Getenv("postgres_url"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	
	postgres :=postgres.NewPostgres(db)
	fileHandler :=handler.NewFileHandler(postgres)
	
	mux.Handle("POST /files/upload",middleware.CheckFileType(http.HandlerFunc( fileHandler.UploadFile)))
	mux.HandleFunc("GET /files/download/{filename}", fileHandler.DownloadFile)

	log.Println("Server is listening on port :", os.Getenv("server_url"))
	log.Fatal(http.ListenAndServe(os.Getenv("server_url"), mux))

}
