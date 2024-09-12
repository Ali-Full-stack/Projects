package postgres

import (
	"database/sql"
	"fmt"
)

func ConnectPostgres(postgres_url string) (*sql.DB, error){
	db, err :=sql.Open("postgres", postgres_url)
	if err != nil {
		return nil, fmt.Errorf("failed to open database : %v",err)
	}

	if err :=db.Ping(); err !=nil {
		return nil, fmt.Errorf("database connection failed : %v",err)
	}
	return db, err
}