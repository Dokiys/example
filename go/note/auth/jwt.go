package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
)

type Claims[T any] struct {
	jwt.RegisteredClaims
	Data T `json:"data,omitempty"`
}

func GenJwtToken[T any](key []byte, data T, expire time.Duration) (string, error) {
	claims := Claims[T]{
		Data: data,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expire)),
		},
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(key)
}

func ParseJwtToken[T any](keyFunc jwt.Keyfunc, token string, data T) (*Claims[T], error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims[T]{
		Data: data,
	}, keyFunc)
	if err != nil {
		return nil, err
	}
	if tokenClaims == nil || !tokenClaims.Valid {
		return nil, errors.New("parse token failed")
	}

	claims, ok := tokenClaims.Claims.(*Claims[T])
	if !ok {
		return nil, errors.New("token claims type invalid")
	}

	return claims, nil
}
