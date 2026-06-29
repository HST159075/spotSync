package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	ID   uint   `json:"id"`
	Role string `json:"role"`
}

func GenerateToken(id uint, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   id,
		"role": role,
		"exp":  time.Now().Add(24 * time.Hour).Unix(),
	})
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func ParseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	claims := token.Claims.(jwt.MapClaims)
	return &Claims{
		ID:   uint(claims["id"].(float64)),
		Role: claims["role"].(string),
	}, nil
}