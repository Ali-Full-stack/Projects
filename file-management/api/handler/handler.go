package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"
	"upload-file/api/middleware"
	"upload-file/internal/model"
	"upload-file/internal/postgres"

	"github.com/google/uuid"
)

var MaxUploadSize int64 = 10 * 1024 * 1024

type FileHandler struct {
	Postgres *postgres.Postgres
}

func NewFileHandler(p *postgres.Postgres) *FileHandler {
	return &FileHandler{Postgres: p}
}

func (f *FileHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	file, ok := r.Context().Value("file").(middleware.FileDetails)
	if !ok {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	defer file.Data.Close()

	newFilePath := "./internal/user-file/" + file.Name
	newFile, err := os.Create(newFilePath)
	if err != nil {
		http.Error(w, "Unable to create new file", http.StatusInternalServerError)
		return
	}
	defer newFile.Close()
	_, err = io.Copy(newFile, file.Data)
	if err != nil {
		http.Error(w, "Unable to write data into  new file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(model.FileResponse{
		ID:   uuid.New().String(),
		Name: file.Name,
		Path: newFilePath,
		Created_at: time.Now().Format(time.ANSIC),
	})

}
func    (f *FileHandler) DownloadFile(w http.ResponseWriter, r *http.Request) {

}
