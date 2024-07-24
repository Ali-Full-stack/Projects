package model



type ReportDetails struct{
	Channel    string `json:"channel" bson:"channel"`
	Title    string `json:"title" bson:"title"`
	Reporter  string  `json:"reporter" bson:"reporter"`
	About    string    `json:"about" bson:"about"`
	Date  string `json:"date" bson:"date"`
}

type ReportResponse struct{
	Message  string `json:"message"`
}

type ReportFilter struct{
	FilterName  string `json:"filterName" bson:"filterName"`
	FilterValue  string `json:"filterValue" bson:"filterValue"`
}
type ReportStatus struct{
	Status string `json:"status" bson:"status"`
}

type EmployeeDetail struct{
	Fullname  string `json:"fullname"  bson:"fullname"`
	Role    string `json:"role"  bson:"role"`
	TVChannel  string  `json:"tvChannel"  bson:"tvChannel"`
	Password string `json:"password" bson:"password"`
}

type EmployeeID struct{
	Id  string  `json:"id" bson:"id"`
	Status string `json:"status" bson:"status"`
}

type EmployeeToken struct{
	Token  string `json:"token" bson:"token"`
}