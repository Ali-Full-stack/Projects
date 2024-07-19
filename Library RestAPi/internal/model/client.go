package model

type ClientInfo struct {
	FullName    string      `json:"fullname" bson:"fulname"`
	Email       string      `json:"email" bson:"email"`
	Phone       string      `json:"phone" bson:"phone"`
	Address     Address     `json:"address" bson:"address"`
	RentedBooks []*RentBook `json:"rentedbooks" bson:"rentedbooks"`
}

type ClientID struct {
	Id string `json:"id" bson:"id"`
}

type RentBook struct {
	Title       string `json:"title" bson:"title"`
	Category    string `json:"category" bson:"category"`
	Price_day   string `json:"price_day" bson:"price_day"`
	Duration    int32  `json:"duration" bson:"duration"`
	Given_date  string `json:"given_date" bson:"given_date"`
	Return_date string `json:"return_date" bson:"return_date"`
}

type Address struct {
	Country     string `json:"country" bson:"country"`
	City        string `json:"city" bson:"city"`
	Home_number string `json:"home_number" bson:"home_number"`
}
type ClientCode struct{
	Email string `json:"email" bson:"email"`
	Code int 	`json:"code" bson:"code"`
}
type ClientLogin struct{
	Email string `json:"email" bson:"email"`
	Id string `json:"id" bson:"id"`
}

type ClientToken struct{
	Status string `json:"status" bson:"status"`
	Token string `json:"token" bson:"token"`
}