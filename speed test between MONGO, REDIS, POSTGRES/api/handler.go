package api

import (
	"encoding/json"
	"log"
	"net/http"
	"speed-test/internal/model"
	"speed-test/internal/mongodb"
	"speed-test/internal/postgres"
	"speed-test/internal/redis"
	"sync"
)

type Handler struct {
	Postgres *postgres.Postgres
	MongoDB  *mongodb.Mongo
	Redis    *redis.RedisClient
}

func NewHandler(p *postgres.Postgres, m *mongodb.Mongo, r *redis.RedisClient) *Handler {
	return &Handler{Postgres: p, MongoDB: m, Redis: r}
}

func (h *Handler) CreateListOfFilms(w http.ResponseWriter, r *http.Request) {
	var listfilms model.ListOfFilms
	if err := json.NewDecoder(r.Body).Decode(&listfilms); err != nil {
		http.Error(w, "Invalid Request : incorrect film informtion", http.StatusBadRequest)
		return
	}

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		_, err := h.MongoDB.AddFilmsToMongo(listfilms, &wg)
		if err != nil {
			log.Println("Failed to add films to MongoDB:", err)
		}
	}()

	go func() {
		_, err := h.Postgres.AddFilmsToPostgres(listfilms.Films, &wg)
		if err != nil {
			log.Println("Failed to add films to PostgresDB:", err)
		}
	}()

	go func() {
		_, err := h.Redis.AddFilmsToRedis(listfilms.Films, &wg)
		if err != nil {
			log.Println("Failed to add films to Redis:", err)
		}
	}()
	wg.Wait()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(model.Response{Message: "Films added to all database succesfully ."})
}

func (h *Handler) GetAllFilmsByFilter(w http.ResponseWriter, r *http.Request) {}
