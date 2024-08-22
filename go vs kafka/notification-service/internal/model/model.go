package model

type OrderInfo struct{
	Title string `json:"title"  bson:"title"`
	Status  string `json:"status"  bson:"status"`
	Time string `json:"time" bson:"time"`
}
