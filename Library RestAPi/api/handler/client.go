package handler

import (
	"bookstore/internal/model"
	"bookstore/internal/repo"
	"encoding/json"
	"log"
	"net/http"
)

type ClientHandler struct {
	ClientRepo *repo.ClientRepo
}

func NewClientHandler(cr *repo.ClientRepo) *ClientHandler {
	return &ClientHandler{ClientRepo: cr}
}
// @Router  				/rent [post]
// @Summary 			Addes Rented Books
// @Description 		This method assign  Rented Books to Client 
// @Security 				BearerAuth
// @Tags					 RENT
// @accept					json
// @Produce				  json
// @Param 					role    header    string    true    "Role"
// @Param 					id    	path        string    true    "Client ID"
// @Param 					body    body    model.RentBook    true  "Rented Book"
// @Success					201 	{object}   model.Response		"Rented book added to client succesfully ."
// @Failure					 400 {object} error "Invalid Request: incorrect book information !!"
// @Failure					 500 {object} error  "Request Denied : Unable to Add RentedBook to Client !"
// @Failure					 403 {object} error "Unauthorized access"
func (c *ClientHandler) RentNewBookFromLibrary(w http.ResponseWriter, r *http.Request) {
	clientID := r.Header.Get("id")
	var RentBook model.RentBook
	if err := json.NewDecoder(r.Body).Decode(&RentBook); err != nil {
		http.Error(w, "Invalid Request: incorrect book information !!", http.StatusBadRequest)
		return
	}
	resp, err := c.ClientRepo.AddRentedBookToClient(clientID, RentBook)
	if err != nil {
		log.Println(err)
		http.Error(w, "Request Denied : Unable to Add RentedBook to Client !", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(resp)

}

func (c *ClientHandler) ReturnBookToLibrary(w http.ResponseWriter, r *http.Request) {

}

func (c *ClientHandler) PayForRentedBook(w http.ResponseWriter, r *http.Request) {}
