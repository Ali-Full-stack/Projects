package api

import (
	"auth/auth"
	"auth/internal/model"
	"auth/internal/redis"
	"auth/internal/repo"
	"auth/pkg"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"math/rand"
)

type UserHandler struct {
	UserRepo *repo.UserRepo
	Redis    *redis.RedisClient
}

func NewUserHandler(r *repo.UserRepo, redis *redis.RedisClient) *UserHandler {
	return &UserHandler{UserRepo: r, Redis: redis}
}

// @Router /register [post]
// @Summary Register New User
// @Description This endpoint registers a new user
// @Tags AUTH
// @Accept json
// @Produce json
// @Param body body model.UserRequest true "UserRequest"
// @Success 200 {object} model.UserResponse 
// @Failure 400 {object} error "Invalid request"
// @Failure 500 {object} error "Request Denied !"
func (u *UserHandler) RegisterNewUser(w http.ResponseWriter, r *http.Request) {
	var userReq model.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&userReq); err != nil {
		log.Println("ERROR: on decoding request body !!")
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}
	userReq.Code = 10000 + rand.Intn(90000)
	err := u.Redis.AddUserToRedis(r.Context(), userReq)
	if err != nil {
		log.Println("ERROR: adding user to redis !!", err)
		http.Error(w, "Request denied !", http.StatusInternalServerError)
		return
	}
	err = pkg.SendEmail(userReq.Email, pkg.SendClientCode(userReq.Code, userReq.Name))
	if err != nil {
		log.Println("ERROR: sending email to user !!", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(model.UserResponse{Message: "Confirm the code we have sent to your Email !"})
}

func (u *UserHandler) VerifyCode(w http.ResponseWriter, r *http.Request) {
	var req model.UserCode
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("ERROR: on decoding request body !!")
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	resp, err := u.Redis.VerifyCodeAndGetUser(r.Context(), req.Email, req.Code)
	if err != nil {
		log.Println("ERROR: verifying  code in redis !!")
		http.Error(w, "Request denied !", http.StatusInternalServerError)
		return
	}

	res, err := u.UserRepo.AddNewUserToDatabase(*resp)
	if err != nil {
		log.Println("ERROR: adding user  into database !!", err)
		http.Error(w, "Request denied !", http.StatusInternalServerError)
		return
	}

	if err := pkg.SendEmail(resp.Email, pkg.SendClientID(res.Id, resp.Name)); err != nil {
		log.Println("ERROR: sending email to client :", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)

}

func (u *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var userlogin model.UserLogin
	if err := json.NewDecoder(r.Body).Decode(&userlogin); err != nil {
		log.Println("ERROR: on decoding request body !!")
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	err := u.UserRepo.CheckUserFromDatabase(userlogin)
	if err != nil {
		log.Println("ERROR:  checking user from database !!", err)
		http.Error(w, "user not found", http.StatusBadRequest)
		return
	}

	var jwthandler auth.JWTHandler
	accesstoken, err := jwthandler.GenerateToken()
	if err != nil {
		log.Println("ERROR:  generating jwt token  !!", err)
		http.Error(w, "Request Denied !", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(model.UserToken{Status: "Confirmed.", Token: accesstoken})
}

func (u *UserHandler) ForgetPassword(w http.ResponseWriter, r *http.Request) {
	email := r.PathValue("email")

	code, err := u.Redis.CodeForRenewingPassword(r.Context(), email)
	if err != nil {
		log.Println("ERROR: getting code from redis:", err)
		http.Error(w, "Request Denied !", http.StatusInternalServerError)
		return
	}
	err = pkg.SendEmail(email, pkg.SendClientCode(code, ""))
	if err != nil {
		log.Println("ERROR: sending email to user !!", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(model.UserResponse{Message: "Confirm the code we have sent to your Email !"})
}

func (u *UserHandler) NewPassword(w http.ResponseWriter, r *http.Request) {
	codeStr := r.PathValue("code")
	code, _ := strconv.Atoi(codeStr)

	var userLogin model.UpdatePassword
	if err := json.NewDecoder(r.Body).Decode(&userLogin); err != nil {
		log.Println("ERROR: on decoding request body !!")
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	if err := u.Redis.VerifyCode(r.Context(), userLogin.Email, code); err != nil {
		log.Println("ERROR: verifying code in redis:", err)
		http.Error(w, "Code is Incorrect !", http.StatusBadRequest)
		return
	}

	if err := u.UserRepo.UpdatePasswordOfUserInDatabase(userLogin); err != nil {
		log.Println("ERROR: updating password:", err)
		http.Error(w, "Request Denied, Please Try again Later !", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(model.UserResponse{Message: "New Password Has Been Set Succesfully ."})
}
