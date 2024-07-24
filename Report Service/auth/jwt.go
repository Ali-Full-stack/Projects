package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	_ "github.com/joho/godotenv/autoload"
)

type JWTHandler struct {
	Role string `json:"role"`
	Id   string `json:"id"`
}

func (j *JWTHandler) GenerateToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"role": j.Role,
			"id":   j.Id,
			"exp":  time.Now().Add(time.Minute * 30).Unix(),
		})

	accessToken, err := token.SignedString([]byte(os.Getenv("secret_key")))
	if err != nil {
		return "", fmt.Errorf("failed to generate jwt token: %v", err)
	}
	return accessToken, nil
}

func VerifyToken(tokenstring string) (string, error) {
	token, err := jwt.Parse(tokenstring, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("secret_key")), nil
	})
	if err != nil && !token.Valid {
		return "", fmt.Errorf("invalid Token")

	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		role, ok := claims["role"].(string)
		if !ok {
			return "", fmt.Errorf("no role found in the JWT")
		}
		return role, nil
	}
	return "", fmt.Errorf("invalid token")
}
