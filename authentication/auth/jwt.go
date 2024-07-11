package auth

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	_ "github.com/joho/godotenv/autoload"
)

type JWTHandler struct {
	Role    string `json:"role"`
	Id      string `json:"id"`
}

func (j *JWTHandler) GenerateToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"role":    j.Role,
			"id":      j.Id,
			"exp":     time.Now().Add(time.Minute * 10).Unix(),
		})

	accessToken, err := token.SignedString([]byte(os.Getenv("secret_key")))
	if err != nil {
		log.Println("Failed generating access token:", err)
		return "", err
	}
	return accessToken, nil
}

func VerifyToken(tokenstring string) error {
	token, err := jwt.Parse(tokenstring, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("secret_key")), nil
	})
	if err != nil {
		return err
	}
	if !token.Valid {
		return fmt.Errorf("invalid Token")
	}
	return nil
}
