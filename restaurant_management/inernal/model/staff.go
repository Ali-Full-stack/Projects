package model


type EmployeeInfo struct{
	Role  string 	`json:"role" bson:"role" `
	Name string  `json:"name" bson:"name"`
	Email  string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password" `
}

type  EmployeeResponse struct{
	Message string `json:"message"` 
}
type EmployeeID struct{
	ID string `json:"id" bson:"id"`
	Status string `json:"status" bson:"status"`
}

type EmployeeLogin struct{
	Role  string 	`json:"role" bson:"role" `
	Name string  `json:"name" bson:"name"`
}
 