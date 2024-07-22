package model

type GuestOrder struct {
	TableNumber int32   `json:"tableNumber" bson:"tableNumber"`
	TotalGuest  int32   `json:"totalGuest" bson:"totalGuest"`
	Waiter      string  `json:"waiter" bson:"waiter"`
	Drink     Drink    `json:"drink" bson:"drink"`
	Starter     Starter `json:"starter" bson:"starter"`
	Pizza       Pizza   `json:"pizza" bson:"pizza"`
	Main        Main    `json:"main" bson:"main"`
	Dessert     Dessert `json:"dessert" bson:"dessert"`
	Time        string   `json:"time" bson:"time"`
}

type Drink struct{
	Drink      map[string]int32 `json:"drink"  bson:"drink"`   
	Status      string                   `json:"status" bson:"status"`
}
type Starter struct {
	Salads      map[string]int32 `json:"salads" bson:"salads"`
	Status      string                   `json:"status" bson:"status"`
}

type Pizza struct {
	Pizza       map[string]int32 `json:"pizza" bson:"pizza"`
	Status      string                   `json:"status" bson:"status"`
}

type Main struct {
	Meals       map[string]int32 `json:"meal" bson:"meal"`
	Status      string                   `json:"status" bson:"status"`
}

type Dessert struct {
	Dessert     map[string]int32 `json:"dessert" bson:"dessert"`
	Coffee      map[string]int32 `json:"coffee"  bson:"coffee"`
	Status      string                   `json:"status" bson:"status"`
}

type FireCourse struct {
	TableNumber int32  `json:"tableNumber" bson:"tableNumber"`
	Waiter      string `json:"waiter" bson:"waiter"`
	Course      string `json:"course" bson:"course"`
	Time        string   `json:"time"  bson:"time"`
}

type GetStatus struct{
	TableNumber int32  `json:"tableNumber" bson:"tableNumber"`
	Course string   `json:"course" bson:"course"`  
}

type StatusResponse struct{
	TableNumber int32  `json:"tableNumber" bson:"tableNumber"`
	Course string   `json:"course" bson:"course"`  
	Status      string                   `json:"status" bson:"status"`
}
