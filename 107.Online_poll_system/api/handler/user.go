package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"poll-service/auth/token"
	"poll-service/internal/model"

	"github.com/google/uuid"
)

func (h *Handler) UserRegister(w http.ResponseWriter, r *http.Request) {
	var userRegister model.UserInfo
	if err := json.NewDecoder(r.Body).Decode(&userRegister); err != nil {
		http.Error(w, "invalid user information !!", http.StatusBadRequest)
		return
	}
	userRegister.ID = uuid.New().String()

	err := h.Postgres.AddUserIntoPostgres(userRegister)
	if err != nil {
		log.Println(err)
		http.Error(w, "failed registration proccess !!", http.StatusInternalServerError)
		return
	}
	if err :=h.Redis.AddUserForLogin(userRegister); err != nil {
		log.Println(err)
		http.Error(w, "failed registration proccess !!", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(model.UserID{ID: userRegister.ID})
}

func (h *Handler)  UserLogin(w http.ResponseWriter, r *http.Request){
	var userLogin model.UserLogin
	if err := json.NewDecoder(r.Body).Decode(&userLogin); err != nil {
		http.Error(w, "invalid login information !!", http.StatusBadRequest)
		return
	}
	name, err :=h.Redis.VerifyUserLogin(userLogin)
	if err != nil {
		log.Println(err)
		http.Error(w, "incorrect ID or  password !", http.StatusNotFound)
		return
	}

	token, err :=token.GenerateToken(userLogin.ID, name)
	if err != nil {
		http.Error(w, "failed  login proccess !!", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(model.UserToken{Token: token})

}

