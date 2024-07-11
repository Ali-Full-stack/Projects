package redis

import (
	"auth/internal/model"
	"context"
	"fmt"
	"time"

	r "github.com/redis/go-redis/v9"
	"math/rand"
)

type RedisClient struct {
	Client *r.Client
}

func ConnectRedis() *RedisClient {
	return &RedisClient{Client: r.NewClient(&r.Options{Addr: "127.0.0.1:6379"})}
}

func (r *RedisClient) AddUserToRedis(ctx context.Context, req model.UserRequest) error {

	err := r.Client.HSet(ctx, req.Email, map[string]interface{}{
		"name":     req.Name,
		"email":    req.Email,
		"password": req.Password,
		"code":     req.Code,
	}).Err()
	if err != nil {
		return fmt.Errorf("error HSET:  %v", err)
	}
	r.Client.Expire(ctx, req.Email, 3*time.Minute)

	return nil
}

func (r *RedisClient) VerifyCodeAndGetUser(ctx context.Context, email string, usercode int) (*model.User, error) {

	code, err := r.Client.HGet(ctx, email, "code").Int()
	if err != nil {
		return nil, fmt.Errorf("error HGET:%v", err)
	}

	if code == usercode {
		result, err := r.Client.HGetAll(ctx, email).Result()
		if err != nil {
			return nil, fmt.Errorf("error HGETALL: %v", err)
		}
		return &model.User{
			Name:     result["name"],
			Email:    result["email"],
			Password: result["password"],
		}, nil
	}
	return nil, fmt.Errorf("error NOT MATCH")
}

func (r *RedisClient) CodeForRenewingPassword(ctx context.Context,email string)(int, error) {
	code := 10000 + rand.Intn(90000)
	err :=r.Client.Set(ctx, email, code, 3 * time.Minute).Err()
	if err != nil {
		return 0, fmt.Errorf("error: SET in redis :%v",err)
	}
	return code, nil
}

func (r *RedisClient) VerifyCode(ctx context.Context, email string, usercode int)error{
	code, err :=r.Client.Get(ctx, email).Int()
	if err != nil {
		return fmt.Errorf("error GET in redis: %v",err)
	}

	if code ==usercode {
		return nil
	}
	return fmt.Errorf("code mismatch")
}