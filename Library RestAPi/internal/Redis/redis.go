package redis

import (
	"bookstore/internal/model"
	"context"
	"fmt"
	"time"

	r "github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client *r.Client
}

func ConnectRedis() *RedisClient {
	return &RedisClient{Client: r.NewClient(&r.Options{Addr: "127.0.0.1:6379"})}
}

func (r *RedisClient) AddUserToRedis(ctx context.Context, req model.ClientInfo, code int64) error {

	err := r.Client.HSet(ctx, req.Email, map[string]interface{}{
		"fullname":   req.FullName,
		"email":      req.Email,
		"Phone":      req.Phone,
		"country":    req.Address.Country,
		"city":       req.Address.City,
		"homeNumber": req.Address.Home_number,
		"code":       code,
	}).Err()
	if err != nil {
		return fmt.Errorf("error HSET:  %v", err)
	}
	r.Client.Expire(ctx, req.Email, 5*time.Minute)

	return nil
}

func (r *RedisClient) VerifyCodeAndGetUser(ctx context.Context, email string, clientcode int) (*model.ClientInfo, error) {

	code, err := r.Client.HGet(ctx, email, "code").Int()
	if err != nil {
		return nil, fmt.Errorf("error HGET:%v", err)
	}

	if code == clientcode {
		result, err := r.Client.HGetAll(ctx, email).Result()
		if err != nil {
			return nil, fmt.Errorf("error HGETALL: %v", err)
		}
		return &model.ClientInfo{
			FullName: result["fullname"],
			Email:    result["email"],
			Phone:    result["phone"],
			Address: model.Address{
				Country:     result["country"],
				City:        result["city"],
				Home_number: result["homeNumber"],
			},
		}, nil
	}
	return nil, fmt.Errorf("error NOT MATCH")
}

func (r *RedisClient) AddClientForLogin(ctx context.Context, email, id string) error {
	err := r.Client.HSet(ctx, email, map[string]interface{}{
		"email": email,
		"id":    id,
	}).Err()
	if err != nil {
		return fmt.Errorf("error HSET:  %v", err)
	}
	return nil
}

func (r *RedisClient) CheckClientInLogin(ctx context.Context, client model.ClientLogin) error {
	result, err := r.Client.HGetAll(ctx, client.Email).Result()
	if err != nil {
		return fmt.Errorf("error HGETALL: %v", err)
	}
	if result["id"] == client.Id {
		return nil
	}
	return fmt.Errorf("client information is incorrect ")
}
