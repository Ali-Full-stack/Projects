package model

type BookInfo struct{
	ISBN  string `json:"isbn" bson:"isbn"`
	Title string `json:"title" bson:"title"`
	Category string `json:"category" bson:"category"`
	Description string `json:"description" bson:"description"`
	Author  	Author  `json:"author" bson:"author"`
	RentDetails  RentDetails `json:"rentdetails" bson:"rentdetails"`
}

type Author struct{
	FullName 	string `json:"fullname" bson:"fullname"`
	Age           int32  `json:"age" bson:"age"`
	Contact Contact `json:"contact" bson:"contact"`
}

type Contact struct{
	Email  string `json:"email"`
	Facebook string `json:"facebook" bson:"facebook"` 
}

type RentDetails struct{
	Quantity  int32 `json:"quantity" bson:"quantity"`
	Price_day string `json:"price_day" bson:"price_day"`
	Status string `json:"status" bson:"status"`
}

type Response struct{
	Message  string  `json:"message" bson:"message"`
}

type BookStatus struct{
	Title string `json:"title" bson:"title"`
	Status string `json:"status" bson:"status"`
}
type AuthorTotalBook struct{
	Author string `json:"author"  bson:"author"`
	Total_books  int32 `json:"total_books" bson:"total_books"`
}