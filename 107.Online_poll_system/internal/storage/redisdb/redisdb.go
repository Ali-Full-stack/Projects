package redisdb

import (
	"context"
	"fmt"
	"log"
	"poll-service/internal/model"

	r "github.com/redis/go-redis/v9"
)

type RedisClient struct {
	*r.Client
}

func ConnectRedis(redis_url string) *RedisClient {
	return &RedisClient{r.NewClient(&r.Options{Addr: redis_url})}
}

func (r *RedisClient) AddUserForLogin(req model.UserInfo) error {
	err := r.Client.HSet(context.Background(), req.ID, map[string]interface{}{
		"name":     req.Fullname,
		"password": req.Password,
	}).Err()
	if err != nil {
		return fmt.Errorf("in redis: failed HSET:  %v", err)
	}
	return nil
}

func (r *RedisClient) VerifyUserLogin(req model.UserLogin)(string, error){
	result, err :=r.Client.HGetAll(context.Background(), req.ID,).Result()
	if err != nil {
		log.Println("in redis:",err)
		return "", fmt.Errorf("failed to get user login from redis : %v",err)
	}
	if result["password"] != req.Password {
		return "", fmt.Errorf("invalid password")
	}
	return result["name"], nil
}