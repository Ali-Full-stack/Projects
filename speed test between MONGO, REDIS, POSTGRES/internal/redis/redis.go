package redis

import (
	"context"
	"fmt"
	"speed-test/internal/model"
	"sync"
	"time"

	r "github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client *r.Client
}

func ConnectRedis(redis_url string) *RedisClient {
	return &RedisClient{Client: r.NewClient(&r.Options{Addr: redis_url})}
}

func (r *RedisClient) AddFilmsToRedis(req []model.Film, wg *sync.WaitGroup) (*model.Response, error) {
	defer wg.Done()
	start := time.Now()

	for _, v := range req {
		_ , err := r.Client.LPush(context.Background(), "films", map[string]interface{}{
			"title":    v.Title,
			"director": v.Director,
			"genre":    v.Genre,
			"rank":     v.Rank,
		}).Result()
		if err != nil {
			return nil, fmt.Errorf("error: adding film to Redis : %v", err)
		}
	}
	fmt.Println("Redis Speed [ POST ]: ", time.Since(start))
	return &model.Response{Message: "Succesfully added to Redis ."}, nil
}
