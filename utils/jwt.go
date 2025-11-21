package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var JWTSecret = []byte(os.Getenv("JWT_SECRET"))

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateAuthToken(userID uint) (string, error) {
	// 在函数内部读取，确保 .env 已加载
	secret := []byte(os.Getenv("JWT_SECRET"))

	Claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims)
	return token.SignedString(secret)
}
