package handler

import (
	"bookstore/auth"
	redis "bookstore/internal/Redis"
	"bookstore/internal/model"
	"bookstore/internal/repo"
	"bookstore/pkg"
	"encoding/json"
	"log"
	"net/http"

	"math/rand"
)

type AuthHandler struct {
	AuthRepo    *repo.AuthRepo
	RedisClient *redis.RedisClient
}

func NewAuthHandler(ar *repo.AuthRepo, rc *redis.RedisClient) *AuthHandler {
	return &AuthHandler{AuthRepo: ar, RedisClient: rc}
}
// @Router  				/register [post]
// @Summary 			Register New Client
// @Description 		This method registers New Clients 
// @Security 				BearerAuth
// @Tags					 AUTH
// @accept					json
// @Produce				  json
// @Param 					role    header    string    true    "Role"
// @Param 					body    body    model.ClientInfo    true  "client Information"
// @Success					201 	{object}   model.Response		"Confirm the code we have sent to your Email !"
// @Failure					 400 {object} error  "Invalid Request"
// @Failure					 500 {object} error  "Request denied !"
// @Failure					 403 {object} error "Unauthorized access"
func (a *AuthHandler) RegisterClient(w http.ResponseWriter, r *http.Request) {
	var userReq model.ClientInfo
	if err := json.NewDecoder(r.Body).Decode(&userReq); err != nil {
		log.Println("ERROR: on decoding request body !!")
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}
	code := 10000 + rand.Intn(90000)
	err := a.RedisClient.AddUserToRedis(r.Context(), userReq, int64(code))
	if err != nil {
		log.Println("ERROR: adding user to redis !!", err)
		http.Error(w, "Request denied !", http.StatusInternalServerError)
		return
	}
	err = pkg.SendEmail(userReq.Email, pkg.SendClientCode(code, userReq.FullName))
	if err != nil {
		log.Println("ERROR: sending email to user !!", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(model.Response{Message: "Confirm the code we have sent to your Email !"})
}
////////////////////////////////////////////////////////////////////////////
// @Router  				/verification [post]
// @Summary 			Confirm Client's Code
// @Description 		This method Verifies  Client's Code And sends ID back 
// @Security 				BearerAuth
// @Tags					 AUTH
// @accept					json
// @Produce				  json
// @Param 					role    header    string    true    "Role"
// @Param 					body    body    model.ClientCode    true  "Client Code"
// @Success					201 	{object}   model.ClientID		
// @Failure					 400 {object} error  "Invalid Request"
// @Failure					 500 {object} error  "Request denied !"
// @Failure					 403 {object} error "Unauthorized access"
func (a *AuthHandler) VerifyClientCode(w http.ResponseWriter, r *http.Request) {
	var req model.ClientCode
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("ERROR: on decoding request body !!")
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	client, err := a.RedisClient.VerifyCodeAndGetUser(r.Context(), req.Email, req.Code)
	if err != nil {
		log.Println("ERROR: verifying  code in redis !!")
		http.Error(w, "Request denied !", http.StatusInternalServerError)
		return
	}

	ClientID, err := a.AuthRepo.AddNewClientToMongoDB(*client)
	if err != nil {
		log.Println("ERROR: adding user  into MongoDB !!", err)
		http.Error(w, "Request denied !", http.StatusInternalServerError)
		return
	}
	err =a.RedisClient.AddClientForLogin(r.Context(), req.Email, ClientID.Id)
	if err != nil {
		log.Println("error: Unable to add client to redis for login")
	}
	if err := pkg.SendEmail(client.Email, pkg.SendClientResponse(ClientID.Id, client.FullName)); err != nil {
		log.Println("ERROR: sending email to client :", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ClientID)
}
////////////////////////////////////////////////////////////////////////////
// @Router  				/login [get]
// @Summary 			LOGIN
// @Description 		This method Verifies  Client's login And Returns Token  back 
// @Security 				BearerAuth
// @Tags					 AUTH
// @accept					json
// @Produce				  json
// @Param 					role    header    string    true    "Role"
// @Param 					body    body    model.ClientLogin    true  "Client Login"
// @Success					201 	{object}   model.ClientToken		
// @Failure					 400 {object} error  "Invalid Request"
// @Failure					 404 {object} error  "Client ID not found:   Sign Up  ."
// @Failure					 500 {object} error  "Request denied !"
// @Failure					 403 {object} error "Unauthorized access"
func (a *AuthHandler) LoginClient(w http.ResponseWriter, r *http.Request) {
	var clientlogin model.ClientLogin
	if err := json.NewDecoder(r.Body).Decode(&clientlogin); err != nil {
		log.Println("ERROR: on decoding request body !!")
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	err := a.RedisClient.CheckClientInLogin(r.Context(),clientlogin)
	if err != nil {
		log.Println("ERROR:  checking client from redis !!", err)
		http.Error(w, "Client ID not found:   Sign Up  .", http.StatusNotFound)
		return
	}

	jwthandler := auth.JWTHandler{
		Role: "client",
		Id: clientlogin.Id,
	}
	accesstoken, err := jwthandler.GenerateToken()
	if err != nil {
		log.Println("ERROR:  generating jwt token  !!", err)
		http.Error(w, "Request Denied !", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(model.ClientToken{Status: "Confirmed.", Token: accesstoken})
}
