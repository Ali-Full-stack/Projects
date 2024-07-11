package storage

import (
	"database/sql"
	"log"
)

func OpenSql(drname, url string) (*sql.DB, error) {
	db, err := sql.Open(drname, url)
	if err != nil {
		log.Println("Unable to open database:", err)
		return nil, err
	}

	if err := db.Ping(); err != nil {
		log.Println("Unable to connect to database:", err)
		return nil, err
	}

	return db, nil
}
