package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"speed-test/internal/model"
	"sync"
	"time"

	sq "github.com/Masterminds/squirrel"
)

type Postgres struct {
	DB *sql.DB
}

func NewPostgres(db *sql.DB) *Postgres {
	return &Postgres{DB: db}
}

func OpenSql(driverName, url string) (*sql.DB, error) {
	db, err := sql.Open(driverName, url)
	if err != nil {
		log.Println("failed to open database:", err)
		return nil, err
	}

	if err = db.Ping(); err != nil {
		log.Println("Unable to connect to database:", err)
		return nil, err
	}

	return db, err
}

func (p *Postgres) AddFilmsToPostgres(req []model.Film, wg *sync.WaitGroup) (*model.Response, error) {
	defer wg.Done()
	start := time.Now()
	for _, film := range req {
		query, args, err := sq.
			Insert("films").Columns("title", "genre", "director", "rank").Values(film.Title, film.Genre, film.Director, film.Rank).PlaceholderFormat(sq.Dollar).ToSql()
		if err != nil {
			return nil, fmt.Errorf("error: Making query in squirrel: %v", err)
		}
		_, err = p.DB.Exec(query, args...)
		if err != nil {
			return nil, fmt.Errorf("error: inserting films to postgres: %v", err)
		}
	}
	fmt.Println("Postgres  speed [ POST ]: ", time.Since(start))
	return &model.Response{Message: "Succesfully added to Postgres ."}, nil
}
