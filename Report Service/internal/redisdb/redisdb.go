package redisdb

import (
	"context"
	"fmt"
	"rabbitmq-topic/internal/model"
	"time"

	r "github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client *r.Client
}

func ConnectRedis() *RedisClient {
	return &RedisClient{Client: r.NewClient(&r.Options{Addr: "127.0.0.1:6379"})}
}

func (rc *RedisClient) Write(p []byte) (n int, err error) {
	logEntry := string(p)
	err = rc.Client.RPush(context.Background(), "logs", logEntry).Err()
	if err != nil {
		return 0, err
	}
	rc.Client.Expire(context.Background(), "logs", 24 * time.Hour)
	return len(p), nil
}

func (rc *RedisClient) AddNewEmployeeIntoRedis(id string , employee model.EmployeeDetail)(error){

	err :=rc.Client.HSet(context.Background(), id, map[string]interface{}{
		"role" : employee.Role,
		"fullname": employee.Fullname,
		"tvChannel": employee.TVChannel,
		"password": employee.Password,
	}).Err()
	if err != nil {
		return fmt.Errorf("failed to insert new employee into Redis : %v",err)
	}
	return nil
}

func (rc *RedisClient) VerifyEmployeeLoginFromRedis(id, empPassword string)(string,error){
	result, err :=rc.Client.HGetAll(context.Background(), id).Result()
	if err != nil {
		return "", fmt.Errorf("failed to  verify employee login from Redis : %v",err)
	}
	role :=result["role"]
	password :=result["password"]
	if password != empPassword {
		return "",fmt.Errorf("password is incorrect ")
	}
	return role,nil
}