package redisdb

import (
	"context"
	"fmt"
	"restaurant-service/inernal/model"

	r "github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client *r.Client
}

func ConnectRedis() *RedisClient {
	return &RedisClient{Client: r.NewClient(&r.Options{Addr: "127.0.0.1:6379"})}
}

func (r *RedisClient) AddEmployeeForLogin(id string, emp model.EmployeeInfo)(error){
	err := r.Client.HSet(context.Background(), id, map[string]interface{}{
		"role" : emp.Role,
		"name" : emp.Name,
		 "password" : emp.Password,
	} )
	if err != nil {
		return fmt.Errorf("error : failed to add employee info into redisDB: %v",err)
	}
	return nil
}

func (r *RedisClient) VerifyEmployeeAndReturnInfo(id string, password string )(*model.EmployeeLogin , error){
	result, err   :=r.Client.HGetAll(context.Background(), id).Result()
	if err != nil {
		return nil , fmt.Errorf("error : failed to get employee info from redisDB: %v",err)
	}
	if  result["password"] != password {
		return nil, fmt.Errorf("error: Password is incorrect ")
	}
	return &model.EmployeeLogin{Role: result["role"], Name: result["name"]}, nil
}