package token

import (
	"api-gateway/pkg/config"
	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	ID   string `json:"id"`
	Role string `json:"role"`
	jwt.StandardClaims
}

func ExtractClaims(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Load().SECRET_KEY_ACCESS), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, err
	}

	return claims, nil
}
