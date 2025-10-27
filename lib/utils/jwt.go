package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Claims struct {
	User string `json:"user"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateJWT(Email string, Role string, Secret string, hours int) (string, error) {
	expiredTime := time.Now().Add(time.Duration(hours) * time.Hour)
	claims := Claims{User: Email, Role: Role, RegisteredClaims: jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expiredTime),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(Secret))
	if err != nil {
		return "", err
	}
	return tokenString, err
}
