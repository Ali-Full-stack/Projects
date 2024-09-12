package postgres

import (
	"fmt"
	"poll-service/auth/hash"
	"poll-service/internal/model"
)

func (p *Postgres) AddUserIntoPostgres(req model.UserInfo) error {
	hashPassword, err := hash.GenerateHash(req.Password)
	if err != nil {
		return fmt.Errorf("failed to create hash : %v", err)
	}

	query := `insert into users(id, fullname, email, hash_password)
	values($1, $2, $3, $4)`

	_, err = p.DB.Exec(query, req.ID, req.Fullname, req.Email, hashPassword)
	if err != nil {
		return fmt.Errorf("failed to insert user information into postgres : %v",err)
	}
	return nil
}

