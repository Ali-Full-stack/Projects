package model

type FileResponse struct{
	ID string 	`json:"id"`
	Name string `json:"name"`
	Path string `json:"path"`
	Created_at string `json:"created_at"` 
}
