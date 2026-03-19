package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID int64  `json:"uid"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func SignJWT(secret string, userID int64, role string, ttl time.Duration) (string, error) {
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(secret))
}

func ParseJWT(secret, tokenStr string) (Claims, error) {
	var c Claims
	t, err := jwt.ParseWithClaims(tokenStr, &c, func(token *jwt.Token) (any, error) {
		return []byte(secret), nil
	})

	if err != nil || !t.Valid {
		return Claims{}, errors.New("invalid token")
	}
	return c, nil
}