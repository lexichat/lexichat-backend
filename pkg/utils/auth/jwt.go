package utils

import (
	"lexichat-backend/pkg/config"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte(config.JWT_SECRET_KEY)

type Claims struct {
    UserID string `json:"userId"`
    jwt.StandardClaims
}

func GenerateJWT(userID string) (string, error) {
    claims := &Claims{
        UserID: userID,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: 0, // Never expire
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtKey)
}