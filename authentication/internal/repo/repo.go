package repo

import (
	"auth/auth"
	"auth/internal/model"
	"database/sql"
	"fmt"
	"time"

	sq"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

type UserRepo struct {
	DB *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{DB: db}
}

func (u *UserRepo) AddNewUserToDatabase(req model.User) (*model.UserId, error) {
	id := uuid.New().String()
	time := time.Now().Format(time.ANSIC)
	hash, err := auth.GenerateHash(req.Password)
	if err != nil {
		return nil, fmt.Errorf("generate hash error :%v", err)
	}

	query, args, err := sq.Insert("users").Columns("id, name, email, hash_password, created_at").Values(id, req.Name, req.Email, hash, time).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("insert error : %v", err)
	}
	_, err = u.DB.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("db exec error : %v", err)
	}

	return &model.UserId{Status: "registration succesfull", Id: id}, nil
}


func (u *UserRepo) CheckUserFromDatabase(req model.UserLogin)error{

	query, args , err :=sq.Select("hash_password").From("users").Where(sq.Eq{"id" : req.Id}).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("error select :%v",err)
	}
	var hash string
	err =u.DB.QueryRow(query, args...).Scan(&hash)
	if err != nil {
		return fmt.Errorf("error queryrow %v",err)
	}

	err =auth.ValidHash(hash, req.Password)
	if err != nil {
		return fmt.Errorf("error validation hash: %v",err)
	}
	return nil
}

func (u *UserRepo) UpdatePasswordOfUserInDatabase(req model.UpdatePassword)error{

	hash, err := auth.GenerateHash(req.Password)
	if err != nil {
		return fmt.Errorf("error: generating hash password: %v",err)
	}
	query , args, err :=sq.Update("users").Set("hash_password", hash).Where(sq.Eq{"email": req.Email}).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("error update password: %v",err)
	}
	_, err =u.DB.Exec(query,args...)
	if err != nil {
		return fmt.Errorf("error exec DB: %v",err)
	}
	return nil
}