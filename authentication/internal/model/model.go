package model

type UserRequest struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Code     int    `json:"code"`
}

type UserResponse struct {
	Message string `json:"status"`
	Code    int    `json:"code,omitempty"`
}
type UserToken struct{
	Status string `json:"status"`
	Token string `json:"token"`
}
type UserLogin struct{
	Id string `json:"id"`
	Password string `json:"password"`
}
type UpdatePassword struct{
	Email string `json:"email"`
	Password string  `json:"password"`
}

type UserId struct {
	Id     string `json:"id"`
	Status string `json:"status"`
}

type UserCode struct{
	Email string `json:"email"`
	Code int 	`json:"code"`
}

type User struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Hash       string `json:"-"`
	Created_at string `json:"created_at"`
}
