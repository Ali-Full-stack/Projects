package model

type TaskInfo struct{
	Title  string `json:"title" bson:"title"`
	Description string `json:"description" bson:"description"`
	Status string `json:"status" bson:"status"`
	AssignedTo  Employee `json:"assignedto" bson:"assignedto"`
	DueDate  string `json:"duedate" bson:"duedate"`
	SubTasks  []*SubTask `json:"subtasks" bson:"subtasks"`   
}

type Employee struct{
	FullName   string `json:"fullname" bson:"fullname"`
	Email  	 	string `json:"email" bson:"email"` 
 }

 type ChangeEmployee struct{
	Title string `json:"title" bson:"title"`
	AssignedTo  Employee `json:"assignedto" bson:"assignedto"`
 }
 
 type SubTask struct{
	Title string `json:"title" bson:"title"`
	Status string `json:"status" bson:"status"`
 }
 type FilterSubtask struct{
	Title string `json:"title" bson:"title"`
	SubTask  SubTask `json:"subtask" bson:"subtask"`   
 }
 
 type TaskResponse struct{
	Message string `json:"message"`
 }

