package utils

import (
	"errors"
	"time"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("aji-ganteng")

type JWTClaim struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateToken(userID int64) (string, error){
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &JWTClaim{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ValidateToken(signedToken string) (*JWTClaim, error){
	token, err := jwt.ParseWithClaims(
			signedToken,
			&JWTClaim{},
			func(token *jwt.Token) (interface{}, error){
				return jwtSecret, nil
			},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaim)
	if !ok || !token.Valid{
		return nil, errors.New("token tidak valid")
	}
	return claims, nil
}
