package model

type PollInfo struct {
	ID          string        `json:"id" bson:"_id"`
	Title       string        `json:"title" bson:"title"`
	Author      string        `json:"author" bson:"author"`
	Description string        `json:"description" bson:"description"`
	Options     []Option      `json:"options" bson:"options"`
	Comments    []UserComments `json:"comments" bson:"comments"`
	Rates       Rate          `json:"rates" bson:"rates"`
}

type Option struct {
	Option string   `json:"option" bson:"option"`
	Votes  int64    `json:"votes" bson:"votes"`
	User   []string `json:"user" bson:"user"`
}

type UserComments struct {
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
}

type Rate struct {
	Rank  float64      `json:"rank" bson:"rank"`
	Users []RatedUser  `json:"users" bson:"users"`
}

type RatedUser struct {
	Name string  `json:"name" bson:"name"`
	Rate float64 `json:"rate" bson:"rate"`
}