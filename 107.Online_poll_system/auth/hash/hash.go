package hash

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func GenerateHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", fmt.Errorf("bcrypt error: %v", err)
	}

	return string(bytes), nil
}

func ValidHash(hash, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return err
	}
	return nil
}
