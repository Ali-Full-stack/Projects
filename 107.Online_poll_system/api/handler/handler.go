package handler

import (
	"poll-service/internal/storage/mongodb"
	"poll-service/internal/storage/postgres"
	"poll-service/internal/storage/redisdb"
)

type Handler struct {
	Mongo *mongodb.MongoRepo
	Postgres *postgres.Postgres
	Redis     *redisdb.RedisClient
}

func NewHandler(m *mongodb.MongoRepo, p *postgres.Postgres, r *redisdb.RedisClient)*Handler{
	return &Handler{
		Mongo:m,
		Postgres: p,
		Redis:  r,
	}
}



