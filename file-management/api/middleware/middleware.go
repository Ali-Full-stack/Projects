package middleware

import (
	"context"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
)

type FileDetails struct{
	Name string
	Data multipart.File
}

var Extensions = map[string]bool{
	".pdf":  true,
	".txt":  true,
	".png":  true,
	".json": true,
	".svg":  true,
	".jpeg": true,
	".go":   true,
}
var MaxUploadSize int64 = 10 * 1024 * 1024

func CheckFileType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(MaxUploadSize)
		if err != nil {
			http.Error(w, "Unable to parse request file: Invalid Size !! ", http.StatusBadRequest)
			return
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "No file uploaded", http.StatusBadRequest)
			return
		}
		defer file.Close()

		fileName := header.Filename

		fileExt := GetExtension(header)
		if !Extensions[fileExt] {
			http.Error(w, " file type not allowed", http.StatusBadRequest)
			return
		}
		ctx := context.WithValue(r.Context(), "file", FileDetails{
			Name: fileName,
			Data: file,
		})

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetExtension(header *multipart.FileHeader) string {
	fileDetails := header.Header.Get("Content-Disposition")
	slice := strings.Split(fileDetails, ";")
	fileType := filepath.Ext(slice[2])
	extension := strings.Trim(fileType, " \"")
	return extension
}
