package postgres

import "database/sql"

type Postgres struct {
	*sql.DB
}

func NewPostgres(db *sql.DB)*Postgres{
	return &Postgres{db}
}