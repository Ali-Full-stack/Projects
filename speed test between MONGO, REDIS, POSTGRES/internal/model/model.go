package model



type Film struct{
	Title string  `json:"title"`
	Genre string `json:"genre"`
	Director string `json:"director"`
	Rank  float32 `json:"Rank"`
}
type ListOfFilms struct{
	Films []Film `json:"films"` 
}

type Filter struct{
	Filter string `json:"filter"`
	Value  string `json:"value"` 
}

type Response struct{
	Message string `json:"message"`
}
