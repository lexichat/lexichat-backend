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

func ParseJWT(tokenString string, secret string) (bool, *Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return false, nil, nil
		}
		return false, nil, err
	}

	if !token.Valid {
		return false, nil, nil
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return false, nil, nil
	}

	return true, claims, nil
}

func GetUserIdFromToken(tokenString string) (string, error) {
    isValid, claims, err := ParseJWT(tokenString, string(jwtKey))
    if err != nil {
        return "", err
    }
    if !isValid || claims == nil {
        return "", err
    }
    return claims.UserID, nil
}
