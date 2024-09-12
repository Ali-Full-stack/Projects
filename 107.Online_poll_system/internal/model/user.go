package model

type UserInfo struct {
	ID      string     `json:"id" bson:"id"`
	Fullname     string `json:"fullname" bson:"fullname"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
	Rate     float32   `json:"rank" bson:"rank"`
}

type UserID struct{
	ID string `json:"id" bson:"id"`
}

type UserLogin struct{
	ID string `json:"id" bson:"id"`
	Password string `json:"password" bson:"password"`
}

type UserToken struct {
	Token  string `json:"token" bson:"token"`
}