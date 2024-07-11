package model


type Error400 struct{
	Error 		error  `json:"error"`
}

type Error500 struct{
	InternalError  error `json:"internalerror"`
}
