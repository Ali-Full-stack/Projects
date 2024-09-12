package postgres

import (
	"database/sql"
	"fmt"
)

type Postgres struct{
	*sql.DB
}

func ConnectPostgres(postgres_url string) (*Postgres, error){
	db, err :=sql.Open("postgres", postgres_url)
	if err != nil {
		return nil, fmt.Errorf("failed to open database : %v",err)
	}
	if err :=db.Ping(); err != nil {
		return nil, fmt.Errorf("unable to ping database : %v",err)
	}
	return &Postgres{db}, nil
}
